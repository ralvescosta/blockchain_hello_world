package blockchain

import (
	pkgBlock "blockchain/pkg/block"
	"blockchain/pkg/repository"
)

type BlockChain struct {
	LastHash []byte
	repo     *repository.Repository
}

func (pst *BlockChain) Add(data string) error {
	block, err := pst.repo.GetLastBlock()
	if err != nil {
		return err
	}

	err, newBlock := pkgBlock.NewBlock(data, block.Hash, 0)
	if err != nil {
		return err
	}

	_, err = pst.repo.UpdateBlock(newBlock)
	if err != nil {
		return err
	}

	pst.LastHash = newBlock.Hash
	return err
}

func NewBlockchain(repo *repository.Repository) (*BlockChain, error) {
	err, firstBlock := pkgBlock.NewBlock("First", []byte{}, 0)
	if err != nil {
		return nil, err
	}

	block, err := repo.GetOrCreateFirstBlock(firstBlock)
	if err != nil {
		return nil, err
	}

	return &BlockChain{block.Hash, repo}, nil
}

type BlockChainIterator struct {
	CurrentHash []byte
	repo        *repository.Repository
}

func (pst *BlockChain) Iterator() *BlockChainIterator {
	iterator := BlockChainIterator{pst.LastHash, pst.repo}

	return &iterator
}

func (pst *BlockChainIterator) Next() (*pkgBlock.Block, error) {
	block, err := pst.repo.GetBlockByKey(pst.CurrentHash)
	if err != nil {
		return nil, err
	}

	pst.CurrentHash = block.PrevHash

	return block, nil
}
