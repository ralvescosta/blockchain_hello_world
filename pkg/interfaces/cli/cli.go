package cli

import (
	"log"
	"os"
	"runtime"
	"strconv"

	pkgBlock "blockchain/pkg/block"
	"blockchain/pkg/blockchain"
)

type CommandLine struct {
	blockchain *blockchain.BlockChain
}

//printUsage will display what options are availble to the user
func (cli *CommandLine) printUsageCommand() {
	log.Println("Usage: ")
	log.Println(" add -block <BLOCK_DATA> - add a block to the chain")
	log.Println(" print - prints the blocks in the chain")
}

//validateArgs ensures the cli was given valid input
func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsageCommand()
		runtime.Goexit()
	}
}

//addBlock allows users to add blocks to the chain via the cli
func (pst *CommandLine) addBlockCommand(data string) {
	err := pst.blockchain.Add(data)
	if err != nil {
		log.Panic(err)
	}
	log.Println("Added Block!")
}

//printChain will display the entire contents of the blockchain
func (pst *CommandLine) printChainCommand() {
	iterator := pst.blockchain.Iterator()

	for {
		block, err := iterator.Next()
		if err != nil {
			log.Panic(err)
		}

		log.Println(block.ToString())
		pow := pkgBlock.NewProofOfWork(block)
		_, isValid := pow.Validate()
		log.Printf("Is Pow Valid: %s\n", strconv.FormatBool(isValid))

		if len(block.PrevHash) == 0 {
			break
		}
	}
}

func (pst *CommandLine) stateMachine(args []string) {
	switch args[3] {
	case "add":
		addBlockData := args[5]
		if addBlockData == "" {
			pst.printUsageCommand()
			runtime.Goexit()
		}
		pst.addBlockCommand(addBlockData)
	case "print":
		pst.printChainCommand()
	default:
		pst.printUsageCommand()
	}
}

func (pst *CommandLine) Run() {
	pst.validateArgs()
	pst.stateMachine(os.Args)
}

func NewCommandLine(blockchain *blockchain.BlockChain) *CommandLine {
	return &CommandLine{blockchain}
}
