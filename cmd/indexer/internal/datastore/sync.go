package datastore

type SyncStore interface {
	GetBlockNumber(protocol, network string) (blockNo uint64, err error)
	SetBlockNumber(protocol, network string, blockNo uint64) error
}
