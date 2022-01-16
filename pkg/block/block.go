package block

import (
	"fmt"
	"time"
)

type Block struct {
	Id        int64
	Timestamp int64
	Hash      []byte
	Data      []byte
	PrevHash  []byte
	Nonce     int
}

func (pst Block) NextId() int64 {
	return pst.Id + 1
}

func (pst Block) ToString() string {
	return fmt.Sprintf("\n[Block] %d\n[Timestamp] %d\n[Data] %x\n[Hash] %x\n[PrevHash] %x\n", pst.Id, pst.Timestamp, pst.Data, pst.Hash, pst.PrevHash)
}

func NewBlock(data string, prevHash []byte, index int64) (*Block, error) {
	block := &Block{
		Id:        index,
		Timestamp: time.Now().UnixMilli(),
		Data:      []byte(data),
		Hash:      []byte{},
		PrevHash:  prevHash,
		Nonce:     0,
	}

	pow := NewProofOfWork(block)
	nonce, hash, err := pow.Exec()

	if err != nil {
		return nil, err
	}

	block.Hash = hash
	block.Nonce = nonce

	return block, nil
}
