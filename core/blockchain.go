package core

import (
	"encoding/hex"
	"fmt"

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
func InitBlockChain(address string) (resChain *BlockChain) {

	db := chaindb.InitDB()
	resChain = &BlockChain{
		Height:   0,
		LastHash: []byte{0},
		ChainDB:  db}

	// If a BlockChain can be found, use it, otherwise make a new one
	if db.HasChain() {
		resChain.LastHash = db.GetLastHash()
		resChain.Height = db.GetBlockWithHash(resChain.LastHash).Index + 1
	} else {
		fmt.Println("No existing BlockChain found in", chaindb.Dir)
		genesisBlock := createGenesisBlock(address)
		fmt.Println("Genesis block signed")

		resChain.saveNewLastBlock(genesisBlock)
	}

	return

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

/*GetSpendableOutputs gets the utxos associated with an address that it needs to access to spend a given amount
@param address - the address in question
@param amount - the amount that the address is attempting to send
@return int - the spendable amount
@return map - the utxos
*/
func (bc *BlockChain) GetSpendableOutputs(address string, amount int) (int, map[string][]int) {
	utxns := bc.GetUnspentTransactions(address)
	utxos := make(map[string][]int)
	sum := 0

Work:
	for _, tx := range utxns {
		txID := hex.EncodeToString(tx.ID)
		for i, txo := range tx.Outputs {
			if txo.CanBeUnlocked(address) && sum < amount {
				sum += txo.Amount
				utxos[txID] = append(utxos[txID], i)

				if sum >= amount {
					break Work
				}
			}
		}
	}
	return sum, utxos
}

// func (bc *Blockchain) GetUTXO(address string) []types.TxOutput {

// }

/*GetUnspentTransactions gets the transactions that contain utxos owned by an address
@params address - the address in question
@return an array of the utxos
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
					utxns = append(utxns, *tx)
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
