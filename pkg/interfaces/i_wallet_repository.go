package interfaces

import "blockchain/pkg/wallet"

type IWalletRepository interface {
	GetAllWallets() (*wallet.Wallets, error)
	InsertNewWallet(wlt wallet.Wallet) error
}
