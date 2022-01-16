package repository

import (
	"context"

	"github.com/go-redis/redis/v8"

	pkgBlock "blockchain/pkg/block"
)

type Repository struct {
	db *redis.Client
}

func (pst Repository) GetLastBlock() (*pkgBlock.Block, error) {
	lastBlockSerialized, err := pst.db.Get(context.Background(), "last_block").Result()
	if shouldReturnRedisError(err) {
		return nil, err
	}

	block, err := pkgBlock.Deserialize([]byte(lastBlockSerialized))

	return block, err
}

func (pst Repository) GetBlockByKey(key []byte) (*pkgBlock.Block, error) {
	blockSerialized, err := pst.db.Get(context.Background(), string(key)).Bytes()
	if shouldReturnRedisError(err) {
		return nil, err
	}
	if err != nil && err.Error() == redis.Nil.Error() && len(blockSerialized) == 0 {
		return nil, nil
	}

	block, err := pkgBlock.Deserialize(blockSerialized)

	return block, err
}

func (pst Repository) GetOrCreateFirstBlock(firstBlock *pkgBlock.Block) (*pkgBlock.Block, error) {
	lastBlockSerialized, err := pst.db.Get(context.Background(), "last_block").Bytes()
	if shouldReturnRedisError(err) {
		return nil, err
	}

	if len(lastBlockSerialized) > 0 {
		block, err := pkgBlock.Deserialize(lastBlockSerialized)
		return block, err
	}

	err = pst.txnCreateNewBlock(firstBlock)
	if err != nil {
		return nil, err
	}

	return firstBlock, err
}

func (pst Repository) InsertNewBlock(block *pkgBlock.Block) (*pkgBlock.Block, error) {
	err := pst.txnCreateNewBlock(block)

	return block, err
}

func (pst Repository) Dispose() {
	defer pst.db.Close()
}

func (pst Repository) txnCreateNewBlock(block *pkgBlock.Block) error {
	serialized, err := block.Serialize()
	if err != nil {
		return err
	}
	ctx := context.Background()
	err = pst.db.Watch(ctx, func(tx *redis.Tx) error {
		_, err := tx.Pipelined(ctx, func(pipe redis.Pipeliner) error {
			err = pipe.Set(ctx, string(block.Hash), serialized, 0).Err()
			if err != nil {
				return err
			}
			err = pipe.Set(ctx, "last_block", serialized, 0).Err()
			return err
		})
		return err
	})

	return err
}

func NewRepository(db *redis.Client) *Repository {
	return &Repository{db}
}

func shouldReturnRedisError(err error) bool {
	if err != nil && err.Error() != redis.Nil.Error() {
		return true
	}
	return false
}
