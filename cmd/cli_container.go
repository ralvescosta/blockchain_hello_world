package cmd

import (
	"log"

	"github.com/go-redis/redis/v8"

	"blockchain/pkg/blockchain"
	"blockchain/pkg/interfaces/cli"
	"blockchain/pkg/repository"
)

type CliContainer struct {
	dbConnection *redis.Client
	repository   *repository.Repository
	chain        *blockchain.BlockChain
	cli          *cli.CommandLine
}

func (pst CliContainer) Close() {
	pst.repository.Dispose()
}

func NewCliContainer() CliContainer {

	dbConnection := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	repository := repository.NewRepository(dbConnection)

	chain, err := blockchain.NewBlockchain(repository)
	if err != nil {
		log.Panic(err)
	}

	cli := cli.NewCommandLine(chain)

	return CliContainer{
		dbConnection,
		repository,
		chain,
		cli,
	}
}
