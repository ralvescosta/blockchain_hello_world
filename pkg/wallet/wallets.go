package wallet

type Wallets struct {
	Wallets map[string]*Wallet
}

func (ws *Wallets) Add() (string, error) {
	wallet, err := NewWallet()
	if err != nil {
		return "", err
	}

	addressInByte, err := wallet.Address()
	if err != nil {
		return "", err
	}

	address := string(addressInByte)

	ws.Wallets[address] = wallet

	return address, nil
}

func (pst Wallets) GetWallet(address string) Wallet {
	return *pst.Wallets[address]
}

func (pst Wallets) GetAllAddresses() []string {
	var addresses []string

	for address := range pst.Wallets {
		addresses = append(addresses, address)
	}
	return addresses
}

func NewWallets() (*Wallets, error) {
	wallets := Wallets{
		Wallets: make(map[string]*Wallet),
	}

	return &wallets, nil
}
