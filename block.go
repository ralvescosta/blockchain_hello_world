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
}

func (pst *Block) DeriveHash() {
	concatenated := bytes.Join([][]byte{pst.Data, pst.PrevHash}, []byte{})

	hash := sha256.Sum256(concatenated)

	pst.Hash = hash[:]
}

func (pst Block) ToString() string {
	return fmt.Sprintf("[Block] %d\n[Time] %d\n[Hash] %s\n[PrevHash] %s", pst.Id, pst.Timestamp, pst.Hash, pst.PrevHash)
}

func NewBlock(data string, prevHash []byte, index uint64) *Block {
	block := Block{
		Id:        index,
		Timestamp: time.Now().Unix(),
		Hash:      []byte{},
		Data:      []byte(data),
		PrevHash:  prevHash,
	}

	block.DeriveHash()

	return &block
}
