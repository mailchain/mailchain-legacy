package clients

import "context"

//go:generate mockgen -source=block.go -package=clientstest -destination=./clientstest/block_mock.go

// BlockByNumber is a client that gets the blocks by number
type BlockByNumber interface {
	Get(ctx context.Context, blockNo uint64) (blk interface{}, err error)
}
