package blockchain

import (
	pkgBlock "blockchain/pkg/block"
	"blockchain/pkg/repositories"
	txn "blockchain/pkg/transaction"
	"encoding/hex"
)

type Status struct {
	New     bool
	Already bool
}

type BlockChain struct {
	LastHash []byte
	Status   Status
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
				return unspentTxs
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

func blockchainAlreadyExist(address string, repo *repositories.Repository, lastBlock *pkgBlock.Block) (*BlockChain, error) {
	status := Status{Already: true}
	return &BlockChain{lastBlock.Hash, status, repo}, nil
}
func createNewBlockChain(address string, repo *repositories.Repository) (*BlockChain, error) {
	coinbaseTxn := txn.CoinbaseTxn(address, "First Transaction from Genesis")
	coinbaseBlock, err := pkgBlock.NewBlock([]*txn.Transaction{coinbaseTxn}, []byte{}, 0)
	if err != nil {
		return nil, err
	}

	_, err = repo.InsertNewBlock(coinbaseBlock)
	if err != nil {
		return nil, err
	}

	status := Status{New: true}
	return &BlockChain{coinbaseBlock.Hash, status, repo}, nil
}

func NewBlockchain(address string, repo *repositories.Repository) (*BlockChain, error) {
	block, err := repo.GetLastBlock()
	if err != nil {
		return nil, err
	}
	if block != nil {
		return blockchainAlreadyExist(address, repo, block)
	}

	return createNewBlockChain(address, repo)
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
