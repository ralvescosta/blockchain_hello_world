package main

type BlockChain struct {
	Blocks []*Block
}

func (pst *BlockChain) Add(data string) {
	index := uint64(len(pst.Blocks) - 1)
	prevBlock := pst.Blocks[index]

	newBlock := NewBlock(data, prevBlock.Hash, index)

	pst.Blocks = append(pst.Blocks, newBlock)
}

func NewBlockchain() BlockChain {
	return BlockChain{[]*Block{NewBlock("Chain Initialized", []byte{}, 0)}}
}
