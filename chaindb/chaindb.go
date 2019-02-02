package chaindb

/* Database interfacing */

import (
	"github.com/danitello/go-blockchain/chaindb/dbutil"
	"github.com/danitello/go-blockchain/common/errutil"
	"github.com/danitello/go-blockchain/core/types"
	"github.com/dgraph-io/badger"
)

/*ChainDB is the database for a BlockChain
@param database - a badger db instance
*/
type ChainDB struct {
	database *badger.DB
}

const (
	// Dir - path to block data
	Dir = "./tmp/blocks"

	// LastHashKey is the db key -> value is hash of most recent block in db
	LastHashKey = "lastHashKey"
)

/*InitDB instantiates a new ChainDB instance from the specified directory
@return a pointer to the new db instance
*/
func InitDB() *ChainDB {
	opts := badger.DefaultOptions
	opts.Dir = Dir
	opts.ValueDir = Dir
	bdb, err := badger.Open(opts)
	errutil.HandleErr(err)
	db := ChainDB{bdb}
	return &db
}

/*HasChain determines whether the ChainDB instance has a previously initiated BlockChain
@return whether it does or not
*/
func (db *ChainDB) HasChain() bool {
	var exists bool
	err := db.database.View(func(txn *badger.Txn) error {
		if _, err := txn.Get([]byte(LastHashKey)); err == badger.ErrKeyNotFound {
			exists = false
			return err
		}

		exists = true
		return nil
	})
	errutil.HandleErr(err)
	return exists
}

/*GetLastHash gets the hash of the most recent Block in the database
@return - the hash
*/
func (db *ChainDB) GetLastHash() (lastHash []byte) {
	err := db.database.View(func(txn *badger.Txn) (err error) {
		item, err := txn.Get([]byte(LastHashKey))
		errutil.HandleErr(err)

		lastHash, err = item.Value()
		return
	})
	errutil.HandleErr(err)

	return
}

/*GetBlockWithHash gets a Block from the database, given it's hash
@param hash - the hash of the desired Block
@return - the Block
*/
func (db *ChainDB) GetBlockWithHash(hash []byte) (resBlock *types.Block) {
	err := db.database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(hash))
		errutil.HandleErr(err)

		value, err := item.Value()
		resBlock = types.DeserializeBlock(value)

		return err
	})
	errutil.HandleErr(err)

	return
}

/*SaveNewLastBlock saves a new Block into the database and updates the last hash value
@param newBlock - the Block
*/
func (db *ChainDB) SaveNewLastBlock(newBlock *types.Block) {
	err := db.database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, dbutil.Serialize(newBlock))
		errutil.HandleErr(err)

		err = txn.Set([]byte(LastHashKey), newBlock.Hash)
		return err
	})

	errutil.HandleErr(err)
}

/*CloseDB closes the badgerdb */
func (db *ChainDB) CloseDB() {
	db.database.Close()
}
