package main

type BlockChain struct {
	Blocks []*Block
}

func (pst *BlockChain) Add(data string) error {
	index := len(pst.Blocks) - 1
	prevBlock := pst.Blocks[index]

	err, newBlock := NewBlock(data, prevBlock.Hash, prevBlock.Id+1)

	if err != nil {
		return err
	}

	pst.Blocks = append(pst.Blocks, newBlock)
	return nil
}

func NewBlockchain() BlockChain {
	_, block := NewBlock("Chain Initialized", []byte{}, 0)

	return BlockChain{[]*Block{block}}
}
