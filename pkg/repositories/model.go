package repositories

import (
	"encoding/json"
	"fmt"
	"time"

	pkgBlock "blockchain/pkg/block"
)

type BlockModel struct {
	Id        int64  `json:"id"`
	Timestamp int64  `json:"timestamp"`
	Data      string `json:"data"`
	Hash      string `json:"hash"`
	PrevHash  string `json:"prev_hash"`
	Nonce     int    `json:"nonce"`
	CreatedAt string `json:"created_at"`
}

func (pst BlockModel) ToBlock() *pkgBlock.Block {
	return &pkgBlock.Block{
		Id:        pst.Id,
		Timestamp: pst.Timestamp,
		Data:      []byte(pst.Data),
		Hash:      []byte(pst.Hash),
		PrevHash:  []byte(pst.PrevHash),
		Nonce:     pst.Nonce,
	}
}

func (pst BlockModel) MarshalBinary() ([]byte, error) {
	return json.Marshal(pst)
}

func BlockToModel(block *pkgBlock.Block) (BlockModel, error) {
	return BlockModel{
		Id:        int64(block.Id),
		Timestamp: block.Timestamp,
		Data:      string(block.Data),
		Hash:      fmt.Sprintf("%x", block.Hash),
		PrevHash:  string(block.PrevHash),
		Nonce:     block.Nonce,
		CreatedAt: time.Unix(0, block.Timestamp*1000000).String(),
	}, nil
}
