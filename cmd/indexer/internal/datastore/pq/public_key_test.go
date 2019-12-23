package pq

import (
	"context"
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/secp256k1/secp256k1test"
	"github.com/mailchain/mailchain/internal/protocols"
	"github.com/mailchain/mailchain/internal/protocols/ethereum"
	"github.com/stretchr/testify/assert"
)

var addressBytes = []byte{0xd5, 0xab, 0x4c, 0xe3, 0x60, 0x5c, 0xd5, 0x90, 0xdb, 0x60, 0x9b, 0x6b, 0x5c, 0x89, 0x1, 0xfd, 0xb2, 0xef, 0x7f, 0xe6}

type UnknownPublicKey struct{}

func (pk UnknownPublicKey) Bytes() []byte {
	return []byte("unknown public key")
}

func (pk UnknownPublicKey) Kind() string {
	return "unknown"
}

func (pk UnknownPublicKey) Verify(message, sig []byte) bool {
	return true
}

func TestPublicKeyStore_PutPublicKey(t *testing.T) {
	type args struct {
		ctx       context.Context
		protocol  string
		network   string
		address   []byte
		publicKey crypto.PublicKey
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
				protocols.Ethereum,
				ethereum.Mainnet,
				addressBytes,
				secp256k1test.SofiaPublicKey,
			},
			func() mock {
				db, m, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}
				m.ExpectExec(regexp.QuoteMeta(`UPDATE public_keys SET public_key_type = $1, public_key = $2, updated_at = $3 WHERE protocol = $4 AND network = $5 AND address = $6`)).
					WithArgs(uint64(1), secp256k1test.SofiaPublicKey.Bytes(), anyTime{}, uint8(1), uint8(1), addressBytes).
					WillReturnResult(sqlmock.NewResult(1, 1))

				return mock{db, m}
			}(),
			false,
		},
		{
			"err-update-failed",
			args{
				context.Background(),
				protocols.Ethereum,
				ethereum.Mainnet,
				addressBytes,
				secp256k1test.SofiaPublicKey,
			},
			func() mock {
				db, m, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}
				m.ExpectExec(regexp.QuoteMeta(`UPDATE public_keys SET public_key_type = $1, public_key = $2, updated_at = $3 WHERE protocol = $4 AND network = $5 AND address = $6`)).
					WithArgs(uint64(1), secp256k1test.SofiaPublicKey.Bytes(), anyTime{}, uint8(1), uint8(1), addressBytes).
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
				addressBytes,
				secp256k1test.SofiaPublicKey,
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
		{
			"err-public-key-type",
			args{
				context.Background(),
				protocols.Ethereum,
				ethereum.Mainnet,
				addressBytes,
				func() crypto.PublicKey {
					return &UnknownPublicKey{}
				}(),
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
			s := PublicKeyStore{db: sqlx.NewDb(tt.mock.db, "postgres")}
			if err := s.PutPublicKey(tt.args.ctx, tt.args.protocol, tt.args.network, tt.args.address, tt.args.publicKey); (err != nil) != tt.wantErr {
				t.Errorf("PublicKeyStore.PutPublicKey() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err := tt.mock.sqlmock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestPublicKeyStore_GetPublicKey(t *testing.T) {
	assert := assert.New(t)
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
		pubKey  crypto.PublicKey
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
				addressBytes,
			},
			func() mock {
				db, m, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}
				m.ExpectQuery(regexp.QuoteMeta(`SELECT public_key_type, public_key FROM public_keys WHERE protocol = $1 AND network = $2 AND address = $3`)).
					WithArgs(uint8(1), uint8(1), addressBytes).
					WillReturnRows(
						sqlmock.NewRows([]string{"public_key_type", "public_key"}).
							AddRow(uint8(1), secp256k1test.SofiaPublicKey.Bytes()))

				return mock{db, m}
			}(),
			result{
				secp256k1test.SofiaPublicKey,
				false,
			},
		},
		{
			"err-select-failed",
			args{
				context.Background(),
				protocols.Ethereum,
				ethereum.Mainnet,
				addressBytes,
			},
			func() mock {
				db, m, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}
				m.ExpectQuery(regexp.QuoteMeta(`SELECT public_key_type, public_key FROM public_keys WHERE protocol = $1 AND network = $2 AND address = $3`)).
					WithArgs(uint8(1), uint8(1), addressBytes).
					WillReturnError(sql.ErrNoRows)

				return mock{db, m}
			}(),
			result{
				nil,
				true,
			},
		},
		{
			"err-protocol-network",
			args{
				context.Background(),
				protocols.Ethereum,
				"unknown",
				addressBytes,
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
		{
			"err-public-key-type",
			args{
				context.Background(),
				protocols.Ethereum,
				ethereum.Mainnet,
				addressBytes,
			},
			func() mock {
				db, m, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}
				m.ExpectQuery(regexp.QuoteMeta(`SELECT public_key_type, public_key FROM public_keys WHERE protocol = $1 AND network = $2 AND address = $3`)).
					WithArgs(uint8(1), uint8(1), addressBytes).
					WillReturnRows(
						sqlmock.NewRows([]string{"public_key_type", "public_key"}).
							AddRow(uint8(0), secp256k1test.SofiaPublicKey.Bytes()))

				return mock{db, m}
			}(),
			result{
				nil,
				true,
			},
		},
		{
			"err-public-key",
			args{
				context.Background(),
				protocols.Ethereum,
				ethereum.Mainnet,
				addressBytes,
			},
			func() mock {
				db, m, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}
				m.ExpectQuery(regexp.QuoteMeta(`SELECT public_key_type, public_key FROM public_keys WHERE protocol = $1 AND network = $2 AND address = $3`)).
					WithArgs(uint8(1), uint8(1), addressBytes).
					WillReturnRows(
						sqlmock.NewRows([]string{"public_key_type", "public_key"}).
							AddRow(uint8(1), (&UnknownPublicKey{}).Bytes()))

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
			s := PublicKeyStore{db: sqlx.NewDb(tt.mock.db, "postgres")}
			publicKey, err := s.GetPublicKey(tt.args.ctx, tt.args.protocol, tt.args.network, tt.args.address)

			if (err != nil) != tt.result.wantErr {
				t.Errorf("PublicKeyStore.GetPublicKey() error = %v, wantErr %v", err, tt.result.wantErr)
			}

			if !tt.result.wantErr && !assert.Equal(tt.result.pubKey.Bytes(), publicKey.Bytes()) {
				t.Errorf("PublicKeyStore.GetPublicKey() = %v, want %v", publicKey.Bytes(), tt.result.pubKey.Bytes())
			}

			if err := tt.mock.sqlmock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
