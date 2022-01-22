package cmd

import (
	"github.com/go-redis/redis/v8"

	"blockchain/pkg/interfaces"
	"blockchain/pkg/presentation/cli"
	"blockchain/pkg/presentation/controllers"
	"blockchain/pkg/repositories"
)

type CliContainer struct {
	blockchainDbConnection *redis.Client
	blockchainRepository   interfaces.IBlockchainRepository

	walletDbConnection *redis.Client
	walletRepository   interfaces.IWalletRepository

	cli *cli.CommandLine
}

func (pst CliContainer) Close() {
	pst.blockchainRepository.Dispose()
}

func NewCliContainer() CliContainer {

	blockchainDbConnection := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	blockchainRepository := repositories.NewBlockchainRepository(blockchainDbConnection)

	walletDbConnection := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       1,
	})
	walletRepository := repositories.NewWallletRepository(walletDbConnection)

	createBlockchainController := controllers.NewCreateBlockchainController(blockchainRepository)
	getBalanceController := controllers.NewGetBalanceController(blockchainRepository)
	printChainController := controllers.NewPrintChainController(blockchainRepository)
	printUsageController := controllers.NewPrintUsageController()
	sendControler := controllers.NewSendController(blockchainRepository)
	cli := cli.NewCommandLine(
		createBlockchainController,
		getBalanceController,
		printChainController,
		printUsageController,
		sendControler,
	)

	return CliContainer{
		blockchainDbConnection,
		blockchainRepository,

		walletDbConnection,
		walletRepository,

		cli,
	}
}
