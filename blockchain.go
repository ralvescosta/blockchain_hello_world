package main

type BlockChain struct {
	LastHash []byte
	repo     *Repository
}

const dbPath = "./tmp/blocks"

func (pst *BlockChain) Add(data string) error {
	block, err := pst.repo.GetBlockByKey([]byte("lh"))
	if err != nil {
		return err
	}

	err, newBlock := NewBlock(data, block.Hash, 0)

	_, err = pst.repo.UpdateBlock(newBlock)
	if err != nil {
		return err
	}

	pst.LastHash = newBlock.Hash
	return err
}

func NewBlockchain(repo *Repository) (*BlockChain, error) {
	var lastHash []byte
	// opts := badger.DefaultOptions(dbPath)
	// db, err := badger.Open(opts)

	err, firstBlock := NewBlock("Chain Initialized", []byte{}, 0)
	if err != nil {
		return nil, err
	}

	repo.GetOrCreateFirstBlock(firstBlock)

	return &BlockChain{lastHash, repo}, nil
}

type BlockChainIterator struct {
	CurrentHash []byte
	repo        *Repository
}

func (pst *BlockChain) Iterator() *BlockChainIterator {
	iterator := BlockChainIterator{pst.LastHash, pst.repo}

	return &iterator
}

func (pst *BlockChainIterator) Next() (*Block, error) {
	block, err := pst.repo.GetBlockByKey(pst.CurrentHash)
	if err != nil {
		return nil, err
	}

	pst.CurrentHash = block.PrevHash

	return block, nil
}
