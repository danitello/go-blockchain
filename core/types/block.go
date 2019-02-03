package types

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/danitello/go-blockchain/common/errutil"
	"github.com/danitello/go-blockchain/common/hexutil"
)

/*Block is a block in the blockchain
@param Index - index of this Block in the BlockChain
@param Nonce - integer that completes hash of Block for successful signing
@param Difficulty - determines the target value to sign the Block
@param Hash - the hash of this block
@param PrevHash - the hash of the previous Block
@param TimeStamp - the time this Blocks proof
@param Transactions - the transactions contained in this Block
*/
type Block struct {
	Index        int
	Nonce        int
	Difficulty   int
	Hash         []byte
	PrevHash     []byte
	TimeStamp    []byte
	Transactions []*Transaction
}

/*InitBlock initializes a new Block
@param txns - the Transactions to be contained in the Block
@param prevHash - the hash of the previous Block in the chain
@param prevIndex - the index of the previous Block in the chain
@return a new Block
*/
func InitBlock(txns []*Transaction, prevHash []byte, prevIndex int) *Block {
	newBlock := &Block{
		Index:        prevIndex + 1,
		Nonce:        0,
		Difficulty:   12,
		Hash:         []byte{},
		Transactions: txns,
		PrevHash:     prevHash,
		TimeStamp:    []byte(time.Now().String())}
	newBlock.runProof()
	return newBlock
}

/*runProof creates a new proof for the given Block, adding it's Hash and Nonce metadata */
func (b *Block) runProof() {
	target := new(big.Int).Lsh(big.NewInt(1), uint(256-b.Difficulty)) // Left shift, 256 is number of bits in a hash
	var hash [32]byte
	var bigIntHash big.Int

	// Block.Nonce was initalized to 0
	for b.Nonce < math.MaxInt64 {
		hash, bigIntHash = b.computeHash(true)

		// If the bigIntHash is less than the target, we have found the nonce
		if bigIntHash.Cmp(target) == -1 {
			b.Hash = hash[:]
			b.TimeStamp = []byte(time.Now().String())
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

	target := new(big.Int).Lsh(big.NewInt(1), uint(256-b.Difficulty))

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
	return bytes.Join([][]byte{b.PrevHash, b.hashTransactions(), hexutil.ToHex(int64(b.Nonce)), hexutil.ToHex(int64(b.Difficulty))}, []byte{})
}

/*hashTransactions creates a hashed representation of the Transactions in a Block
@return the hash
*/
func (b *Block) hashTransactions() []byte {
	var txHashes [][]byte
	var resHash [32]byte

	// Get hash of each tx
	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.ID)
	}

	// Get final hash
	resHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return resHash[:]
}

/*DeserializeBlock converts a []byte into a Block for database compatibility
@param data - the []byte representation of a Block
@returns a types.Block representation of a Block
*/
func DeserializeBlock(data []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)
	errutil.HandleErr(err)

	return &block
}
