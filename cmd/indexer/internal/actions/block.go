package actions

import "context"

//go:generate mockgen -source=block.go -package=actionstest -destination=./actionstest/block_mock.go

// Block processes an individual block
type Block interface {
	Run(ctx context.Context, protocol, network string, blk interface{}) error
}
