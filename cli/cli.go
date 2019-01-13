package cli

import (
	"fmt"

	"github.com/danitello/go-blockchain/consensus"
	"github.com/danitello/go-blockchain/core"
)

/*PrintChain prints the chain from newest to oldest Block
@param bc - the BlockChain in question
*/
func PrintChain(bc *core.BlockChain) {
	iter := bc.Iterator()

	for {
		currBlock := iter.Next()

		//fmt.Printf("Prev Hash: %x\n", block.PrevHash)
		fmt.Printf("Data: %s\n", currBlock.Data)
		fmt.Printf("Hash: %x\n", currBlock.Hash)
		fmt.Println("Verified:", consensus.ValidateProof(currBlock))
		//fmt.Println(block.Nonce)
		fmt.Println()

		if len(currBlock.PrevHash) == 0 {
			break
		}
	}
}
