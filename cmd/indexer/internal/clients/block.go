package clients

import "context"

// BlockByNumber is a client that gets the blocks by number
type BlockByNumber interface {
	Get(ctx context.Context, blockNo uint64) (blk interface{}, err error)
}
