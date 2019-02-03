package core

import (
	"encoding/hex"
	"fmt"
	"log"

	"github.com/danitello/go-blockchain/chaindb"
	"github.com/danitello/go-blockchain/core/types"
)

const (
	/*genesisData is the data that will go into the genesis Block by default */
	genesisData = "Genesis"
)

/*BlockChain is a complete blockchain
@param Height - how many Blocks it has (highest Index + 1)
@param LastHash - the hash of the most recent Block added to this BlockChain
@param DB - a badger database instance
*/
type BlockChain struct {
	Height   int
	LastHash []byte
	ChainDB  *chaindb.ChainDB
}

/*InitBlockChain instantiates a new instance of a BlockChain
@return the current working BlockChain
*/
func InitBlockChain(address string) *BlockChain {

	db := chaindb.InitDB()
	resChain := &BlockChain{
		Height:   0,
		LastHash: []byte{0},
		ChainDB:  db}

	// If a BlockChain can be found, use it, otherwise make a new one
	if db.HasChain() {
		log.Panic(fmt.Sprintf("BlockChain already exists in %s", chaindb.Dir))
	} else {
		genesisBlock := createGenesisBlock(address)
		fmt.Println("Genesis block signed")

		resChain.saveNewLastBlock(genesisBlock)
	}

	return resChain

}

/*GetBlockChain gets an existing BlockChain from the database
@return the BlockChain
*/
func GetBlockChain() *BlockChain {
	db := chaindb.InitDB()
	resChain := &BlockChain{
		Height:   0,
		LastHash: []byte{0},
		ChainDB:  db}

	if db.HasChain() {
		resChain.LastHash = db.GetLastHash()
		resChain.Height = db.GetBlockWithHash(resChain.LastHash).Index + 1
	} else {
		log.Panic("Error: No BlockChain exists")
	}

	return resChain
}

/*AddBlock adds a new Block to a given BlockChain
@param data - the data to be contained in the Block
*/
func (bc *BlockChain) AddBlock(txns []*types.Transaction) {
	// Create a new block and save it
	newBlock := types.InitBlock(txns, bc.LastHash, bc.Height-1)
	bc.saveNewLastBlock(newBlock)
}

/*saveNewLastBlock saves the new Block to db, and updates BlockChain struct
@param newBlock - the Block to save
*/
func (bc *BlockChain) saveNewLastBlock(newBlock *types.Block) {

	// Update DB
	bc.ChainDB.SaveNewLastBlock(newBlock)

	// Update chain
	bc.LastHash = newBlock.Hash

}

/*createGenesisBlock creates the first Block
@return the Block
*/
func createGenesisBlock(address string) *types.Block {
	cbtx := types.CoinbaseTx(address)
	return types.InitBlock([]*types.Transaction{cbtx}, []byte{}, -1) // prevHash empty
}

/*GetSpendableOutputs gets the utxos associated with an address up to what it needs to access in order to spend a given amount
@param address - the address in question
@param amount - the amount that the address is attempting to send
@return int - the total spendable amount in the utxos returned - in case this is greater than the amount attempting to be sent
@return map - the utxos
*/
func (bc *BlockChain) GetSpendableOutputs(address string, amount int) (int, map[string][]int) {
	utxns := bc.GetUnspentTransactions(address)
	utxos := make(map[string][]int)
	txoSum := 0

Work:
	for _, tx := range utxns {
		txID := hex.EncodeToString(tx.ID)
		for i, txo := range tx.Outputs {
			if txo.CanBeUnlocked(address) && txoSum < amount {
				txoSum += txo.Amount
				utxos[txID] = append(utxos[txID], i)

				if txoSum >= amount {
					break Work
				}
			}
		}
	}
	return txoSum, utxos
}

/*GetUnspentTransactions gets the transactions that contain utxos owned by an address
@params address - the address in question
@return an array of the utxns
*/
func (bc *BlockChain) GetUnspentTransactions(address string) []types.Transaction {
	var utxns []types.Transaction
	stxos := make(map[string][]int) // Tracking indices of spent txos in txns
	iter := bc.Iterator()

	// Go through each txo of each txn
	for {
		block := iter.Next()
		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for i, txo := range tx.Outputs {
				if stxos[txID] != nil {
					for _, stxosIdx := range stxos[txID] {
						if stxosIdx == i {
							continue Outputs // continue if this txo idx is already in the map for this txID
						}
					}
				}
				if txo.CanBeUnlocked(address) {
					exists := false
					txid := tx.ID
					for _, utxn := range utxns {
						if hex.EncodeToString(utxn.ID) == hex.EncodeToString(txid) {
							exists = true
						}
					}
					if !exists {
						utxns = append(utxns, *tx)
					}
				}
			}

			if !tx.IsCoinbase() {
				for _, txin := range tx.Inputs {
					if txin.CanUnlock(address) {
						txinID := hex.EncodeToString(txin.TxID)
						stxos[txinID] = append(stxos[txinID], txin.OutputIndex) // add the txo idx to the map if the address has a txin w/ reference
					}
				}
			}
		}
		if len(block.PrevHash) == 0 {
			break
		}
	}
	return utxns
}

/*CreateTransaction makes a new Transaction to be added to a Block
@param from - sender
@param to - recipient
@param amount - amount to send
@return the Transaction
*/
func (bc *BlockChain) CreateTransaction(from, to string, amount int) *types.Transaction {
	txoSum, utxos := bc.GetSpendableOutputs(from, amount)
	return types.CreateTransaction(from, to, amount, txoSum, utxos)
}
