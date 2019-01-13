package chainDB

/* Database interfacing */

import (
	"github.com/danitello/go-blockchain/core/util"
	"github.com/dgraph-io/badger"
)

const (
	/*Dir - path to block data */
	Dir = "./tmp/blocks"
	/*LastHashKey is the db key -> value is hash of most recent block in db */
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
	util.HandleErr(err)

	return db
}
