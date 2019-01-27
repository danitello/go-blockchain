package types

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"

	"github.com/danitello/go-blockchain/common/hexutil"
)

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
	difficulty int
}

/*InitBlock creates a new Block
@param data - the data to be contained in the Block
@param prevHash - the hash of the previous Block in the chain
@return a new Block
*/
func InitBlock(data string, prevHash []byte) *Block {
	newBlock := &Block{[]byte{}, []byte(data), prevHash, 0, 12}
	newBlock.runProof()
	return newBlock
}

/*runProof creates a new proof for the given Block, adding it's Hash and Nonce metadata */
func (b *Block) runProof() {
	target := new(big.Int).Lsh(big.NewInt(1), uint(256-b.difficulty)) // Left shift, 256 is number of bits in a hash
	var hash [32]byte
	var bigIntHash big.Int // Makes a difference

	// Block.Nonce was initalized to 0
	for b.Nonce < math.MaxInt64 {
		hash, bigIntHash = b.computeHash(true)

		// If the bigIntHash is less than the target, we have found the nonce
		if bigIntHash.Cmp(target) == -1 {
			b.Hash = hash[:]
			fmt.Println()
			break
		} else {
			b.Nonce++
		}
	}
	fmt.Println("New block signed")
}

/*ValidateProof confirms that a given Block has been signed correctly and thus is a valid Block in the BlockChain
using the Nonce that has been computed for it
@return whether the Block has been signed correctly or not
*/
func (b *Block) ValidateProof() bool {
	var bigIntHash big.Int
	_, bigIntHash = b.computeHash(false)

	target := new(big.Int).Lsh(big.NewInt(1), uint(256-b.difficulty))

	return bigIntHash.Cmp(target) == -1
}

/*computeHash calculates the Hash for the given Block
@param print - whether to print outputs
@return *[32]byte version of hash
@return *big.Int version of hash
*/
func (b *Block) computeHash(print bool) ([32]byte, big.Int) {
	var bigIntHash big.Int

	hash := sha256.Sum256(b.compileProofData())
	if print {
		fmt.Printf("\r%x", hash)
	}
	bigIntHash.SetBytes(hash[:])

	return hash, bigIntHash
}

/*compileProofData creates the comprehensive data slice that will be hashed during the POW
@return a [][]byte containing the final data
*/
func (b *Block) compileProofData() []byte {
	return bytes.Join([][]byte{b.PrevHash, b.Data, hexutil.ToHex(int64(b.Nonce)), hexutil.ToHex(int64(b.difficulty))}, []byte{})
}
