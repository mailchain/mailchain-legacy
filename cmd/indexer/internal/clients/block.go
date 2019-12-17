package clients

import "context"

type BlockByNumber interface {
	Get(ctx context.Context, blockNo uint64) (blk interface{}, err error)
}
