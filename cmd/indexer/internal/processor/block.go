package processor

import "context"

type Block interface {
	Run(ctx context.Context, protocol, network string, blk interface{}) error
}
