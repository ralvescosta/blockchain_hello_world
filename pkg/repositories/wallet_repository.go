package repositories

import (
	"blockchain/pkg/interfaces"
	"blockchain/pkg/wallet"
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"
)

type wallletRepository struct {
	db *redis.Client
}

func (pst wallletRepository) GetAllWallets() (*wallet.Wallets, error) {
	ctx := context.Background()
	keys, err := pst.db.Keys(ctx, "*").Result()
	if err != nil {
		return nil, err
	}

	wallets := wallet.Wallets{}
	for _, key := range keys {
		s, err := pst.db.Get(ctx, key).Bytes()
		if err != nil {
			return nil, err
		}

		var model *WalletModel
		err = json.Unmarshal(s, &model)
		if err != nil {
			return nil, err
		}
		wallets.Wallets[key] = model.ToWallet()
	}

	return &wallets, nil
}

func (pst wallletRepository) InsertNewWallet(wlt wallet.Wallet) error {
	ctx := context.Background()
	addr, err := wlt.Address()
	if err != nil {
		return err
	}

	_, err = pst.db.Set(ctx, string(addr), ToWalletModel(wlt), 0).Result()

	return err
}

func NewWallletRepository(db *redis.Client) interfaces.IWalletRepository {
	return &wallletRepository{db}
}
