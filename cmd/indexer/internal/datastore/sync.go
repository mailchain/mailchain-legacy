package datastore

import "context"

type SyncStore interface {
	GetBlockNumber(ctx context.Context, protocol, network string) (blockNo uint64, err error)
	PutBlockNumber(ctx context.Context, protocol, network string, blockNo uint64) error
}
