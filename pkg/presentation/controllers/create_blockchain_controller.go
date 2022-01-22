package controllers

import (
	"blockchain/pkg/blockchain"
	"blockchain/pkg/interfaces"
	"log"
)

type CreateBlockchainController struct {
	blockchainRepository interfaces.IBlockchainRepository
}

func (pst CreateBlockchainController) Create(address string) {
	chain, err := blockchain.NewBlockchain(address, pst.blockchainRepository)
	defer pst.blockchainRepository.Dispose()
	if err != nil {
		log.Fatal(err)
	}
	if chain.Status.Already {
		log.Println("blockchain already created")
		return
	}
	log.Println("Finished creating chain")
}

func NewCreateBlockchainController(blockchainRepository interfaces.IBlockchainRepository) CreateBlockchainController {
	return CreateBlockchainController{blockchainRepository}
}
