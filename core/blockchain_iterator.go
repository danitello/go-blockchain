package core

import (
	"fmt"

	"github.com/danitello/go-blockchain/core/types"
	"github.com/danitello/go-blockchain/core/util"
	"github.com/dgraph-io/badger"
)

/*BlockChainIterator reverse traverses a given BlockChain
@param currentHash - the hash of the Block that the Iterator is currently on in the chain
@param db - the badger database associated with the chain
*/
type BlockChainIterator struct {
	currentHash []byte
	db          *badger.DB
}

/*Iterator creates a new BlockChainIterator for a BlockChain instance
@param bc - the BlockChain in question
@return a new BlockChainIterator
*/
func (bc *BlockChain) Iterator() *BlockChainIterator {
	return &BlockChainIterator{bc.LastHash, bc.DB}
}

/*Next retrievies the next (older) Block in the chain
@param iter - the current Iterator
@return the next Block
*/
func (iter *BlockChainIterator) Next() (resBlock *types.Block) {
	resBlock = nil
	// Get the Block represented by the CurrentHash
	err := iter.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(iter.currentHash))
		util.HandleErr(err)

		value, err := item.Value()
		resBlock = util.DeserializeBlock(value)

		return err
	})
	util.HandleErr(err)

	// Update iterator
	iter.currentHash = resBlock.PrevHash

	return
}

/*Iterate and print the chain using badgerDB built in iterator -wip
@param bc - the BlockChain in question
@return error encountered
*/
func (bc *BlockChain) Iterate() error {
	err := bc.DB.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			//k := item.Key()
			value, err := item.Value()
			fmt.Printf("%x %s", util.DeserializeBlock(value).Hash, util.DeserializeBlock(value).Data)
			fmt.Println()
			util.HandleErr(err)
		}
		return nil
	})

	return err
}
