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
	genesisData = "Genesis"
)

// BlockChain is a complete blockchain
type BlockChain struct {
	Height   int
	LastHash []byte
	ChainDB  *chaindb.ChainDB
}

// InitBlockChain instantiates a new instance of a BlockChain
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

// GetBlockChain gets an existing BlockChain from the database
func GetBlockChain() *BlockChain {
	db := chaindb.InitDB()

	if !db.HasChain() {
		log.Panic("Error: No BlockChain exists")
	}
	resChain := &BlockChain{
		Height:   0,
		LastHash: []byte{0},
		ChainDB:  db}
	resChain.LastHash = db.ReadLastHash()
	resChain.Height = db.ReadBlockWithHash(resChain.LastHash).Index + 1

	return resChain
}

// AddBlock adds a new Block to a given BlockChain
func (bc *BlockChain) AddBlock(txns []*types.Transaction) {
	// Create a new block and save it
	newBlock := types.InitBlock(txns, bc.LastHash, bc.Height-1)
	bc.saveNewLastBlock(newBlock)
}

// saveNewLastBlock saves the new Block to db, and updates BlockChain struct
func (bc *BlockChain) saveNewLastBlock(newBlock *types.Block) {

	// Update DB
	bc.ChainDB.WriteNewLastBlock(newBlock)

	// Update chain
	bc.LastHash = newBlock.Hash
	//bc.UpdateUTXOSet(newBlock)
	bc.Reindex()

}

// createGenesisBlock creates the first Block
func createGenesisBlock(address string) *types.Block {
	cbtx := types.CoinbaseTx(address)
	return types.InitBlock([]*types.Transaction{cbtx}, []byte{}, -1) // prevHash empty
}

// GetUTXO gets the all the utxos in the chain
func (bc *BlockChain) GetUTXO() map[string]types.TxOutputs {
	UTXO := make(map[string]types.TxOutputs)
	spentTXO := make(map[string][]int)
	iter := bc.Iterator()

	for {
		block := iter.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

			// Txos in first block in question are all unspent
		Outputs:
			for outIdx, txo := range tx.Outputs {
				if spentTXO[txID] != nil {
					for _, spentOutIdx := range spentTXO[txID] {
						if spentOutIdx == outIdx {
							continue Outputs // continue if this txo idx is already in the map for this txID
						}
					}
				}
				txos := UTXO[txID]
				txos.Outputs = append(txos.Outputs, txo)
				UTXO[txID] = txos
			}

			if !tx.IsCoinbase() {
				for _, txin := range tx.Inputs {
					txID := hex.EncodeToString(txin.TxID)
					spentTXO[txID] = append(spentTXO[txID], txin.OutputIdx) // add the txo idx to the map if the pubKeyHash has a txin w/ reference
				}
			}
		}

		if len(block.PrevHash) == 0 {
			break
		}
	}

	return UTXO
}

// CreateTransaction makes a new Transaction to be added to a Block
func (bc *BlockChain) CreateTransaction(from, to string, amount int) *types.Transaction {
	// Get wallet info using address
	wallets, err := wallet.InitWallets()
	errutil.Handle(err)
	w := wallets.GetWallet(from)
	pubKeyHash := wallet.HashPubKey(w.PublicKey)

	utxos, txoSum := bc.GetUTXOWithPubKey(pubKeyHash, amount)
	newTx := types.CreateTransaction(from, to, pubKeyHash, amount, txoSum, utxos)
	bc.SignTransaction(newTx, w.PrivateKey)
	return newTx
}

// SignTransaction gathers necessary data and initiates the flow for signing a tx
func (bc *BlockChain) SignTransaction(tx *types.Transaction, privKey ecdsa.PrivateKey) {
	prevTxs := make(map[string]types.Transaction)

	for _, txin := range tx.Inputs {
		prevTx, err := bc.GetTransactionWithID(txin.TxID)
		errutil.Handle(err)
		prevTxs[hex.EncodeToString(prevTx.ID)] = prevTx
	}

	tx.Sign(privKey, prevTxs)
}

// VerifyTransaction gathers necessary data and initiates the flow for verifying a tx
func (bc *BlockChain) VerifyTransaction(tx *types.Transaction) bool {
	if tx.IsCoinbase() {
		return true
	}

	prevTxs := make(map[string]types.Transaction)

	for _, txin := range tx.Inputs {
		prevTx, err := bc.GetTransactionWithID(txin.TxID)
		errutil.Handle(err)
		prevTxs[hex.EncodeToString(prevTx.ID)] = prevTx
	}

	return tx.Verify(prevTxs)
}

// GetTransactionWithID searches the bc for a Transaction with a given ID
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
