package controllers

import (
	"blockchain/pkg/blockchain"
	"blockchain/pkg/interfaces"
	"fmt"
	"log"
)

type GetBalanceController struct {
	blockchainRepository interfaces.IBlockchainRepository
}

func (pst GetBalanceController) Get(address string) {
	chain, err := blockchain.NewBlockchain(address, pst.blockchainRepository)
	defer pst.blockchainRepository.Dispose()
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

func NewGetBalanceController(blockchainRepository interfaces.IBlockchainRepository) GetBalanceController {
	return GetBalanceController{blockchainRepository}
}
