package consensus

/** Proof of Work implementation */

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"

	"github.com/danitello/go-blockchain/core/types"
	"github.com/danitello/go-blockchain/core/util"
)

/*InitProof creates a new proof for the given Block
@param Block - the Block for which a proof must be determined
*/
func InitProof(b *types.Block) {
	target := new(big.Int).Lsh(big.NewInt(1), uint(256-b.Difficulty)) // Left shift, 256 is number of bits in a hash
	var hash [32]byte
	var bigIntHash big.Int // Makes a difference

	// Block.Nonce was initalized to 0
	for b.Nonce < math.MaxInt64 {
		compiledProofData := compileProofData(b)
		hash = sha256.Sum256(compiledProofData)
		fmt.Printf("\r%x", hash)

		// If the bigIntHash is less than the target, we have found the nonce
		bigIntHash.SetBytes(hash[:])
		if bigIntHash.Cmp(target) == -1 {
			b.Hash = hash[:]
			break
		} else {
			b.Nonce++
		}
	}
	fmt.Println()
}

/*ValidateProof confirms that a given Block has been signed correctly
@param b - the Block in question
@return whether the Block has been signed correctly or not
*/
func ValidateProof(b *types.Block) bool {
	var bigIntHash big.Int
	target := new(big.Int).Lsh(big.NewInt(1), uint(256-b.Difficulty))

	compileProofData := compileProofData(b)
	hash := sha256.Sum256(compileProofData)
	bigIntHash.SetBytes(hash[:])

	return bigIntHash.Cmp(target) == -1
}

/*CompileProofData creates the comprehensive data slice that will be hashed during the POW
@param b - the Block in question
@return a [][]byte containing the final data
*/
func compileProofData(b *types.Block) []byte {
	return bytes.Join([][]byte{b.PrevHash, b.Data, util.ToHex(int64(b.Nonce)), util.ToHex(int64(b.Difficulty))}, []byte{})
}
