package types

/*Block is a block in the blockchain
@param Hash - the hash of this block
@param Data - the data in this block
@param PrevHash - the hash of the previous Block
@param Nonce - integer that completes hash of Block for successful signing
@param Difficulty - determines the target value to sign the Block
*/
type Block struct {
	Hash       []byte
	Data       []byte
	PrevHash   []byte
	Nonce      int
	Difficulty int
}

/*InitBlock creates a new Block
@param data - the data to be contained in the Block
@param prevHash - the hash of the previous Block in the chain
@return a new Block
*/
func InitBlock(data string, prevHash []byte) *Block {
	newBlock := &Block{[]byte{}, []byte(data), prevHash, 0, 12}
	return newBlock
}
