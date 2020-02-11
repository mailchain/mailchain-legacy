package pq

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/mailchain/mailchain/internal/protocols"
	"github.com/mailchain/mailchain/internal/protocols/ethereum"
	"github.com/stretchr/testify/assert"
)

func TestSyncStore_GetblockNumber(t *testing.T) {
	type args struct {
		ctx      context.Context
		protocol string
		network  string
	}
	type mock struct {
		db      *sql.DB
		sqlmock sqlmock.Sqlmock
	}
	type result struct {
		blockNo uint64
		wantErr bool
	}
	tests := []struct {
		name   string
		args   args
		mock   mock
		result result
	}{
		{
			"success",
			args{
				context.Background(),
				protocols.Ethereum,
				ethereum.Mainnet,
			},
			func() mock {
				db, m, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}
				m.ExpectQuery(regexp.QuoteMeta(`SELECT block_no FROM sync WHERE protocol = $1 AND network = $2`)).
					WithArgs(uint8(1), uint8(1)).
					WillReturnRows(sqlmock.NewRows([]string{"block_no"}).AddRow(uint8(144)))

				return mock{db, m}
			}(),
			result{
				uint64(144),
				false,
			},
		},
		{
			"err-select-failed",
			args{
				context.Background(),
				protocols.Ethereum,
				ethereum.Mainnet,
			},
			func() mock {
				db, m, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}
				m.ExpectQuery(regexp.QuoteMeta(`SELECT block_no FROM sync WHERE protocol = $1 AND network = $2`)).
					WithArgs(uint8(1), uint8(1)).
					WillReturnError(sql.ErrNoRows)

				return mock{db, m}
			}(),
			result{
				0,
				true,
			},
		},
		{
			"err-protocol-network",
			args{
				context.Background(),
				protocols.Ethereum,
				"unknown",
			},
			func() mock {
				db, m, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}

				return mock{db, m}
			}(),
			result{
				0,
				true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := SyncStore{db: sqlx.NewDb(tt.mock.db, "postgres")}
			blockNo, err := s.GetBlockNumber(tt.args.ctx, tt.args.protocol, tt.args.network)

			if (err != nil) != tt.result.wantErr {
				t.Errorf("SyncStore.GetBlockNumber() error = %v, wantErr %v", err, tt.result.wantErr)
			} else {
				assert.Equal(t, tt.result.blockNo, blockNo)
			}

			if err := tt.mock.sqlmock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestSyncStore_PutBlockNumber(t *testing.T) {
	type args struct {
		ctx      context.Context
		protocol string
		network  string
		blockNo  uint64
	}
	type mock struct {
		db      *sql.DB
		sqlmock sqlmock.Sqlmock
		time    func() time.Time
	}
	tests := []struct {
		name    string
		args    args
		mock    mock
		wantErr bool
	}{
		{
			"success",
			args{
				context.Background(),
				protocols.Ethereum,
				ethereum.Mainnet,
				144,
			},
			func() mock {
				db, m, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}

				now := time.Now()

				m.ExpectExec(regexp.QuoteMeta(`INSERT INTO sync (protocol,network,block_no,created_at,updated_at) VALUES ($1,$2,$3,$4,$5) ON CONFLICT (protocol, network) DO UPDATE SET block_no = $6, updated_at = $7`)).
					WithArgs(uint8(1), uint8(1), uint64(144), now, now, uint64(144), now).
					WillReturnResult(sqlmock.NewResult(1, 1))

				return mock{db, m, func() time.Time { return now }}
			}(),
			false,
		},
		{
			"err-update-failed",
			args{
				context.Background(),
				protocols.Ethereum,
				ethereum.Mainnet,
				144,
			},
			func() mock {
				db, m, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}

				now := time.Now()

				m.ExpectExec(regexp.QuoteMeta(`INSERT INTO sync (protocol,network,block_no,created_at,updated_at) VALUES ($1,$2,$3,$4,$5) ON CONFLICT (protocol, network) DO UPDATE SET block_no = $6, updated_at = $7`)).
					WithArgs(uint8(1), uint8(1), uint64(144), now, now, uint64(144), now).
					WillReturnError(sql.ErrNoRows)

				return mock{db, m, func() time.Time { return now }}
			}(),
			true,
		},
		{
			"err-protocol-network",
			args{
				context.Background(),
				protocols.Ethereum,
				"unknown",
				144,
			},
			func() mock {
				db, m, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}

				now := time.Now()

				return mock{db, m, func() time.Time { return now }}
			}(),
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := SyncStore{db: sqlx.NewDb(tt.mock.db, "postgres"), now: tt.mock.time}
			if err := s.PutBlockNumber(tt.args.ctx, tt.args.protocol, tt.args.network, tt.args.blockNo); (err != nil) != tt.wantErr {
				t.Errorf("SyncStore.PutBlockNumber() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err := tt.mock.sqlmock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
