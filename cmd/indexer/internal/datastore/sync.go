package datastore

import "context"

//go:generate mockgen -source=sync.go -package=datastoretest -destination=./datastoretest/sync_mock.go

type SyncStore interface {
	GetBlockNumber(ctx context.Context, protocol, network string) (blockNo uint64, err error)
	PutBlockNumber(ctx context.Context, protocol, network string, blockNo uint64) error
}
