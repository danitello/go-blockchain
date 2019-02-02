package main

import (
	"os"

	"github.com/danitello/go-blockchain/cli"
)

func main() {
	defer os.Exit(0)
	cli.Run()

}
