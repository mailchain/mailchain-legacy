package pq

import (
	"context"
	"database/sql"
	"math/big"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/mailchain/mailchain/cmd/internal/datastore"
	"github.com/mailchain/mailchain/internal/protocols"
	"github.com/mailchain/mailchain/internal/protocols/ethereum"
	"github.com/stretchr/testify/assert"
)

var (
	fromAddress = []byte{0x54, 0xab, 0x4c, 0xe3, 0x60, 0x5c, 0xd5, 0x90, 0xdb, 0x60, 0x9b, 0x6b, 0x5c, 0x89, 0x1, 0xfd, 0xb2, 0xef, 0x7f, 0xe6}
	toAddress   = []byte{0x55, 0xab, 0x4c, 0xe3, 0x60, 0x5c, 0xd5, 0x90, 0xdb, 0x60, 0x9b, 0x6b, 0x5c, 0x89, 0x1, 0xfd, 0xb2, 0xef, 0x7f, 0xe6}
	toAddress2  = []byte{0x57, 0xab, 0x4c, 0xe3, 0x60, 0x5c, 0xd5, 0x90, 0xdb, 0x60, 0x9b, 0x6b, 0x5c, 0x89, 0x1, 0xfd, 0xb2, 0xef, 0x7f, 0xe6}
)

