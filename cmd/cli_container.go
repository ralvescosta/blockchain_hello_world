package cmd

import (
	"github.com/go-redis/redis/v8"

	"blockchain/pkg/interfaces/cli"
	"blockchain/pkg/repositories"
)

type CliContainer struct {
	blockchainDbConnection *redis.Client
	blockchainRepository   *repositories.BlockchainRepository

	walletDbConnection *redis.Client
	walletRepository   *repositories.WallletRepository

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

	cli := cli.NewCommandLine(blockchainRepository)

	return CliContainer{
		blockchainDbConnection,
		blockchainRepository,

		walletDbConnection,
		walletRepository,

		cli,
	}
}
