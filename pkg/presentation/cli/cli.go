package cli

import (
	"os"
	"runtime"
	"strconv"

	"blockchain/pkg/presentation/controllers"
)

type CommandLine struct {
	createBlockchainController controllers.CreateBlockchainController
	getBalanceController       controllers.GetBalanceController
	printChainController       controllers.PrintChainController
	printUsageController       controllers.PrintUsageController
	sendController             controllers.SendController
}

func (pst *CommandLine) Exit() {
	pst.printUsageController.Print()
	runtime.Goexit()
}

//validateArgs ensures the cli was given valid input
func (pst *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		pst.Exit()
	}
}

func (pst *CommandLine) stateMachine(args []string) {
	switch args[3] {
	case "getbalance":
		address := args[5]
		if address == "" {
			pst.Exit()
		}
		pst.getBalanceController.Get(address)
	case "createblockchain":
		address := args[5]
		if address == "" {
			pst.Exit()
		}
		pst.createBlockchainController.Create(address)
	case "printchain":
		pst.printUsageController.Print()
	case "send":
		from := args[5]
		to := args[7]
		amount, err := strconv.Atoi(args[9])
		if err != nil {
			pst.Exit()
		}
		pst.sendController.Send(from, to, amount)
	default:
		pst.printUsageController.Print()
	}
}

func (pst *CommandLine) Run() {
	pst.validateArgs()
	pst.stateMachine(os.Args)
}

func NewCommandLine(
	createBlockchainController controllers.CreateBlockchainController,
	getBalanceController controllers.GetBalanceController,
	printChainController controllers.PrintChainController,
	printUsageController controllers.PrintUsageController,
	sendController controllers.SendController,
) *CommandLine {
	return &CommandLine{
		createBlockchainController,
		getBalanceController,
		printChainController,
		printUsageController,
		sendController,
	}
}
