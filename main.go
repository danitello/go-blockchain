package main

import (
	"fmt"
	"os"

	"github.com/danitello/go-blockchain/cli"
	"github.com/danitello/go-blockchain/core"
)

func main() {
	defer os.Exit(0)

	// Get BlockChain
	chain := core.InitBlockChain()
	defer chain.DB.Close()

	// Start CLI
	cl := cli.CommandLine{BC: chain}
	fmt.Println()
	cl.Run()

}
