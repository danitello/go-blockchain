package core

import (
	"github.com/danitello/go-blockchain/chaindb"
	"github.com/danitello/go-blockchain/core/types"
)

// BlockChainIterator reverse traverses a given BlockChain
type BlockChainIterator struct {
	currentHash []byte
	db          *chaindb.ChainDB
}

// Iterator creates a new BlockChainIterator for a BlockChain instance
func (bc *BlockChain) Iterator() *BlockChainIterator {
	return &BlockChainIterator{bc.LastHash, bc.ChainDB}
}

// Next retrievies the next (older) Block in the chain
func (iter *BlockChainIterator) Next() (resBlock *types.Block) {
	// Get the Block represented by the CurrentHash
	resBlock = iter.db.ReadBlockWithHash(iter.currentHash)

	// Update iterator
	iter.currentHash = resBlock.PrevHash

	return
}
