package cmd

import (
	"log"
	"os"

	"github.com/dgraph-io/badger"

	"blockchain/pkg/blockchain"
	"blockchain/pkg/interfaces/cli"
	"blockchain/pkg/repository"
)

type CliContainer struct {
	dbConnection *badger.DB
	repository   *repository.Repository
	chain        *blockchain.BlockChain
	cli          *cli.CommandLine
}

func (pst CliContainer) Close() {
	pst.repository.Dispose()
}

func NewCliContainer() CliContainer {

	opts := badger.DefaultOptions(os.Getenv("DB_PATH"))
	dbConnection, err := badger.Open(opts)
	if err != nil {
		log.Panic(err)
	}
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
