package controllers

import (
	"blockchain/pkg/blockchain"
	"blockchain/pkg/interfaces"
	txn "blockchain/pkg/transaction"
	"log"
)

type SendController struct {
	blockchainRepository interfaces.IBlockchainRepository
}

func (pst SendController) Send(from, to string, amount int) {
	chain, err := blockchain.NewBlockchain(from, pst.blockchainRepository)
	defer pst.blockchainRepository.Dispose()
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

	log.Println("Success!")
}

func NewSendController(blockchainRepository interfaces.IBlockchainRepository) SendController {
	return SendController{blockchainRepository}
}
