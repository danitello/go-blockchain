package cli

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"runtime"
	"strconv"

	"github.com/danitello/go-blockchain/wallet"

	"github.com/danitello/go-blockchain/common/errutil"

	"github.com/danitello/go-blockchain/core"
	"github.com/danitello/go-blockchain/core/types"
)

// Run starts the cli and processes the args
func Run() {
	// Check if there are args (first arg is the "main" subcommand)
	if len(os.Args) < 2 {
		printHelp()
		runtime.Goexit()
	}

	// Commands
	balanceCommand := flag.NewFlagSet("balance", flag.ExitOnError)
	createWalletCommand := flag.NewFlagSet("create-wallet", flag.ExitOnError)
	initChainCommand := flag.NewFlagSet("init-chain", flag.ExitOnError)
	helpCommand := flag.NewFlagSet("help", flag.ExitOnError)
	addressListCommand := flag.NewFlagSet("address-list", flag.ExitOnError)
	printCommand := flag.NewFlagSet("print-chain", flag.ExitOnError)
	reindexCommand := flag.NewFlagSet("reindex", flag.ExitOnError)
	sendCommand := flag.NewFlagSet("send", flag.ExitOnError)

	// Subcommands (pointers)
	balanceAddress := balanceCommand.String("address", "", "(Required) The address to get balance of.")
	initChainCommandAddress := initChainCommand.String("address", "", "(Required) The address to init the chain with.")
	sendCommandFrom := sendCommand.String("from", "", "(Required) The address to send from.")
	sendCommandTo := sendCommand.String("to", "", "(Required) The address to send to.")
	sendCommandAmount := sendCommand.String("amount", "", "(Required) The amount to send.")

	// Parse relevant commands
	switch os.Args[1] {
	case "balance":
		balanceCommand.Parse(os.Args[2:])
	case "create-wallet":
		createWalletCommand.Parse(os.Args[2:])
	case "help":
		helpCommand.Parse(os.Args[2:])
	case "init-chain":
		initChainCommand.Parse(os.Args[2:])
	case "address-list":
		addressListCommand.Parse(os.Args[2:])
	case "print-chain":
		printCommand.Parse(os.Args[2:])
	case "reindex":
		reindexCommand.Parse(os.Args[2:])
	case "send":
		sendCommand.Parse(os.Args[2:])
	default:
		printHelp()
		runtime.Goexit()
	}

	// Check for and evaluate used commands
	if balanceCommand.Parsed() {
		if *balanceAddress == "" {
			balanceCommand.Usage()
			runtime.Goexit()
		}

		getBalance(*balanceAddress)
	}

	if createWalletCommand.Parsed() {
		createWallet()
	}

	if helpCommand.Parsed() {
		printHelp()
	}

	if initChainCommand.Parsed() {
		if *initChainCommandAddress == "" {
			initChainCommand.Usage()
			fmt.Println()
			runtime.Goexit() // Give badgerdb time to garbage collect
		}

		initChain(*initChainCommandAddress)

	}

	if addressListCommand.Parsed() {
		addressList()
	}

	if printCommand.Parsed() {
		printChain()
	}

	if reindexCommand.Parsed() {
		reindex()
	}

	if sendCommand.Parsed() {
		// Make sure the required input was submitted
		if *sendCommandFrom == "" || *sendCommandTo == "" || *sendCommandAmount == "" {
			sendCommand.Usage()
			fmt.Println()
			runtime.Goexit()
		}

		amt, err := strconv.Atoi(*sendCommandAmount)
		errutil.Handle(err)
		send(*sendCommandFrom, *sendCommandTo, amt)
	}

}

// addressList iterates through current Wallets and prints each Wallet address
func addressList() {
	ws, _ := wallet.InitWallets()
	addresses := ws.GetAddresses()
	for _, address := range addresses {
		fmt.Println(address)
	}
}

// getBalance prints the balance of the given address
func getBalance(address string) {
	if !wallet.ValidateAddress(address) {
		log.Panic("Invalid address")
	}

	bc := core.GetBlockChain()
	defer bc.ChainDB.CloseDB()

	pubKeyHash := wallet.GetPubKeyHashFromAddress(address)

	_, balance := bc.GetUTXOWithPubKey(pubKeyHash, math.MaxInt32)

	fmt.Printf("Balance of %s: %d\n", address, balance)
}

// createWallet instantiates current Wallets and adds a new Wallet to it, then prints out the address
func createWallet() {
	ws, _ := wallet.InitWallets()
	fmt.Println(ws.CreateWallet())
	ws.SaveToFile()
}

// initChain initializes a new BlockChain with a given address
func initChain(address string) {
	if !wallet.ValidateAddress(address) {
		log.Panic("Invalid address")
	}
	bc := core.InitBlockChain(address)
	defer bc.ChainDB.CloseDB()
}

// printChain prints the chain from newest to oldest Block
func printChain() {
	bc := core.GetBlockChain()
	iter := bc.Iterator()

	for {
		currBlock := iter.Next()

		fmt.Printf("Block\t %d\n", currBlock.Index)
		fmt.Println("----------")
		fmt.Printf("Hash: %x\n", currBlock.Hash)
		fmt.Printf("Mined Date: %s\n", currBlock.TimeStamp)
		fmt.Println("Verified:", currBlock.ValidateProof())
		for _, tx := range currBlock.Transactions {
			fmt.Println(tx)
		}
		fmt.Println()

		// Reached the beginning of the chain
		if len(currBlock.PrevHash) == 0 {
			break
		}
	}
}

// printHelp prints the instructions for the cli
func printHelp() {
	fmt.Println("Usage: go run main.go <command>")
	fmt.Println()
	fmt.Println("where <command> is one of:")
	fmt.Println("\taddress-list, balance, create-wallet, help, init-chain, print-chain, reindex, send")
	fmt.Println()
	//fmt.Println("./main.go <command> h\t\tquick help on <command>")

}

// reindex reindexes UTXO set
func reindex() {
	bc := core.GetBlockChain()
	defer bc.ChainDB.CloseDB()
	bc.Reindex()

	count := bc.CountUTX()
	fmt.Printf("Reindex complete! There are %d transactions in the UTXO set.\n", count)
}

// send initiates the addition of a Transaction to the chain given a sender, reciever, and amount
func send(from, to string, amount int) {
	if !wallet.ValidateAddress(from) {
		log.Panic("Invalid from address")
	}
	if !wallet.ValidateAddress(to) {
		log.Panic("Invalid to address")
	}
	var txns []*types.Transaction
	bc := core.GetBlockChain()
	defer bc.ChainDB.CloseDB()
	txns = append(txns, types.CoinbaseTx(from), bc.CreateTransaction(from, to, amount))
	bc.AddBlock(txns)
}
