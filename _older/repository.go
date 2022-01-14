package main

import "github.com/dgraph-io/badger"

type Repository struct {
	db *badger.DB
}

func (pst Repository) GetBlockByKey(key []byte) (*Block, error) {
	var block *Block

	err := pst.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}

		err = item.Value(func(val []byte) error {
			err, block = Deserialize(val)
			if err != nil {
				return err
			}

			return nil
		})

		return err
	})
	if err != nil {
		return nil, err
	}

	return block, err
}

func (pst Repository) GetOrCreateFirstBlock(firstBlock *Block) (*Block, error) {
	var blockToReturn *Block

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
			err, blockToReturn = Deserialize(val)
			return err
		})

		return err
	})

	return blockToReturn, err
}

func (pst Repository) FindOrCreateBlock(block *Block) (*Block, error) {
	var blockToReturn *Block

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
			err, blockToReturn = Deserialize(val)
			return err
		})

		return err
	})

	return blockToReturn, err
}

func (pst Repository) UpdateBlock(block *Block) (*Block, error) {
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

func NewRepository(db *badger.DB) Repository {
	return Repository{db}
}
