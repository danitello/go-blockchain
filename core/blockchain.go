package core

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"log"

	"github.com/danitello/go-blockchain/common/errutil"
	"github.com/danitello/go-blockchain/wallet"

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

/*GetUTXOs gets the utxos that are owned by pubKeyHash
@params pubKeyHash - the pubKeyHash in question
@return a map of txID -> utxoIdx
*/
func (bc *BlockChain) GetUTXOs(pubKeyHash []byte) map[string][]int {
	UTXOs := make(map[string][]int)
	spentTXOs := make(map[string][]int)
	iter := bc.Iterator()

	for {
		block := iter.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

			// Txos in first block in question are all unspent
		Outputs:
			for outIdx, txo := range tx.Outputs {
				if spentTXOs[txID] != nil {
					for _, spentOutIdx := range spentTXOs[txID] {
						if spentOutIdx == outIdx {
							continue Outputs // continue if this txo idx is already in the map for this txID
						}
					}
				}
				if txo.IsLockedWithKey(pubKeyHash) {
					txoIdxs := UTXOs[txID]
					txoIdxs = append(txoIdxs, outIdx)
					UTXOs[txID] = txoIdxs
				}

			}

			if !tx.IsCoinbase() {
				for _, txin := range tx.Inputs {
					txID := hex.EncodeToString(txin.TxID)
					spentTXOs[txID] = append(spentTXOs[txID], txin.OutputIdx) // add the txo idx to the map if the pubKeyHash has a txin w/ reference
				}
			}
		}

		if len(block.PrevHash) == 0 {
			break
		}
	}

	return UTXOs
}

/*GetSpendableOutputs gets utxos owned by a pub key hash with a total balance up to a given amount
@param pubKeyHash - the pub key hash in question
@param max - the max value to compute to
@return map - the utxos
@return int - the balance
*/
func (bc *BlockChain) GetSpendableOutputs(pubKeyHash []byte, max int) (map[string][]int, int) {
	UTXOs := bc.GetUTXOs(pubKeyHash)
	spendableUTXOs := make(map[string][]int)
	balance := 0

	for txID, txoIdxs := range UTXOs {
		for _, txoIdx := range txoIdxs {
			id, _ := hex.DecodeString(txID)
			tx, err := bc.GetTransactionWithID(id)
			errutil.HandleErr(err)
			for i, txo := range tx.Outputs {
				if i == txoIdx {
					sputxos := spendableUTXOs[txID]
					sputxos = append(sputxos, txoIdx)
					spendableUTXOs[txID] = sputxos
					balance += txo.Amount

					if balance > max {
						break
					}
				}
			}
		}
	}
	return spendableUTXOs, balance
}

/*CreateTransaction makes a new Transaction to be added to a Block
@param from - sender
@param to - recipient
@param amount - amount to send
@return the Transaction
*/
func (bc *BlockChain) CreateTransaction(from, to string, amount int) *types.Transaction {
	// Get wallet info using address
	wallets, err := wallet.InitWallets()
	errutil.HandleErr(err)
	w := wallets.GetWallet(from)
	pubKeyHash := wallet.HashPubKey(w.PublicKey)

	utxos, txoSum := bc.GetSpendableOutputs(pubKeyHash, amount)
	newTx := types.CreateTransaction(from, to, pubKeyHash, amount, txoSum, utxos)
	bc.SignTransaction(newTx, w.PrivateKey)
	return newTx
}

/*SignTransaction gathers necessary data and initiates the flow for signing a tx
@param tx - the tx to sign
@param privKey - of the signer
*/
func (bc *BlockChain) SignTransaction(tx *types.Transaction, privKey ecdsa.PrivateKey) {
	prevTxs := make(map[string]types.Transaction)

	for _, txin := range tx.Inputs {
		prevTx, err := bc.GetTransactionWithID(txin.TxID)
		errutil.HandleErr(err)
		prevTxs[hex.EncodeToString(prevTx.ID)] = prevTx
	}

	tx.Sign(privKey, prevTxs)
}

/*VerifyTransaction gathers necessary data and initiates the flow for verifying a tx
@param tx - the tx
@return whether it is valid
*/
func (bc *BlockChain) VerifyTransaction(tx *types.Transaction) bool {
	if tx.IsCoinbase() {
		return true
	}

	prevTxs := make(map[string]types.Transaction)

	for _, txin := range tx.Inputs {
		prevTx, err := bc.GetTransactionWithID(txin.TxID)
		errutil.HandleErr(err)
		prevTxs[hex.EncodeToString(prevTx.ID)] = prevTx
	}

	return tx.Verify(prevTxs)
}

/*GetTransactionWithID searches the bc for a Transaction with a given ID
@param id - the id
@return Transaction - the Transaction
@return error - any error
*/
func (bc *BlockChain) GetTransactionWithID(id []byte) (types.Transaction, error) {
	iter := bc.Iterator()

	for {
		block := iter.Next()

		for _, tx := range block.Transactions {
			if bytes.Compare(tx.ID, id) == 0 {
				return *tx, nil
			}
		}

		if len(block.PrevHash) == 0 {
			break
		}
	}

	return types.Transaction{}, errors.New("Transaction not found")
}
