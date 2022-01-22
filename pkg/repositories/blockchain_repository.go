package repositories

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"

	pkgBlock "blockchain/pkg/block"
	"blockchain/pkg/interfaces"
)

type blockchainRepository struct {
	db *redis.Client
}

func (pst blockchainRepository) GetLastBlock() (*pkgBlock.Block, error) {
	s, err := pst.db.Get(context.Background(), "last_block").Bytes()
	if shouldReturnRedisError(err) {
		return nil, err
	}

	if err != nil && err.Error() == redis.Nil.Error() && s == nil {
		return nil, nil
	}

	var model *BlockModel
	json.Unmarshal(s, &model)

	return model.ToBlock(), nil
}

func (pst blockchainRepository) GetBlockByKey(key []byte) (*pkgBlock.Block, error) {
	s, err := pst.db.Get(context.Background(), string(key)).Bytes()
	if shouldReturnRedisError(err) {
		return nil, err
	}
	if err != nil && err.Error() == redis.Nil.Error() && s == nil {
		return nil, nil
	}

	var model *BlockModel
	json.Unmarshal(s, &model)

	return model.ToBlock(), nil
}

func (pst blockchainRepository) GetOrCreateFirstBlock(firstBlock *pkgBlock.Block) (*pkgBlock.Block, error) {
	s, err := pst.db.Get(context.Background(), "last_block").Bytes()
	if shouldReturnRedisError(err) {
		return nil, err
	}

	var lastModel *BlockModel
	json.Unmarshal(s, &lastModel)
	if lastModel != nil {
		return lastModel.ToBlock(), nil
	}

	err = pst.txnCreateNewBlock(firstBlock)
	if err != nil {
		return nil, err
	}

	return firstBlock, err
}

func (pst blockchainRepository) InsertNewBlock(block *pkgBlock.Block) (*pkgBlock.Block, error) {
	err := pst.txnCreateNewBlock(block)

	return block, err
}

func (pst blockchainRepository) Dispose() {
	defer pst.db.Close()
}

func (pst blockchainRepository) txnCreateNewBlock(block *pkgBlock.Block) error {
	model, err := BlockToModel(block)
	if err != nil {
		return err
	}

	ctx := context.Background()
	err = pst.db.Watch(ctx, func(tx *redis.Tx) error {
		_, err := tx.Pipelined(ctx, func(pipe redis.Pipeliner) error {
			err = pipe.Set(ctx, fmt.Sprintf("%x", block.Hash), model, 0).Err()
			if err != nil {
				return err
			}
			err = pipe.Set(ctx, "last_block", model, 0).Err()
			return err
		})
		return err
	})

	return err
}

func NewBlockchainRepository(db *redis.Client) interfaces.IBlockchainRepository {
	return &blockchainRepository{db}
}

func shouldReturnRedisError(err error) bool {
	if err != nil && err.Error() != redis.Nil.Error() {
		return true
	}
	return false
}
