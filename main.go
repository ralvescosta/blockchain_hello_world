package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"

	"github.com/dgraph-io/badger"
)

type CommandLine struct {
	blockchain *BlockChain
}

//printUsage will display what options are availble to the user
func (cli *CommandLine) printUsage() {
	fmt.Println("Usage: ")
	fmt.Println(" add -block <BLOCK_DATA> - add a block to the chain")
	fmt.Println(" print - prints the blocks in the chain")
}

//validateArgs ensures the cli was given valid input
func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		//go exit will exit the application by shutting down the goroutine
		// if you were to use os.exit you might corrupt the data
		runtime.Goexit()
	}
}

//addBlock allows users to add blocks to the chain via the cli
func (cli *CommandLine) addBlock(data string) {
	cli.blockchain.Add(data)
	fmt.Println("Added Block!")
}

//printChain will display the entire contents of the blockchain
func (cli *CommandLine) printChain() {
	iterator := cli.blockchain.Iterator()

	for {
		block, err := iterator.Next()
		if err != nil {
			log.Panic(err)
		}

		log.Printf("Previous hash: %x\n", block.PrevHash)
		log.Printf("data: %s\n", block.Data)
		log.Printf("hash: %x\n", block.Hash)
		pow := NewProofOfWork(block)
		_, isValid := pow.Validate()
		log.Printf("Pow: %s\n", strconv.FormatBool(isValid))
		log.Println()
		// This works because the Genesis block has no PrevHash to point to.
		if len(block.PrevHash) == 0 {
			break
		}
	}
}

//run will start up the command line
func (cli *CommandLine) run() {
	cli.validateArgs()

	addBlockCmd := flag.NewFlagSet("add", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("print", flag.ExitOnError)
	addBlockData := addBlockCmd.String("block", "", "Block data")

	switch os.Args[1] {
	case "add":
		err := addBlockCmd.Parse(os.Args[2:])
		log.Panic(err)

	case "print":
		err := printChainCmd.Parse(os.Args[2:])
		log.Panic(err)

	default:
		cli.printUsage()
		runtime.Goexit()
	}
	// Parsed() will return true if the object it was used on has been called
	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			runtime.Goexit()
		}
		cli.addBlock(*addBlockData)
	}
	if printChainCmd.Parsed() {
		cli.printChain()
	}
}

func main() {
	defer os.Exit(0)

	opts := badger.DefaultOptions(dbPath)
	db, err := badger.Open(opts)
	repository := NewRepository(db)
	defer repository.Dispose()

	chain, err := NewBlockchain(&repository)
	if err != nil {
		log.Panic(err)
	}

	cli := CommandLine{chain}

	cli.run()
}
