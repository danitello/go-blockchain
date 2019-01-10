package main

import (
	"fmt"

	"github.com/danitello/go-blockchain/consensus"

	"github.com/danitello/go-blockchain/core"
)

func main() {
	chain := core.InitBlockChain("This is just the beginning.")
	chain.AddBlock("More Block.")
	chain.AddBlock("$1b")

	fmt.Println()
	for _, block := range chain.Blocks {
		//fmt.Printf("Prev Hash: %x\n", block.PrevHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Println("Verified:", consensus.ValidateProof(block))
		fmt.Println(block.Nonce)
	}

}
