package controllers

import (
	"blockchain/pkg/interfaces"
	"fmt"
	"log"
)

type ListAddressesController struct {
	walletRepository interfaces.IWalletRepository
}

func (pst ListAddressesController) List() {
	wallets, err := pst.walletRepository.GetAllWallets()
	if err != nil {
		log.Fatal(err)
	}

	addresses := wallets.GetAllAddresses()

	for _, address := range addresses {
		fmt.Println(address)
	}
}

func NewListAddressesController(walletRepository interfaces.IWalletRepository) ListAddressesController {
	return ListAddressesController{walletRepository}
}
