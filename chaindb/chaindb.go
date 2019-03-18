package chaindb

// Database interfacing

import (
	"log"

	"github.com/danitello/go-blockchain/common/byteutil"
	"github.com/danitello/go-blockchain/common/errutil"
	"github.com/danitello/go-blockchain/core/types"
	"github.com/dgraph-io/badger"
)

// ChainDB is the database for a BlockChain
type ChainDB struct {
	Database *badger.DB
}

const (
	// Dir - path to block data
	Dir = "./tmp/blocks"

	// LastHashKey is the db key -> value is hash of most recent block in db
	LastHashKey = "lastHashKey"
)

// InitDB instantiates a new ChainDB instance from the specified directory
func InitDB() *ChainDB {
	opts := badger.DefaultOptions
	opts.Dir = Dir
	opts.ValueDir = Dir
	bdb, err := badger.Open(opts)
	errutil.Handle(err)
	db := ChainDB{bdb}
	return &db
}

// HasChain determines whether the ChainDB instance has a previously initiated BlockChain
func (db *ChainDB) HasChain() bool {
	var exists bool
	err := db.Database.View(func(txn *badger.Txn) error {
		if _, err := txn.Get([]byte(LastHashKey)); err == badger.ErrKeyNotFound {
			exists = false
			return err
		}

		exists = true
		return nil
	})
	if err != nil {
		log.Println(err)
	}
	return exists
}

// ReadLastHash gets the hash of the most recent Block in the database
func (db *ChainDB) ReadLastHash() (lastHash []byte) {
	err := db.Database.View(func(txn *badger.Txn) (err error) {
		item, err := txn.Get([]byte(LastHashKey))
		errutil.Handle(err)

		lastHash, err = item.Value()
		return
	})
	errutil.Handle(err)

	return
}

// ReadBlockWithHash gets a Block from the database, given it's hash
func (db *ChainDB) ReadBlockWithHash(hash []byte) (resBlock *types.Block) {
	err := db.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(hash))
		errutil.Handle(err)

		value, err := item.Value()
		resBlock = types.DeserializeBlock(value)

		return err
	})
	errutil.Handle(err)

	return
}

// WriteNewLastBlock writes a new Block into the database and updates the last hash value
func (db *ChainDB) WriteNewLastBlock(newBlock *types.Block) {
	err := db.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, byteutil.Serialize(newBlock))
		errutil.Handle(err)

		err = txn.Set([]byte(LastHashKey), newBlock.Hash)
		return err
	})

	errutil.Handle(err)
}

// CloseDB closes the badgerdb
func (db *ChainDB) CloseDB() {
	db.Database.Close()
}
