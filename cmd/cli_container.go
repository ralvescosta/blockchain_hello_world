package cmd

import (
	"github.com/go-redis/redis/v8"

	"blockchain/pkg/interfaces/cli"
	"blockchain/pkg/repositories"
)

type CliContainer struct {
	dbConnection         *redis.Client
	blockchainRepository *repositories.BlockchainRepository
	cli                  *cli.CommandLine
}

func (pst CliContainer) Close() {
	pst.blockchainRepository.Dispose()
}

func NewCliContainer() CliContainer {

	dbConnection := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	blockchainRepository := repositories.NewBlockchainRepository(dbConnection)

	cli := cli.NewCommandLine(blockchainRepository)

	return CliContainer{
		dbConnection,
		blockchainRepository,
		cli,
	}
}
