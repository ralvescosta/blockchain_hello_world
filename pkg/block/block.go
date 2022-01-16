package block

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
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

func (pst Block) NextId() uint64 {
	return pst.Id + 1
}

func (pst Block) ToString() string {
	return fmt.Sprintf("\n[Block] %d\n[Timestamp] %d\n[Data] %x\n[Hash] %x\n[PrevHash] %x\n", pst.Id, pst.Timestamp, pst.Data, pst.Hash, pst.PrevHash)
}

// Converte the entire block to byte, to make it possible to save
func (b *Block) Serialize() ([]byte, error) {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)

	err := encoder.Encode(b)

	if err != nil {
		log.Println("[Err][Block::Serialize]", err)
		return nil, err
	}

	return res.Bytes(), nil
}

// Deserialize the saved block
func Deserialize(data []byte) (*Block, error) {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(&block)

	if err != nil {
		log.Println("[Err][Block::Deserialize]", err)
		return nil, err
	}

	return &block, nil
}

// Create a new Block
func NewBlock(data string, prevHash []byte, index uint64) (*Block, error) {
	block := &Block{
		Id:        index,
		Timestamp: time.Now().UnixMilli(),
		Hash:      []byte{},
		Data:      []byte(data),
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
