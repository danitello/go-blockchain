package types

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/danitello/go-blockchain/common/byteutil"
	"github.com/danitello/go-blockchain/common/errutil"
	"github.com/danitello/go-blockchain/common/hexutil"
)

// Block is a block in the blockchain with
// Index - index of this Block in the BlockChain
// Nonce - integer that completes hash of Block for successful signing
// Difficulty - determines the target value to sign the Block
// Hash - the hash of this block
// PrevHash - the hash of the previous Block
// TimeStamp - the time this Blocks proof
// Transactions - the transactions contained in this Block
type Block struct {
	Index        int
	Nonce        int
	Difficulty   int
	Hash         []byte
	PrevHash     []byte
	TimeStamp    []byte
	Transactions []*Transaction
}

// InitBlock initializes a new Block
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

// runProof creates a new proof for the given Block, adding it's Hash and Nonce metadata
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

// ValidateProof confirms that a given Block has been signed correctly and thus is a valid Block in the BlockChain
// using the Nonce that has been computed for it
func (b *Block) ValidateProof() bool {
	var bigIntHash big.Int
	_, bigIntHash = b.computeHash(false)

	target := new(big.Int).Lsh(big.NewInt(1), uint(256-b.Difficulty))

	return bigIntHash.Cmp(target) == -1
}

// computeHash calculates the Hash for the given Block
func (b *Block) computeHash(print bool) ([32]byte, big.Int) {
	var bigIntHash big.Int

	hash := sha256.Sum256(b.compileProofData())
	if print {
		fmt.Printf("\rBlock Hash: %x", hash)
	}
	bigIntHash.SetBytes(hash[:])

	return hash, bigIntHash
}

// compileProofData creates the comprehensive data slice that will be hashed during the POW
func (b *Block) compileProofData() []byte {
	return bytes.Join([][]byte{b.PrevHash, b.getMerkleTree(), hexutil.ToHex(int64(b.Nonce)), hexutil.ToHex(int64(b.Difficulty))}, []byte{})
}

// getMerkleTree gets the MerkleTree representation of the Transactions in the Block and returns the root
func (b *Block) getMerkleTree() []byte {
	var txs [][]byte

	// Get txs
	for _, tx := range b.Transactions {
		txs = append(txs, byteutil.Serialize(tx))
	}

	// Create MerkleTree
	tree := InitMerkleTree(txs)

	return tree.Root.Data
}

// DeserializeBlock converts a []byte into a Block for database compatibility
func DeserializeBlock(data []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)
	errutil.Handle(err)

	return &block
}
