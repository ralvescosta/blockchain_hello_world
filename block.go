package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"time"
)

type Block struct {
	Id        uint64
	Timestamp int64
	Hash      []byte
	Data      []byte
	PrevHash  []byte
	Nonce     int
}

func (pst *Block) DeriveHash() {
	concatenated := bytes.Join([][]byte{pst.Data, pst.PrevHash}, []byte{})

	hash := sha256.Sum256(concatenated)

	pst.Hash = hash[:]
}

func (pst Block) ToString() string {
	return fmt.Sprintf("\n[Block] %d\n[Timestamp] %d\n[Data] %x\n[Hash] %x\n[PrevHash] %x\n", pst.Id, pst.Timestamp, pst.Data, pst.Hash, pst.PrevHash)
}

func NewBlock(data string, prevHash []byte, index uint64) (error, *Block) {
	block := &Block{
		Id:        index,
		Timestamp: time.Now().UnixMilli(),
		Hash:      []byte{},
		Data:      []byte(data),
		PrevHash:  prevHash,
		Nonce:     0,
	}

	pow := NewProofOfWork(block)

	err, nonce, hash := pow.Exec()

	if err != nil {
		return err, nil
	}

	block.Hash = hash
	block.Nonce = nonce

	return nil, block
}
