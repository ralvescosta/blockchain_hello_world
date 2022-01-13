package main

import (
	"github.com/dgraph-io/badger"
)

type BlockChain struct {
	LastHash []byte
	Database *badger.DB
}

const dbPath = "./tmp/blocks"

func (pst *BlockChain) Add(data string) error {
	var lastHash []byte

	err := pst.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		if err != nil {
			return err
		}

		err = item.Value(func(val []byte) error {
			lastHash = val
			return nil
		})
		return err
	})
	if err != nil {
		return err
	}

	err, newBlock := NewBlock(data, lastHash, 0)

	err = pst.Database.Update(func(transaction *badger.Txn) error {
		err, serialized := newBlock.Serialize()
		if err != nil {
			return err
		}

		err = transaction.Set(newBlock.Hash, serialized)
		if err != nil {
			return err
		}

		err = transaction.Set([]byte("lh"), newBlock.Hash)
		pst.LastHash = newBlock.Hash
		return err
	})
	return err
}

func NewBlockchain() (error, *BlockChain) {
	var lastHash []byte
	opts := badger.DefaultOptions(dbPath)
	db, err := badger.Open(opts)

	if err != nil {
		return err, nil
	}

	err = db.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get([]byte("lh")); err == badger.ErrKeyNotFound {
			err, block := NewBlock("Chain Initialized", []byte{}, 0)
			if err != nil {
				return err
			}

			err, serialized := block.Serialize()
			if err != nil {
				return err
			}

			err = txn.Set(block.Hash, serialized)
			if err != nil {
				return err
			}

			err = txn.Set([]byte("lh"), block.Hash)
			if err != nil {
				return err
			}

			lastHash = block.Hash

			return err
		} else {
			item, err := txn.Get([]byte("lh"))
			if err != nil {
				return err
			}
			err = item.Value(func(val []byte) error {
				lastHash = val
				return nil
			})

			return err
		}
	})

	return nil, &BlockChain{lastHash, db}
}

type BlockChainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

func (chain *BlockChain) Iterator() *BlockChainIterator {
	iterator := BlockChainIterator{chain.LastHash, chain.Database}

	return &iterator
}

func (iterator *BlockChainIterator) Next() (error, *Block) {
	var block *Block

	err := iterator.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iterator.CurrentHash)
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
		return err, nil
	}

	iterator.CurrentHash = block.PrevHash

	return nil, block
}
