package transaction

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
)

const reward = 100

type TxOutput struct {
	Value  int
	PubKey string
}

type TxInput struct {
	ID  []byte
	Out int
	Sig string
}

type Transaction struct {
	ID      []byte
	Inputs  []TxInput  // receiver
	Outputs []TxOutput // crypto owner
}

func (tx *Transaction) SetID() error {
	var encoded bytes.Buffer
	var hash [32]byte

	encoder := gob.NewEncoder(&encoded)
	err := encoder.Encode(tx)
	if err != nil {
		return err
	}

	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]

	return nil
}

func (in *TxInput) CanUnlock(pubKey string) bool {
	return in.Sig == pubKey
}
func (out *TxOutput) CanBeUnlocked(sig string) bool {
	return out.PubKey == sig
}

func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].ID) == 0 && tx.Inputs[0].Out == -1
}

// First Transaction. The first transaction into a blockchain is called "coinbase"
func CoinbaseTxn(toAddress, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Coins to %s", toAddress)
	}

	txIn := TxInput{[]byte{}, -1, data}

	txOut := TxOutput{reward, toAddress}

	tx := Transaction{nil, []TxInput{txIn}, []TxOutput{txOut}}

	return &tx
}
