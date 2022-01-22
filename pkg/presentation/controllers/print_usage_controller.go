package controllers

import "log"

type PrintUsageController struct{}

func (PrintUsageController) Print() {
	log.Println("Usage: ")
	log.Println(" [x] createblockchain -address <ADDRESS> -> creates a blockchain and rewards the mining fee")
	log.Println(" [x] send -from <FROM> -to <TO> -amount <AMOUNT> -> Send amount of coins from one address to another")
	log.Println(" [x] getbalance -address <ADDRESS> -> get balance for ADDRESS")
	log.Println(" [x] printchain -> Prints the blocks in the chain")
	log.Println(" [x] createwallet - Creates a new wallet")
	log.Println(" [x] listaddresses - Lists the addresses in the wallet file")
}

func NewPrintUsageController() PrintUsageController {
	return PrintUsageController{}
}
