package clients

import "context"

//go:generate mockgen -source=block.go -package=clientstest -destination=./clientstest/block_mock.go

// Block is a client that gets the blocks by number and the latest block
type Block interface {
	BlockByNumber(ctx context.Context, blockNo uint64) (blk interface{}, err error)
	LatestBlockNumber(ctx context.Context) (blockNo uint64, err error)
}
