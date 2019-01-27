package core

import (
	"fmt"

	"github.com/danitello/go-blockchain/chaindb"
	"github.com/danitello/go-blockchain/core/types"
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
	ChainDB  *chaindb.ChainDB
}

/*InitBlockChain instantiates a new instance of a BlockChain
@return the current working BlockChain
*/
func InitBlockChain() (newChain *BlockChain) {

	db := chaindb.InitDB()
	newChain = &BlockChain{[]byte{0}, db}

	// If a BlockChain can be found, use it, otherwise make a new one
	if db.HasChain() {
		newChain.LastHash = db.GetLastHash()
	} else {
		fmt.Println("No existing BlockChain found in", chaindb.Dir)
		genesisBlock := createGenesisBlock()
		fmt.Println("Genesis block signed")

		newChain.saveNewLastBlock(genesisBlock)
	}

	return

}

/*AddBlock adds a new Block to a given BlockChain
@param data - the data to be contained in the Block
*/
func (bc *BlockChain) AddBlock(data string) {
	// Get hash of most recent Block in the chain
	lastHash := bc.ChainDB.GetLastHash()

	// Create a new block and save it
	newBlock := types.InitBlock(data, lastHash)
	bc.saveNewLastBlock(newBlock)
}

/*saveNewLastBlock saves the new Block to db, and updates BlockChain struct
@param newBlock - the Block to save
*/
func (bc *BlockChain) saveNewLastBlock(newBlock *types.Block) {

	// Update DB
	bc.ChainDB.SaveNewLastBlock(newBlock)

	// Update chain
	bc.LastHash = newBlock.Hash

}

/*createGenesisBlock creates a genesis Block
@return the Block
*/
func createGenesisBlock() *types.Block {
	return types.InitBlock(genesisData, []byte{}) // prevHash empty
}
