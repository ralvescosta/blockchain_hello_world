package blockchain

import (
	pkgBlock "blockchain/pkg/block"
	"blockchain/pkg/repositories"
)

type BlockChain struct {
	LastHash []byte
	repo     *repositories.Repository
}

func (pst *BlockChain) Add(data string) error {
	block, err := pst.repo.GetLastBlock()
	if err != nil {
		return err
	}

	newBlock, err := pkgBlock.NewBlock(data, block.Hash, block.NextId())
	if err != nil {
		return err
	}

	_, err = pst.repo.InsertNewBlock(newBlock)
	if err != nil {
		return err
	}

	pst.LastHash = newBlock.Hash
	return err
}

func NewBlockchain(repo *repositories.Repository) (*BlockChain, error) {
	firstBlock, err := pkgBlock.NewBlock("First", []byte{}, 0)
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
	repo        *repositories.Repository
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
