package main

import (
	"fmt"

	"github.com/danitello/go-blockchain/cli"
	"github.com/danitello/go-blockchain/core"
)

func main() {
	chain := core.InitBlockChain()
	//chain.AddBlock("More Block.")
	//chain.AddBlock("$1b")

	fmt.Println()
	//chain.Iterate()
	cli.PrintChain(chain)

}
