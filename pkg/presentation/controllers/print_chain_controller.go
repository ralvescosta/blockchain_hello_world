package controllers

import (
	pkgBlock "blockchain/pkg/block"
	"blockchain/pkg/blockchain"
	"blockchain/pkg/interfaces"

	"log"
	"strconv"
)

type PrintChainController struct {
	blockchainRepository interfaces.IBlockchainRepository
}

func (pst PrintChainController) Print() {
	chain, err := blockchain.NewBlockchain("", pst.blockchainRepository)
	defer pst.blockchainRepository.Dispose()
	if err != nil {
		log.Fatal(err)
	}

	iterator := chain.Iterator()

	for {
		block, err := iterator.Next()
		if err != nil {
			log.Panic(err)
		}

		log.Println(block.ToString())
		pow := pkgBlock.NewProofOfWork(block)
		isValid, _ := pow.Validate()
		log.Printf("Is Pow Valid: %s\n", strconv.FormatBool(isValid))

		// if len(block.PrevHash) == 0 {
		// 	break
		// }
	}
}

func NewPrintChainController(blockchainRepository interfaces.IBlockchainRepository) PrintChainController {
	return PrintChainController{blockchainRepository}
}
