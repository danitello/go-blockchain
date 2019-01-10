package core

import (
	"github.com/danitello/go-blockchain/consensus"
	"github.com/danitello/go-blockchain/core/types"
)

/*BlockChain is a complete blockchain
@param Blocks - the constituent blocks in the chain
*/
type BlockChain struct {
	Blocks []*types.Block
}

/*AddBlock adds a new Block to a given BlockChain
@param data - the data to be contained in the Block
*/
func (bc *BlockChain) AddBlock(data string) {
	prevHash := bc.Blocks[len(bc.Blocks)-1].Hash // Hash of last Block in bc
	newBlock := types.InitBlock(data, prevHash)

	// Generate nonce and hash
	consensus.InitProof(newBlock)

	bc.Blocks = append(bc.Blocks, newBlock)
}

/*InitBlockChain creates a new BlockChain
@param genesisData - the genesis Block's data
@return a newly minted BlockChain
*/
func InitBlockChain(genesisData string) *BlockChain {
	genesisBlock := types.InitBlock(genesisData, []byte{}) // prevHash empty

	// Generate nonce and hash
	consensus.InitProof(genesisBlock)

	return &BlockChain{[]*types.Block{genesisBlock}}
}
