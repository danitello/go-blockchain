package cli

import (
	"encoding/hex"
	"flag"
	"fmt"
	"math"
	"os"
	"runtime"
	"strconv"

	"github.com/danitello/go-blockchain/common/errutil"

	"github.com/danitello/go-blockchain/core"
	"github.com/danitello/go-blockchain/core/types"
)

/*Run starts the cli and processes the args*/
func Run() {
	// Check if there are args (first arg is the "main" subcommand)
	if len(os.Args) < 2 {
		printHelp()
		runtime.Goexit()
	}

	// Commands
	initCommand := flag.NewFlagSet("init", flag.ExitOnError)
	balanceCommand := flag.NewFlagSet("balance", flag.ExitOnError)
	sendCommand := flag.NewFlagSet("send", flag.ExitOnError)
	helpCommand := flag.NewFlagSet("help", flag.ExitOnError)
	printCommand := flag.NewFlagSet("print", flag.ExitOnError)

	// Subcommands (pointers)
	initCommandAddress := initCommand.String("address", "", "(Required) The address to init the chain with.")
	balanceAddress := balanceCommand.String("address", "", "(Required) The address to get balance of.")
	sendCommandFrom := sendCommand.String("from", "", "(Required) The address to send from.")
	sendCommandTo := sendCommand.String("to", "", "(Required) The address to send to.")
	sendCommandAmount := sendCommand.String("amount", "", "(Required) The amount to send.")
	//sendCommandHelp := addCommand.String("help", "", "Quick help on 'add' command")
	//helpCommandHelp:= helpCommand.String("help", "", "Quick help on 'help' command")
	//printCommandHelp:= printCommand.String("help", "", "Quick help on 'print' command")

	// Parse relevant commands
	switch os.Args[1] {
	case "init":
		initCommand.Parse(os.Args[2:])
	case "balance":
		balanceCommand.Parse(os.Args[2:])
	case "send":
		sendCommand.Parse(os.Args[2:])
	case "help":
		helpCommand.Parse(os.Args[2:])
	case "print":
		printCommand.Parse(os.Args[2:])
	default:
		printHelp()
		runtime.Goexit()
	}

	// Check for and evaluate used commands
	if initCommand.Parsed() {
		if *initCommandAddress == "" {
			initCommand.Usage()
			fmt.Println()
			runtime.Goexit() // Give badgerdb time to garbage collect
		}

		initBlockChain(*initCommandAddress)

	}
	if sendCommand.Parsed() {
		// Make sure the required input was submitted
		if *sendCommandFrom == "" || *sendCommandTo == "" || *sendCommandAmount == "" {
			sendCommand.Usage()
			fmt.Println()
			runtime.Goexit()
		}

		amt, err := strconv.Atoi(*sendCommandAmount)
		errutil.HandleErr(err)
		sendTransaction(*sendCommandFrom, *sendCommandTo, amt)
	}

	if balanceCommand.Parsed() {
		if *balanceAddress == "" {
			balanceCommand.Usage()
			runtime.Goexit()
		}

		getBalance(*balanceAddress)
	}

	if helpCommand.Parsed() {
		printHelp()
	}

	if printCommand.Parsed() {
		printChain()
	}

}

/*initBlockChain initializes a new BlockChain */
func initBlockChain(address string) {
	bc := core.InitBlockChain(address)
	bc.ChainDB.CloseDB()
}

/*getBalance prints the balance of the given address
@param address - the address in question
*/
func getBalance(address string) {
	bc := core.GetBlockChain()
	defer bc.ChainDB.CloseDB()
	txoSum, _ := bc.GetSpendableOutputs(address, math.MaxInt64)

	fmt.Printf("Balance of %s: %d\n", address, txoSum)
}

/*sendTransaction initiates the addition of a Transaction to the chain
@param from - the sender
@param to - the recipient
@param amount - amount to send
*/
func sendTransaction(from, to string, amount int) {
	var txns []*types.Transaction
	bc := core.GetBlockChain()
	txns = append(txns, bc.CreateTransaction(from, to, amount))
	bc.AddBlock(txns)
}

/*printChain prints the chain from newest to oldest Block
 */
func printChain() {
	bc := core.GetBlockChain()
	iter := bc.Iterator()

	for {
		currBlock := iter.Next()

		fmt.Printf("Block\t %d\n", currBlock.Index)
		fmt.Println("----------")
		fmt.Printf("First TxID: %s\n", hex.EncodeToString(currBlock.Transactions[0].ID))
		fmt.Printf("Hash: %x\n", currBlock.Hash)
		fmt.Printf("Time: %s\n", currBlock.TimeStamp)
		fmt.Println("Verified:", currBlock.ValidateProof())
		fmt.Println()

		// Reached the beginning of the chain
		if len(currBlock.PrevHash) == 0 {
			break
		}
	}
}

/*printHelp prints the instructions for the cli */
func printHelp() {
	fmt.Println("Usage: go run main.go <command>")
	fmt.Println()
	fmt.Println("where <command> is one of:")
	fmt.Println("\tbalance, help, init, print, send")
	fmt.Println()
	//fmt.Println("./main.go <command> h\t\tquick help on <command>")

}
