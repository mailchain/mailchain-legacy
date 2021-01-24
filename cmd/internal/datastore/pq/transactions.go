package pq

import (
	"context"
	"math/big"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/mailchain/mailchain/cmd/internal/datastore"
	"github.com/pkg/errors"
)

type TransactionStore struct {
	db *sqlx.DB
}

type transaction struct {
	Protocol uint8  `db:"protocol"`
	Network  uint8  `db:"network"`
	Hash     []byte `db:"hash"`

	From        []byte `db:"tx_from"`
	To          []byte `db:"tx_to"`
	Data        []byte `db:"tx_data"`
	BlockHash   []byte `db:"tx_block_hash"`
	BlockNumber int64  `db:"tx_block_no"`
	Value       []byte `db:"tx_value"`
	GasUsed     []byte `db:"tx_gas_used"`
	GasPrice    []byte `db:"tx_gas_price"`
}

func NewTransactionStore(db *sqlx.DB) (datastore.TransactionStore, error) {
	return &TransactionStore{db: db}, nil
}

func (s *TransactionStore) PutTransaction(ctx context.Context, protocol, network string, hash []byte, tx *datastore.Transaction) error {
	p, n, err := getProtocolNetworkUint8(protocol, network)
	if err != nil {
		return errors.WithStack(err)
	}

	sql, args, err := squirrel.Insert("transactions").
		Columns("protocol", "network", "hash", "tx_from", "tx_to", "tx_data", "tx_block_no", "tx_block_hash", "tx_value", "tx_gas_used", "tx_gas_price").
		Values(p, n, hash, tx.From, tx.To, tx.Data, tx.BlockNumber, tx.BlockHash, tx.Value.Bytes(), tx.GasUsed.Bytes(), tx.GasPrice.Bytes()).
		PlaceholderFormat(squirrel.Dollar).
		Suffix("ON CONFLICT (protocol, network, hash) DO UPDATE SET tx_from = $11, tx_to = $12, tx_data = $13, tx_block_no = $14, tx_block_hash = $15, tx_value = $16, tx_gas_used = $17, tx_gas_price = $18",
			tx.From, tx.To, tx.Data, tx.BlockNumber, tx.BlockHash, tx.Value.Bytes(), tx.GasUsed.Bytes(), tx.GasPrice.Bytes()).
		ToSql()

	if err != nil {
		return errors.WithStack(err)
	}

	if _, err := s.db.ExecContext(ctx, sql, args...); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (s *TransactionStore) GetTransactionsFrom(ctx context.Context, protocol, network string, address []byte) ([]datastore.Transaction, error) {
	return s.getTransactions(ctx, protocol, network, squirrel.Eq{"tx_from": address})
}

func (s *TransactionStore) GetTransactionsTo(ctx context.Context, protocol, network string, address []byte) ([]datastore.Transaction, error) {
	return s.getTransactions(ctx, protocol, network, squirrel.Eq{"tx_to": address})
}

func (s *TransactionStore) getTransactions(ctx context.Context, protocol, network string, condition squirrel.Eq) ([]datastore.Transaction, error) {
	p, n, err := getProtocolNetworkUint8(protocol, network)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	sql, args, err := squirrel.Select("hash", "tx_from", "tx_to", "tx_data", "tx_block_no", "tx_block_hash", "tx_value", "tx_gas_used", "tx_gas_price").
		From("transactions").
		PlaceholderFormat(squirrel.Dollar).
		Where(squirrel.Eq{"protocol": p}).
		Where(squirrel.Eq{"network": n}).
		Where(condition).
		ToSql()

	if err != nil {
		return nil, errors.WithStack(err)
	}

	tx := []transaction{}
	if err := s.db.SelectContext(ctx, &tx, sql, args...); err != nil {
		return nil, errors.WithStack(err)
	}

	transactions := make([]datastore.Transaction, 0, len(tx))

	for i := range tx {
		transactions = append(transactions, datastore.Transaction{
			From:        tx[i].From,
			To:          tx[i].To,
			Data:        tx[i].Data,
			BlockHash:   tx[i].BlockHash,
			BlockNumber: tx[i].BlockNumber,
			Hash:        tx[i].Hash,
			Value:       *new(big.Int).SetBytes(tx[i].Value),
			GasUsed:     *new(big.Int).SetBytes(tx[i].GasUsed),
			GasPrice:    *new(big.Int).SetBytes(tx[i].GasPrice),
		})
	}

	return transactions, nil
}
