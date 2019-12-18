package processor

import "context"

//go:generate mockgen -source=block.go -package=processortest -destination=./processortest/block_mock.go

type Block interface {
	Run(ctx context.Context, protocol, network string, blk interface{}) error
}
