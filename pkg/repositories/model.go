package repositories

import (
	"encoding/json"
	"fmt"
	"time"

	pkgBlock "blockchain/pkg/block"
	txn "blockchain/pkg/transaction"
)

type TransactionModel struct{}

func (pst TransactionModel) ToTransactions() []*txn.Transaction {
	return []*txn.Transaction{}
}

func ToTransactinsModel(txs []*txn.Transaction) TransactionModel {
	return TransactionModel{}
}

type BlockModel struct {
	Id           int64            `json:"id"`
	Timestamp    int64            `json:"timestamp"`
	Transactions TransactionModel `json:"transactions"`
	Hash         string           `json:"hash"`
	PrevHash     string           `json:"prev_hash"`
	Nonce        int              `json:"nonce"`
	CreatedAt    string           `json:"created_at"`
}

func (pst BlockModel) ToBlock() *pkgBlock.Block {
	return &pkgBlock.Block{
		Id:           pst.Id,
		Timestamp:    pst.Timestamp,
		Transactions: pst.Transactions.ToTransactions(),
		Hash:         []byte(pst.Hash),
		PrevHash:     []byte(pst.PrevHash),
		Nonce:        pst.Nonce,
	}
}

func (pst BlockModel) MarshalBinary() ([]byte, error) {
	return json.Marshal(pst)
}

func BlockToModel(block *pkgBlock.Block) (BlockModel, error) {
	return BlockModel{
		Id:           int64(block.Id),
		Timestamp:    block.Timestamp,
		Transactions: ToTransactinsModel(block.Transactions),
		Hash:         fmt.Sprintf("%x", block.Hash),
		PrevHash:     string(block.PrevHash),
		Nonce:        block.Nonce,
		CreatedAt:    time.Unix(0, block.Timestamp*1000000).String(),
	}, nil
}
