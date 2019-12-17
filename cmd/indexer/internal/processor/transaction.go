package processor

import (
	"bytes"
	"context"

	"github.com/mailchain/mailchain/cmd/indexer/internal/datastore"
	"github.com/mailchain/mailchain/internal/encoding"
)

type Transaction interface {
	Run(ctx context.Context, protocol, network string, tx interface{}, txOpts TransactionOptions) error
}

// TransactionOptions related to different transactions
type TransactionOptions interface{}

func StoreTransaction(ctx context.Context, txStore datastore.TransactionStore, rawTxStore datastore.RawTransactionStore, protocol, network string, tx *datastore.Transaction, rawTx interface{}) error {
	if bytes.HasPrefix(tx.Data, encoding.DataPrefix()) {
		if err := rawTxStore.PutRawTransaction(ctx, protocol, network, tx.Hash, rawTx); err != nil {
			return err
		}
		if err := txStore.PutTransaction(ctx, protocol, network, tx.Hash, tx); err != nil {
			return err
		}
	}
	return nil
}