func TestTransactionStore_GetTransactionsFrom(t *testing.T) {
	type args struct {
		ctx      context.Context
		protocol string
		network  string
		address  []byte
	}
	type mock struct {
		db      *sql.DB
		sqlmock sqlmock.Sqlmock
	}
	type result struct {
		txs     []datastore.Transaction
		wantErr bool
	}
	tests := []struct {
		name   string
		args   args
		mock   mock
		result result
	}{
		{
			"err-select-failed",
			args{
				context.Background(),
				protocols.Ethereum,
				ethereum.Mainnet,
				fromAddress,
			},
			func() mock {
				db, m, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}

				m.ExpectQuery(regexp.QuoteMeta(`SELECT tx_from, tx_to, tx_data, tx_block_hash, tx_value, tx_gas_used, tx_gas_price FROM transactions WHERE protocol = $1 AND network = $2 AND tx_from = $3`)).
					WithArgs(uint8(1), uint8(1), fromAddress).
					WillReturnError(sql.ErrNoRows)

				return mock{db, m}
			}(),
			result{
				nil,
				true,
			},
		},
		{
			"success",
			args{
				context.Background(),
				protocols.Ethereum,
				ethereum.Mainnet,
				fromAddress,
			},
			func() mock {
				db, m, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}

				m.ExpectQuery(regexp.QuoteMeta(`SELECT tx_from, tx_to, tx_data, tx_block_hash, tx_value, tx_gas_used, tx_gas_price FROM transactions WHERE protocol = $1 AND network = $2 AND tx_from = $3`)).
					WithArgs(uint8(1), uint8(1), fromAddress).
					WillReturnRows(
						sqlmock.NewRows([]string{"tx_from", "tx_to", "tx_data", "tx_block_hash", "tx_value", "tx_gas_used", "tx_gas_price"}).
							AddRow(fromAddress, toAddress, []byte{0x1}, []byte{0x1}, []byte{0x1}, []byte{0x1}, []byte{0x1}).
							AddRow(fromAddress, toAddress2, []byte{0x2}, []byte{0x2}, []byte{0x2}, []byte{0x2}, []byte{0x2}),
					)

				return mock{db, m}
			}(),
			result{
				[]datastore.Transaction{
					datastore.Transaction{
						From:      fromAddress,
						To:        toAddress,
						Data:      []byte{0x1},
						BlockHash: []byte{0x1},
						Value:     *big.NewInt(1),
						GasUsed:   *big.NewInt(1),
						GasPrice:  *big.NewInt(1),
					},
					datastore.Transaction{
						From:      fromAddress,
						To:        toAddress2,
						Data:      []byte{0x2},
						BlockHash: []byte{0x2},
						Value:     *big.NewInt(2),
						GasUsed:   *big.NewInt(2),
						GasPrice:  *big.NewInt(2),
					},
				},
				false,
			},
		},
		{
			"err-protocol-network",
			args{
				context.Background(),
				protocols.Ethereum,
				"unknown",
				fromAddress,
			},
			func() mock {
				db, m, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}

				return mock{db, m}
			}(),
			result{
				nil,
				true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := TransactionStore{
				db: sqlx.NewDb(tt.mock.db, "postgres"),
			}

			txs, err := s.GetTransactionsFrom(tt.args.ctx, tt.args.protocol, tt.args.network, tt.args.address)
			if (err != nil) != tt.result.wantErr {
				t.Errorf("TransactionStore.GetTransactionsFrom() err = %v, wantErr %v", err, tt.result.wantErr)
			}

			if !tt.result.wantErr && !assert.Equal(t, tt.result.txs, txs) {
				t.Errorf("TransactionStore.GetTransactionFrom() = %v, want %v", txs, tt.result.txs)
			}

			if err := tt.mock.sqlmock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestTransactionStore_GetTransactionsTo(t *testing.T) {
	type args struct {
		ctx      context.Context
		protocol string
		network  string
		address  []byte
	}
	type mock struct {
		db      *sql.DB
		sqlmock sqlmock.Sqlmock
	}
	type result struct {
		txs     []datastore.Transaction
		wantErr bool
	}
	tests := []struct {
		name   string
		args   args
		mock   mock
		result result
	}{
		{
			"err-select-failed",
			args{
				context.Background(),
				protocols.Ethereum,
				ethereum.Mainnet,
				toAddress,
			},
			func() mock {
				db, m, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}

				m.ExpectQuery(regexp.QuoteMeta(`SELECT tx_from, tx_to, tx_data, tx_block_hash, tx_value, tx_gas_used, tx_gas_price FROM transactions WHERE protocol = $1 AND network = $2 AND tx_to = $3`)).
					WithArgs(uint8(1), uint8(1), toAddress).
					WillReturnError(sql.ErrNoRows)

				return mock{db, m}
			}(),
			result{
				nil,
				true,
			},
		},
		{
			"success-toAddress",
			args{
				context.Background(),
				protocols.Ethereum,
				ethereum.Mainnet,
				toAddress,
			},
			func() mock {
				db, m, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}

				m.ExpectQuery(regexp.QuoteMeta(`SELECT tx_from, tx_to, tx_data, tx_block_hash, tx_value, tx_gas_used, tx_gas_price FROM transactions WHERE protocol = $1 AND network = $2 AND tx_to = $3`)).
					WithArgs(uint8(1), uint8(1), toAddress).
					WillReturnRows(
						sqlmock.NewRows([]string{"tx_from", "tx_to", "tx_data", "tx_block_hash", "tx_value", "tx_gas_used", "tx_gas_price"}).
							AddRow(fromAddress, toAddress, []byte{0x1}, []byte{0x1}, []byte{0x1}, []byte{0x1}, []byte{0x1}),
					)

				return mock{db, m}
			}(),
			result{
				[]datastore.Transaction{
					datastore.Transaction{
						From:      fromAddress,
						To:        toAddress,
						Data:      []byte{0x1},
						BlockHash: []byte{0x1},
						Value:     *big.NewInt(1),
						GasUsed:   *big.NewInt(1),
						GasPrice:  *big.NewInt(1),
					},
				},
				false,
			},
		},
		{
			"success-toAddress2",
			args{
				context.Background(),
				protocols.Ethereum,
				ethereum.Mainnet,
				toAddress2,
			},
			func() mock {
				db, m, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}

				m.ExpectQuery(regexp.QuoteMeta(`SELECT tx_from, tx_to, tx_data, tx_block_hash, tx_value, tx_gas_used, tx_gas_price FROM transactions WHERE protocol = $1 AND network = $2 AND tx_to = $3`)).
					WithArgs(uint8(1), uint8(1), toAddress2).
					WillReturnRows(
						sqlmock.NewRows([]string{"tx_from", "tx_to", "tx_data", "tx_block_hash", "tx_value", "tx_gas_used", "tx_gas_price"}).
							AddRow(fromAddress, toAddress2, []byte{0x2}, []byte{0x2}, []byte{0x2}, []byte{0x2}, []byte{0x2}),
					)

				return mock{db, m}
			}(),
			result{
				[]datastore.Transaction{
					datastore.Transaction{
						From:      fromAddress,
						To:        toAddress2,
						Data:      []byte{0x2},
						BlockHash: []byte{0x2},
						Value:     *big.NewInt(2),
						GasUsed:   *big.NewInt(2),
						GasPrice:  *big.NewInt(2),
					},
				},
				false,
			},
		},
		{
			"err-protocol-network",
			args{
				context.Background(),
				protocols.Ethereum,
				"unknown",
				toAddress,
			},
			func() mock {
				db, m, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}

				return mock{db, m}
			}(),
			result{
				nil,
				true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := TransactionStore{
				db: sqlx.NewDb(tt.mock.db, "postgres"),
			}

			txs, err := s.GetTransactionsTo(tt.args.ctx, tt.args.protocol, tt.args.network, tt.args.address)
			if (err != nil) != tt.result.wantErr {
				t.Errorf("TransactionStore.GetTransactionsTo() err = %v, wantErr %v", err, tt.result.wantErr)
			}

			if !tt.result.wantErr && !assert.Equal(t, tt.result.txs, txs) {
				t.Errorf("TransactionStore.GetTransactionsTo() = %v, want %v", txs, tt.result.txs)
			}

			if err := tt.mock.sqlmock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestTransactionStore_PutTransaction(t *testing.T) {
	type args struct {
		ctx      context.Context
		protocol string
		network  string
		hash     []byte
		tx       *datastore.Transaction
	}
	type mock struct {
		db      *sql.DB
		sqlmock sqlmock.Sqlmock
	}
	tests := []struct {
		name    string
		args    args
		mock    mock
		wantErr bool
	}{
		{
			"success-insert",
			args{
				context.Background(),
				protocols.Ethereum,
				ethereum.Mainnet,
				[]byte{0x1},
				&datastore.Transaction{
					From:      fromAddress,
					To:        toAddress,
					Data:      []byte{0x1},
					BlockHash: blockHash,
					Hash:      []byte{0x1},
					Value:     *big.NewInt(1),
					GasUsed:   *big.NewInt(1),
					GasPrice:  *big.NewInt(1),
				},
			},
			func() mock {
				db, m, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}
				m.ExpectExec(regexp.QuoteMeta(`INSERT INTO transactions (protocol,network,hash,tx_from,tx_to,tx_data,tx_block_hash,tx_value,tx_gas_used,tx_gas_price) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) ON CONFLICT DO UPDATE SET tx_from = $, tx_to = $, tx_data = $, tx_block_hash = $, tx_value = $, tx_gas_used = $, tx_gas_price = $`)).
					WithArgs(uint8(1), uint8(1), []byte{0x1}, fromAddress, toAddress, []byte{0x1}, blockHash, []byte{0x1}, []byte{0x1}, []byte{0x1}, fromAddress, toAddress, []byte{0x1}, blockHash, []byte{0x1}, []byte{0x1}, []byte{0x1}).
					WillReturnResult(sqlmock.NewResult(1, 1))

				return mock{db, m}
			}(),
			false,
		},
		{
			"success-update",
			args{
				context.Background(),
				protocols.Ethereum,
				ethereum.Mainnet,
				[]byte{0x1},
				&datastore.Transaction{
					From:      fromAddress,
					To:        toAddress,
					Data:      []byte{0x2},
					BlockHash: blockHash,
					Hash:      []byte{0x1},
					Value:     *big.NewInt(1),
					GasUsed:   *big.NewInt(1),
					GasPrice:  *big.NewInt(1),
				},
			},
			func() mock {
				db, m, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}
				m.NewRows([]string{"protocol", "network", "hash", "tx_from", "tx_to", "tx_data", "tx_block_hash", "tx_value", "tx_gas_used", "tx_gas_price"}).
					AddRow(uint8(1), uint8(1), []byte{0x1}, fromAddress, toAddress, []byte{0x1}, blockHash, []byte{0x1}, []byte{0x1}, []byte{0x1})

				m.ExpectExec(regexp.QuoteMeta(`INSERT INTO transactions (protocol,network,hash,tx_from,tx_to,tx_data,tx_block_hash,tx_value,tx_gas_used,tx_gas_price) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) ON CONFLICT DO UPDATE SET tx_from = $, tx_to = $, tx_data = $, tx_block_hash = $, tx_value = $, tx_gas_used = $, tx_gas_price = $`)).
					WithArgs(uint8(1), uint8(1), []byte{0x1}, fromAddress, toAddress, []byte{0x2}, blockHash, []byte{0x1}, []byte{0x1}, []byte{0x1}, fromAddress, toAddress, []byte{0x2}, blockHash, []byte{0x1}, []byte{0x1}, []byte{0x1}).
					WillReturnResult(sqlmock.NewResult(1, 1))

				return mock{db, m}
			}(),
			false,
		},
		{
			"err-failure",
			args{
				context.Background(),
				protocols.Ethereum,
				ethereum.Mainnet,
				[]byte{0x1},
				&datastore.Transaction{
					From:      fromAddress,
					To:        toAddress,
					Data:      []byte{0x1},
					BlockHash: blockHash,
					Hash:      []byte{0x1},
					Value:     *big.NewInt(1),
					GasUsed:   *big.NewInt(1),
					GasPrice:  *big.NewInt(1),
				},
			},
			func() mock {
				db, m, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}
				m.ExpectExec(regexp.QuoteMeta(`INSERT INTO transactions (protocol,network,hash,tx_from,tx_to,tx_data,tx_block_hash,tx_value,tx_gas_used,tx_gas_price) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) ON CONFLICT DO UPDATE SET tx_from = $, tx_to = $, tx_data = $, tx_block_hash = $, tx_value = $, tx_gas_used = $, tx_gas_price = $`)).
					WithArgs(uint8(1), uint8(1), []byte{0x1}, fromAddress, toAddress, []byte{0x1}, blockHash, []byte{0x1}, []byte{0x1}, []byte{0x1}, fromAddress, toAddress, []byte{0x1}, blockHash, []byte{0x1}, []byte{0x1}, []byte{0x1}).
					WillReturnError(sql.ErrNoRows)

				return mock{db, m}
			}(),
			true,
		},
		{
			"err-protocol-network",
			args{
				context.Background(),
				protocols.Ethereum,
				"unknown",
				[]byte{0x1},
				&datastore.Transaction{
					From:      fromAddress,
					To:        toAddress,
					Data:      []byte{0x1},
					BlockHash: blockHash,
					Hash:      []byte{0x1},
					Value:     *big.NewInt(1),
					GasUsed:   *big.NewInt(1),
					GasPrice:  *big.NewInt(1),
				},
			},
			func() mock {
				db, m, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}

				return mock{db, m}
			}(),
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := TransactionStore{
				db: sqlx.NewDb(tt.mock.db, "postgres"),
			}

			if err := s.PutTransaction(tt.args.ctx, tt.args.protocol, tt.args.network, tt.args.hash, tt.args.tx); (err != nil) != tt.wantErr {
				t.Errorf("TransactionStore.PutTransaction() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err := tt.mock.sqlmock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
