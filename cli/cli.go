package cli

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/danitello/go-blockchain/consensus"
	"github.com/danitello/go-blockchain/core"
)

/*CommandLine allows the user to interact with the BlockChain through the terminal
@param BC - the chain in question
*/
type CommandLine struct {
	BC *core.BlockChain
}

/*Run starts the cli and retrieves the args sent with the "go run" command*/
func (cl *CommandLine) Run() {
	// Check if there are args (first arg is the main run command)
	if len(os.Args) < 2 {
		printHelp()
		runtime.Goexit()
	}

	// Commands
	addCommand := flag.NewFlagSet("add", flag.ExitOnError)
	helpCommand := flag.NewFlagSet("help", flag.ExitOnError)
	printCommand := flag.NewFlagSet("print", flag.ExitOnError)

	// Subcommands (pointers)
	addCommandData := addCommand.String("block", "", "(Required) The data to be put in the Block.")
	//addCommandHelp := addCommand.String("help", "", "Quick help on 'add' command")
	//helpCommandHelp:= helpCommand.String("help", "", "Quick help on 'help' command")
	//printCommandHelp:= printCommand.String("help", "", "Quick help on 'print' command")

	// Parse relevant commands
	switch os.Args[1] {
	case "add":
		addCommand.Parse(os.Args[2:])
	case "help":
		helpCommand.Parse(os.Args[2:])
	case "print":
		printCommand.Parse(os.Args[2:])
	default:
		printHelp()
		runtime.Goexit()
	}

	// Check for and evaluate used commands
	if addCommand.Parsed() {
		// 'add' requires a subcommand
		if len(addCommand.Args()) == 0 || addCommand.Arg(0) == "help" {
			addCommand.PrintDefaults()
			runtime.Goexit()
		} else if addCommand.Arg(0) == "block" { // want to add a block
			// Make sure the required input was submitted
			if *addCommandData == "" {
				addCommand.PrintDefaults()
				fmt.Println()
				runtime.Goexit() // Give badgerdb time to garbage collect
			}
			cl.addBlock(*addCommandData)
		}
	}

	if helpCommand.Parsed() {
		printHelp()
		runtime.Goexit()
	}

	if printCommand.Parsed() {
		cl.printChain()
	}

}

/*addBlock initiates the addition of a Block to the chain
@param data - the data to be contained in the Block
*/
func (cl *CommandLine) addBlock(data string) {
	cl.BC.AddBlock(data)
}

/*printChain prints the chain from newest to oldest Block
 */
func (cl *CommandLine) printChain() {
	iter := cl.BC.Iterator()

	for {
		currBlock := iter.Next()

		//fmt.Printf("Prev Hash: %x\n", block.PrevHash)
		fmt.Printf("Data: %s\n", currBlock.Data)
		fmt.Printf("Hash: %x\n", currBlock.Hash)
		fmt.Println("Verified:", consensus.ValidateProof(currBlock))
		//fmt.Println(block.Nonce)
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
	fmt.Println("\tadd, help, print")
	fmt.Println()
	//fmt.Println("./main.go <command> -h\t\tquick help on <command>")

}
