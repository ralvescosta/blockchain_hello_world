package blockchain

import (
	pkgBlock "blockchain/pkg/block"
	"blockchain/pkg/repositories"
	txn "blockchain/pkg/transaction"
	"encoding/hex"
)

type BlockChain struct {
	LastHash []byte
	repo     *repositories.Repository
}

func (pst *BlockChain) Add(transactions []*txn.Transaction) error {
	block, err := pst.repo.GetLastBlock()
	if err != nil {
		return err
	}

	newBlock, err := pkgBlock.NewBlock(transactions, block.Hash, block.NextId())
	if err != nil {
		return err
	}

	_, err = pst.repo.InsertNewBlock(newBlock)
	if err != nil {
		return err
	}

	pst.LastHash = newBlock.Hash
	return err
}

func (pst *BlockChain) FindUnspentTransactions(address string) []txn.Transaction {
	var unspentTxs []txn.Transaction

	spentTXNs := make(map[string][]int)

	iter := pst.Iterator()

	for {
		block, err := iter.Next()
		if err != nil {
			break
		}

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Outputs {
				if spentTXNs[txID] != nil {
					for _, spentOut := range spentTXNs[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}
				if out.CanBeUnlocked(address) {
					unspentTxs = append(unspentTxs, *tx)
				}
			}
			if !tx.IsCoinbase() {
				for _, in := range tx.Inputs {
					if in.CanUnlock(address) {
						inTxID := hex.EncodeToString(in.ID)
						spentTXNs[inTxID] = append(spentTXNs[inTxID], in.Out)
					}
				}
			}
			if len(block.PrevHash) == 0 {
				break
			}
		}
	}
	return unspentTxs

}

func (pst *BlockChain) FindUnspentTrancationOutputs(address string) []txn.TxOutput {
	var unspentTrancationOutput []txn.TxOutput
	unspentTransactions := pst.FindUnspentTransactions(address)
	for _, tx := range unspentTransactions {
		for _, out := range tx.Outputs {
			if out.CanBeUnlocked(address) {
				unspentTrancationOutput = append(unspentTrancationOutput, out)
			}
		}
	}

	return unspentTrancationOutput
}

func (pst *BlockChain) FindSpendableTransactionOutputs(address string, amount int) (int, map[string][]int) {
	unspentOuts := make(map[string][]int)
	unspentTxs := pst.FindUnspentTransactions(address)
	accumulated := 0

Work:
	for _, tx := range unspentTxs {
		txID := hex.EncodeToString(tx.ID)
		for outIdx, out := range tx.Outputs {
			if out.CanBeUnlocked(address) && accumulated < amount {
				accumulated += out.Value
				unspentOuts[txID] = append(unspentOuts[txID], outIdx)

				if accumulated >= amount {
					break Work
				}
			}
		}
	}
	return accumulated, unspentOuts
}

func NewBlockchain(repo *repositories.Repository) (*BlockChain, error) {
	firstBlock, err := pkgBlock.NewBlock([]*txn.Transaction{}, []byte{}, 0)
	if err != nil {
		return nil, err
	}

	block, err := repo.GetOrCreateFirstBlock(firstBlock)
	if err != nil {
		return nil, err
	}

	return &BlockChain{block.Hash, repo}, nil
}

type BlockChainIterator struct {
	CurrentHash []byte
	repo        *repositories.Repository
}

func (pst *BlockChain) Iterator() *BlockChainIterator {
	iterator := BlockChainIterator{pst.LastHash, pst.repo}

	return &iterator
}

func (pst *BlockChainIterator) Next() (*pkgBlock.Block, error) {
	block, err := pst.repo.GetBlockByKey(pst.CurrentHash)
	if err != nil {
		return nil, err
	}

	pst.CurrentHash = block.PrevHash

	return block, nil
}
