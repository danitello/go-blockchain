package core

import (
	"fmt"

	"github.com/danitello/go-blockchain/chaindb"
	"github.com/danitello/go-blockchain/chaindb/dbutil"
	"github.com/danitello/go-blockchain/common/errutil"
	"github.com/danitello/go-blockchain/core/types"

	"github.com/dgraph-io/badger"
)

const (
	/*genesisData is the data that will go into the genesis Block by default */
	genesisData = "Genesis"
)

/*BlockChain is a complete blockchain
@param LastHash - the hash of the most recent Block added to this BlockChain
@param DB - a badger database instance
*/
type BlockChain struct {
	LastHash []byte
	DB       *badger.DB
}

/*AddBlock adds a new Block to a given BlockChain
@param bc - the BlockChain in question
@param data - the data to be contained in the Block
*/
func (bc *BlockChain) AddBlock(data string) {
	// Get hash of most recent Block in the chain
	lastHash := bc.getLastHash()

	// Create a new block and save it
	newBlock := types.InitBlock(data, lastHash)
	bc.saveNewLastBlock(newBlock)
}

/*InitBlockChain instantiates a new instance of a BlockChain
@return the current working BlockChain
*/
func InitBlockChain() *BlockChain {
	var lastHash []byte
	db := chaindb.InitDB()

	// if chaindb.HasChain(db) {
	// 	// CHANGE THIS WHEN DB WRAPPED
	// 	err := db.View(func(txn *badger.Txn) error {
	// 		item, err := txn.Get([]byte(chaindb.LastHashKey))
	// 		errutil.HandleErr(err)

	// 		lastHash, err = item.Value()

	// 		return err
	// 	})
	// 	errutil.HandleErr(err)
	// } else {
	// 	fmt.Println("H222")
	// 	fmt.Println("No existing BlockChain found in", chaindb.Dir)
	// 	genesisBlock := createGenesisBlock()
	// 	fmt.Println("Genesis block signed")

	// 	err := db.Update(func(txn *badger.Txn) error {
	// 		err := txn.Set([]byte(chaindb.LastHashKey), genesisBlock.Hash)
	// 		lastHash = genesisBlock.Hash
	// 		errutil.HandleErr(err)
	// 		err = txn.Set(genesisBlock.Hash, dbutil.SerializeBlock(genesisBlock))
	// 		return err
	// 	})
	// 	errutil.HandleErr(err)
	// }

	// If no data exists in the tmp/blocks directory, create a BlockChain
	// else get that BlockChain
	err := db.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get([]byte(chaindb.LastHashKey)); err == badger.ErrKeyNotFound {
			fmt.Println("No existing BlockChain found in", chaindb.Dir)
			genesisBlock := createGenesisBlock()
			fmt.Println("Genesis block signed")

			// Save genesis block
			err = txn.Set([]byte(chaindb.LastHashKey), genesisBlock.Hash)
			errutil.HandleErr(err)
			err = txn.Set(genesisBlock.Hash, dbutil.SerializeBlock(genesisBlock))

			lastHash = genesisBlock.Hash

			return err
		}

		item, err := txn.Get([]byte(chaindb.LastHashKey))
		errutil.HandleErr(err)

		lastHash, err = item.Value()

		return err

	})

	errutil.HandleErr(err)
	newChain := &BlockChain{lastHash, db}

	return newChain

}

/*getLastHash gets the hash of the most recent Block in the chain
@return - the hash
*/
func (bc *BlockChain) getLastHash() []byte {
	var lastHash []byte
	err := bc.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(chaindb.LastHashKey))
		errutil.HandleErr(err)

		lastHash, err = item.Value()

		return err
	})
	errutil.HandleErr(err)

	return lastHash
}

/*saveNewLastBlock saves the new Block to db, and the new most recent hash in db */
func (bc *BlockChain) saveNewLastBlock(newBlock *types.Block) {
	err := bc.DB.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, dbutil.SerializeBlock(newBlock))
		errutil.HandleErr(err)
		err = txn.Set([]byte(chaindb.LastHashKey), newBlock.Hash)

		// Update chain
		bc.LastHash = newBlock.Hash

		return err
	})

	errutil.HandleErr(err)
}

/*createGenesisBlock creates a genesis Block
@return the Block
*/
func createGenesisBlock() *types.Block {
	return types.InitBlock(genesisData, []byte{}) // prevHash empty
}
