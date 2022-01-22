package repositories

import (
	"encoding/json"
	"fmt"
	"time"

	pkgBlock "blockchain/pkg/block"
	txn "blockchain/pkg/transaction"
	"blockchain/pkg/wallet"
)

type BlockModel struct {
	Id           int64              `json:"id"`
	Timestamp    int64              `json:"timestamp"`
	Transactions []TransactionModel `json:"transactions"`
	Hash         string             `json:"hash"`
	PrevHash     string             `json:"prev_hash"`
	Nonce        int                `json:"nonce"`
	CreatedAt    string             `json:"created_at"`
}
type TxOutputModel struct {
	Value  int    `json:"value"`
	PubKey string `json:"pub_key"`
}
type TxInputModel struct {
	ID  []byte `json:"id"`
	Out int    `json:"out"`
	Sig string `json:"sig"`
}
type TransactionModel struct {
	ID      []byte          `json:"id"`
	Inputs  []TxInputModel  `json:"inputs"`
	Outputs []TxOutputModel `json:"outputs"`
}

func (pst BlockModel) ToBlock() *pkgBlock.Block {
	return &pkgBlock.Block{
		Id:           pst.Id,
		Timestamp:    pst.Timestamp,
		Transactions: ToTransactions(pst.Transactions),
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

func ToTransactions(transactionsModel []TransactionModel) []*txn.Transaction {
	var txnTransactions []*txn.Transaction
	var txnInputs []txn.TxInput
	var txnOutputs []txn.TxOutput
	for _, transaction := range transactionsModel {
		for _, input := range transaction.Inputs {
			txnInputs = append(txnInputs, txn.TxInput{
				ID:  input.ID,
				Out: input.Out,
				Sig: input.Sig,
			})
		}

		for _, output := range transaction.Outputs {
			txnOutputs = append(txnOutputs, txn.TxOutput{
				Value:  output.Value,
				PubKey: output.PubKey,
			})
		}

		txnTransactions = append(txnTransactions, &txn.Transaction{
			ID:      transaction.ID,
			Inputs:  txnInputs,
			Outputs: txnOutputs,
		})

		txnInputs = []txn.TxInput{}
		txnOutputs = []txn.TxOutput{}
	}

	return txnTransactions
}

func ToTransactinsModel(txs []*txn.Transaction) []TransactionModel {
	var transactionsModel []TransactionModel
	var inputsModel []TxInputModel
	var outputsModel []TxOutputModel
	for _, transaction := range txs {
		for _, input := range transaction.Inputs {
			inputsModel = append(inputsModel, TxInputModel{
				ID:  input.ID,
				Out: input.Out,
				Sig: input.Sig,
			})
		}

		for _, output := range transaction.Outputs {
			outputsModel = append(outputsModel, TxOutputModel{
				Value:  output.Value,
				PubKey: output.PubKey,
			})
		}

		transactionsModel = append(transactionsModel, TransactionModel{
			ID:      transaction.ID,
			Inputs:  inputsModel,
			Outputs: outputsModel,
		})

		inputsModel = []TxInputModel{}
		outputsModel = []TxOutputModel{}
	}

	return transactionsModel
}

type WalletModel struct{}

func (WalletModel) ToWallet() *wallet.Wallet {
	return &wallet.Wallet{}
}

func ToWalletModel(wlt wallet.Wallet) WalletModel {
	return WalletModel{}
}
