package noop

import (
	"context"

	"github.com/mailchain/mailchain/cmd/internal/datastore"
)

// NewRawTransactionStore that does nothing.
func NewRawTransactionStore() (datastore.RawTransactionStore, error) {
	return &RawTransactionStore{}, nil
}

// RawTransactionStore object.
type RawTransactionStore struct {
}

// PutRawTransaction writes nothing.
func (s RawTransactionStore) PutRawTransaction(ctx context.Context, protocol, network string, hash []byte, tx interface{}) error {
	return nil
}
