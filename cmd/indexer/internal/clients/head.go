package clients

import "context"

//go:generate mockgen -source=head.go -package=clientstest -destination=./clientstest/head_mock.go

// Head is a client that listens to the new block added the blocks by number and the latest block
type Head interface {
	BlockByNumber(ctx context.Context, blockNo uint64) (blk interface{}, err error)
	LatestBlockNumber(ctx context.Context) (blockNo uint64, err error)
}
