package interfaces

import (
	pkgBlock "blockchain/pkg/block"
)

type IBlockchainRepository interface {
	GetLastBlock() (*pkgBlock.Block, error)
	GetBlockByKey(key []byte) (*pkgBlock.Block, error)
	GetOrCreateFirstBlock(firstBlock *pkgBlock.Block) (*pkgBlock.Block, error)
	InsertNewBlock(block *pkgBlock.Block) (*pkgBlock.Block, error)
	Dispose()
}
