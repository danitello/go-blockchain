package chaindb

/* Database interfacing */

import (
	"fmt"

	"github.com/danitello/go-blockchain/common/errutil"
	"github.com/dgraph-io/badger"
)

const (
	// Dir - path to block data
	Dir = "./tmp/blocks"
	// LastHashKey is the db key -> value is hash of most recent block in db
	LastHashKey = "latesthash"
)

/*InitDB instantiates a new badger database instance from the specified directory
@return a pointer to the new db instance
*/
func InitDB() *badger.DB {
	opts := badger.DefaultOptions
	opts.Dir = Dir
	opts.ValueDir = Dir
	db, err := badger.Open(opts)
	errutil.HandleErr(err)

	return db
}

func HasChain(db *badger.DB) bool {
	var exists bool
	err := db.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get([]byte(LastHashKey)); err == badger.ErrKeyNotFound {
			fmt.Println("OKAY")
			exists = false
			return err
		}

		exists = true
		return nil
	})
	fmt.Println("HERE", exists)
	errutil.HandleErr(err)
	return exists
}
