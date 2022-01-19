package cli

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"

	pkgBlock "blockchain/pkg/block"
	"blockchain/pkg/blockchain"
	"blockchain/pkg/repositories"
	txn "blockchain/pkg/transaction"
)

type CommandLine struct {
	repository *repositories.Repository
}

//printUsage will display what options are availble to the user
func (cli *CommandLine) printUsageCommand() {
	log.Println("Usage: ")
	log.Println(" [x] createblockchain -address <ADDRESS> -> creates a blockchain and rewards the mining fee")
	log.Println(" [x] send -from <FROM> -to <TO> -amount <AMOUNT> -> Send amount of coins from one address to another")
	log.Println(" [x] getbalance -address <ADDRESS> -> get balance for ADDRESS")
	log.Println(" [x] printchain -> Prints the blocks in the chain")
}

//validateArgs ensures the cli was given valid input
func (pst *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		pst.printUsageCommand()
		runtime.Goexit()
	}
}

func (pst *CommandLine) createBlockChainCommand(address string) {
	chain, err := blockchain.NewBlockchain(address, pst.repository)
	defer pst.repository.Dispose()
	if err != nil {
		log.Fatal(err)
	}
	if chain.Status.Already {
		log.Println("blockchain already created")
		return
	}
	log.Println("Finished creating chain")
}

func (pst *CommandLine) getBalanceCommand(address string) {
	chain, err := blockchain.NewBlockchain(address, pst.repository)
	defer pst.repository.Dispose()
	if err != nil {
		log.Fatal(err)
	}
	balance := 0
	unspentTrancationOutputs := chain.FindUnspentTrancationOutputs(address)

	for _, out := range unspentTrancationOutputs {
		balance += out.Value
	}

	fmt.Printf("Balance of %s: %d\n", address, balance)
}

func (pst *CommandLine) sendCommand(from, to string, amount int) {
	chain, err := blockchain.NewBlockchain(from, pst.repository)
	defer pst.repository.Dispose()
	if err != nil {
		log.Fatal(err)
	}

	acc, validOutputs := chain.FindSpendableTransactionOutputs(from, amount)
	tx, err := txn.NewTransaction(from, to, amount, acc, validOutputs)
	if err != nil {
		log.Fatal(err)
	}

	err = chain.Add([]*txn.Transaction{tx})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Success!")
}

//printChain will display the entire contents of the blockchain
func (pst *CommandLine) printChainCommand() {
	chain, err := blockchain.NewBlockchain("", pst.repository)
	defer pst.repository.Dispose()
	if err != nil {
		log.Fatal(err)
	}

	iterator := chain.Iterator()

	for {
		block, err := iterator.Next()
		if err != nil {
			log.Panic(err)
		}

		log.Println(block.ToString())
		pow := pkgBlock.NewProofOfWork(block)
		isValid, _ := pow.Validate()
		log.Printf("Is Pow Valid: %s\n", strconv.FormatBool(isValid))

		// if len(block.PrevHash) == 0 {
		// 	break
		// }
	}
}

func (pst *CommandLine) stateMachine(args []string) {
	switch args[3] {
	case "getbalance":
		address := args[5]
		if address == "" {
			pst.printUsageCommand()
			runtime.Goexit()
		}
		pst.getBalanceCommand(address)
	case "createblockchain":
		address := args[5]
		if address == "" {
			pst.printUsageCommand()
			runtime.Goexit()
		}
		pst.createBlockChainCommand(address)
	case "printchain":
		pst.printUsageCommand()
	case "send":
		log.Println("send")
	default:
		pst.printUsageCommand()
	}
}

func (pst *CommandLine) Run() {
	pst.validateArgs()
	pst.stateMachine(os.Args)
}

func NewCommandLine(repo *repositories.Repository) *CommandLine {
	return &CommandLine{repo}
}
