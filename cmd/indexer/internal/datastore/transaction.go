package datastore

import (
	"context"
	"math/big"
)

type Transaction struct {
	From      []byte
	To        []byte
	Data      []byte
	BlockHash []byte
	Hash      []byte
	Value     big.Int
	GasUsed   big.Int
	GasPrice  big.Int
}

type TransactionStore interface {
	PutTransaction(ctx context.Context, protocol, network string, hash []byte, tx *Transaction) error
	GetTransactionsFrom(ctx context.Context, protocol, network string, address []byte) ([]Transaction, error)
	GetTransactionsTo(ctx context.Context, protocol, network string, address []byte) ([]Transaction, error)
}

type RawTransactionStore interface {
	PutRawTransaction(ctx context.Context, protocol, network string, hash []byte, tx interface{}) error
}
