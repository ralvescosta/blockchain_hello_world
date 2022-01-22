package controllers

import (
	"blockchain/pkg/interfaces"
	"fmt"
	"log"
)

type CreateWalletController struct {
	walletRepository interfaces.IWalletRepository
}

func (pst CreateWalletController) Create() {
	wallets, err := pst.walletRepository.GetAllWallets()
	if err != nil {
		log.Fatal(err)
	}

	address, err := wallets.Add()
	if err != nil {
		log.Fatal(err)
	}

	err = pst.walletRepository.InsertNewWallet(wallets.GetWallet(address))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("New address is: %s\n", address)
}

func NewCreateWalletController(walletRepository interfaces.IWalletRepository) CreateWalletController {
	return CreateWalletController{walletRepository}
}
