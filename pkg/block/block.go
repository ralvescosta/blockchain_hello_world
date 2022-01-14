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

func (pst Block) ToString() string {
	return fmt.Sprintf("\n[Block] %d\n[Timestamp] %d\n[Data] %x\n[Hash] %x\n[PrevHash] %x\n", pst.Id, pst.Timestamp, pst.Data, pst.Hash, pst.PrevHash)
}

// Converte the entire block to byte, to make it possible to save
func (b *Block) Serialize() (error, []byte) {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)

	err := encoder.Encode(b)

	if err != nil {
		log.Println("[Err][Block::Serialize]", err)
		return err, nil
	}

	return nil, res.Bytes()
}

// Deserialize the saved block
func Deserialize(data []byte) (error, *Block) {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(&block)

	if err != nil {
		log.Println("[Err][Block::Deserialize]", err)
		return err, nil
	}

	return nil, &block
}

// Create a new Block
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
