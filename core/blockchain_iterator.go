package core

import (
	"github.com/danitello/go-blockchain/chaindb"
	"github.com/danitello/go-blockchain/core/types"
)

/*BlockChainIterator reverse traverses a given BlockChain
@param currentHash - the hash of the Block that the Iterator is currently on in the chain
@param db - the badger database associated with the chain
*/
type BlockChainIterator struct {
	currentHash []byte
	db          *chaindb.ChainDB
}

/*Iterator creates a new BlockChainIterator for a BlockChain instance
@return a new BlockChainIterator
*/
func (bc *BlockChain) Iterator() *BlockChainIterator {
	return &BlockChainIterator{bc.LastHash, bc.ChainDB}
}

/*Next retrievies the next (older) Block in the chain
@param iter - the current Iterator
@return the next Block
*/
func (iter *BlockChainIterator) Next() (resBlock *types.Block) {
	// Get the Block represented by the CurrentHash
	resBlock = iter.db.ReadBlockWithHash(iter.currentHash)

	// Update iterator
	iter.currentHash = resBlock.PrevHash

	return
}
