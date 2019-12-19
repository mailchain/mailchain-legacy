package pq

import (
	"context"
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

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
				"ethereum",
				"mainnet",
				144,
			},
			func() mock {
				db, m, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}
				m.ExpectExec(regexp.QuoteMeta(`UPDATE sync SET block_no = $1, updated_at = $2 WHERE protocol = $3 AND network = $4`)).
					WithArgs(uint64(144), anyTime{}, uint8(1), uint8(1)).
					WillReturnResult(sqlmock.NewResult(1, 1))

				return mock{db, m}
			}(),
			false,
		},
		{
			"err-update-failed",
			args{
				context.Background(),
				"ethereum",
				"mainnet",
				144,
			},
			func() mock {
				db, m, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}
				m.ExpectExec(regexp.QuoteMeta(`UPDATE sync SET block_no = $1, updated_at = $2 WHERE protocol = $3 AND network = $4`)).
					WithArgs(uint64(144), anyTime{}, uint8(1), uint8(1)).
					WillReturnError(sql.ErrNoRows)

				return mock{db, m}
			}(),
			true,
		},

		{
			"err-protocol-network",
			args{
				context.Background(),
				"ethereum",
				"unknown",
				144,
			},
			func() mock {
				db, m, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}

				return mock{db, m}
			}(),
			true,
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := SyncStore{db: sqlx.NewDb(tt.mock.db, "postgres")}
			if err := s.PutBlockNumber(tt.args.ctx, tt.args.protocol, tt.args.network, tt.args.blockNo); (err != nil) != tt.wantErr {
				t.Errorf("SyncStore.PutBlockNumber() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err := tt.mock.sqlmock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
