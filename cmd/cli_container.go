package cmd

import (
	"github.com/go-redis/redis/v8"

	"blockchain/pkg/interfaces/cli"
	"blockchain/pkg/repositories"
)

type CliContainer struct {
	dbConnection *redis.Client
	repository   *repositories.Repository
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
	repository := repositories.NewRepository(dbConnection)

	cli := cli.NewCommandLine(repository)

	return CliContainer{
		dbConnection,
		repository,
		cli,
	}
}
