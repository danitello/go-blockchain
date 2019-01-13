package core

import (
	"fmt"

	"github.com/danitello/go-blockchain/chainDB"
	"github.com/danitello/go-blockchain/consensus"
	"github.com/danitello/go-blockchain/core/types"
	"github.com/danitello/go-blockchain/core/util"

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
	var lastHash []byte

	// Get hash of most recent Block in the chain
	err := bc.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(chainDB.LastHashKey))
		util.HandleErr(err)

		lastHash, err = item.Value()

		return err
	})
	util.HandleErr(err)

	// Create a new block, save the new Block to db, and save new most recent hash in db
	newBlock := types.InitBlock(data, lastHash)
	consensus.InitProof(newBlock)
	fmt.Println("New block signed")

	err = bc.DB.Update(func(txn *badger.Txn) error {
		err = txn.Set(newBlock.Hash, util.SerializeBlock(newBlock))
		util.HandleErr(err)
		err = txn.Set([]byte(chainDB.LastHashKey), newBlock.Hash)

		// Update chain
		bc.LastHash = newBlock.Hash

		return err
	})

	util.HandleErr(err)
}

/*InitBlockChain instantiates a new instance of a BlockChain
@return the current working BlockChain
*/
func InitBlockChain() *BlockChain {
	var lastHash []byte
	db := chainDB.InitDB()

	// If no data exists in the tmp/blocks directory, create a BlockChain
	// else get that BlockChain
	err := db.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get([]byte(chainDB.LastHashKey)); err == badger.ErrKeyNotFound {
			fmt.Println("No existing BlockChain found in", chainDB.Dir)
			genesisBlock := genesis()
			fmt.Println("Genesis block signed")

			// Save genesis block
			err = txn.Set([]byte(chainDB.LastHashKey), genesisBlock.Hash)
			util.HandleErr(err)
			err = txn.Set(genesisBlock.Hash, util.SerializeBlock(genesisBlock))

			lastHash = genesisBlock.Hash

			return err
		}

		item, err := txn.Get([]byte(chainDB.LastHashKey))
		util.HandleErr(err)

		lastHash, err = item.Value()

		return err

	})

	util.HandleErr(err)
	newChain := &BlockChain{lastHash, db}

	return newChain
}

/*genesis creates a genesis Block
@return the Block
*/
func genesis() *types.Block {
	genesisBlock := types.InitBlock(genesisData, []byte{}) // prevHash empty

	// Generate nonce and hash for the new block
	consensus.InitProof(genesisBlock)

	return genesisBlock
}
