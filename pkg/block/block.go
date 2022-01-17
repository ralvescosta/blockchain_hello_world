package block

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"time"

	txn "blockchain/pkg/transaction"
)

type Block struct {
	Id           int64
	Timestamp    int64
	Hash         []byte
	Transactions []*txn.Transaction
	PrevHash     []byte
	Nonce        int
}

func (pst Block) NextId() int64 {
	return pst.Id + 1
}

func (pst Block) ToString() string {
	return fmt.Sprintf("\n[Block] %d\n[Timestamp] %d\n[Transactions] %s\n[Hash] %x\n[PrevHash] %x\n", pst.Id, pst.Timestamp, pst.Transactions, pst.Hash, pst.PrevHash)
}

func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.ID)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]
}

func NewBlock(txs []*txn.Transaction, prevHash []byte, index int64) (*Block, error) {
	block := &Block{
		Id:           index,
		Timestamp:    time.Now().UnixMilli(),
		Transactions: txs,
		Hash:         []byte{},
		PrevHash:     prevHash,
		Nonce:        0,
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
