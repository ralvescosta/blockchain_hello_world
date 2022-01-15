package repository

import (
	"github.com/dgraph-io/badger"

	pkgBlock "blockchain/pkg/block"
)

type Repository struct {
	db *badger.DB
}

func (pst Repository) GetLastBlock() (*pkgBlock.Block, error) {
	var block *pkgBlock.Block

	err := pst.db.View(func(txn *badger.Txn) error {
		lastHashSaved, err := txn.Get([]byte("lh"))
		if err != nil {
			return err
		}

		err = lastHashSaved.Value(func(lastHash []byte) error {
			lastBlock, err := txn.Get(lastHash)
			if err != nil {
				return err
			}
			err = lastBlock.Value(func(val []byte) error {
				block, err = pkgBlock.Deserialize(val)
				return err
			})

			return err
		})

		return err
	})

	return block, err
}

func (pst Repository) GetBlockByKey(key []byte) (*pkgBlock.Block, error) {
	var block *pkgBlock.Block

	err := pst.db.View(func(txn *badger.Txn) error {
		blockSalved, err := txn.Get(key)
		if err != nil {
			return err
		}

		err = blockSalved.Value(func(val []byte) error {
			block, err = pkgBlock.Deserialize(val)
			return err
		})

		return err
	})

	return block, err
}

func (pst Repository) GetOrCreateFirstBlock(firstBlock *pkgBlock.Block) (*pkgBlock.Block, error) {
	var blockToReturn *pkgBlock.Block

	err := pst.db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get(firstBlock.Hash)
		if err != nil && err != badger.ErrKeyNotFound {
			return err
		}

		if err == badger.ErrKeyNotFound {
			err, serialized := firstBlock.Serialize()
			if err != nil {
				return err
			}

			err = txn.Set(firstBlock.Hash, serialized)
			if err != nil {
				return err
			}

			err = txn.Set([]byte("lh"), firstBlock.Hash)
			if err != nil {
				return err
			}

			blockToReturn = firstBlock

			return err
		}

		err = item.Value(func(val []byte) error {
			blockToReturn, err = pkgBlock.Deserialize(val)
			return err
		})

		return err
	})

	return blockToReturn, err
}

func (pst Repository) FindOrCreateBlock(block *pkgBlock.Block) (*pkgBlock.Block, error) {
	var blockToReturn *pkgBlock.Block

	err := pst.db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get(block.Hash)
		if err != nil && err != badger.ErrKeyNotFound {
			return err
		}

		if err == badger.ErrKeyNotFound {
			err, serialized := block.Serialize()
			if err != nil {
				return err
			}

			err = txn.Set(block.Hash, serialized)
			if err != nil {
				return err
			}

			blockToReturn = block

			return err
		}

		err = item.Value(func(val []byte) error {
			blockToReturn, err = pkgBlock.Deserialize(val)
			return err
		})

		return err
	})

	return blockToReturn, err
}

func (pst Repository) UpdateBlock(block *pkgBlock.Block) (*pkgBlock.Block, error) {
	err := pst.db.Update(func(transaction *badger.Txn) error {
		err, serialized := block.Serialize()
		if err != nil {
			return err
		}

		err = transaction.Set(block.Hash, serialized)
		if err != nil {
			return err
		}

		err = transaction.Set([]byte("lh"), block.Hash)
		return err
	})

	return block, err
}

func (pst Repository) Dispose() {
	defer pst.db.Close()
}

func NewRepository(db *badger.DB) *Repository {
	return &Repository{db}
}
