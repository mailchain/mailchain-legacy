package pq

import (
	"context"
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/jmoiron/sqlx"
	"github.com/mailchain/mailchain/cmd/internal/datastore"
	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/cryptotest"
	"github.com/mailchain/mailchain/crypto/secp256k1/secp256k1test"
	"github.com/mailchain/mailchain/internal/protocols"
	"github.com/mailchain/mailchain/internal/protocols/ethereum"
	"github.com/stretchr/testify/assert"
)

var addressBytes = []byte{0xd5, 0xab, 0x4c, 0xe3, 0x60, 0x5c, 0xd5, 0x90, 0xdb, 0x60, 0x9b, 0x6b, 0x5c, 0x89, 0x1, 0xfd, 0xb2, 0xef, 0x7f, 0xe6}
var txHash = []byte("0x98beb27135aa0a25650557005ad962919d6a278c4b3dde7f4f6a3a1e65aa746c")
var blockHash = []byte("0x373d339e45a701447367d7b9c7cef84aab79c2b2714271b908cda0ab3ad0849b")

func TestPublicKeyStore_PutPublicKey(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	type args struct {
		ctx       context.Context
		protocol  string
		network   string
		address   []byte
		publicKey *datastore.PublicKey
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
				addressBytes,
				&datastore.PublicKey{
					PublicKey: secp256k1test.AlicePublicKey,
					BlockHash: blockHash,
					TxHash:    txHash,
				},
			},
			func() mock {
				db, m, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}
				m.ExpectExec(regexp.QuoteMeta(`INSERT INTO public_keys (protocol,network,address,public_key_type,public_key,created_block_hash,updated_block_hash,created_tx_hash,updated_tx_hash) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) ON CONFLICT (protocol, network, address) DO UPDATE SET public_key_type = $10, public_key = $11, updated_block_hash = $12, updated_tx_hash = $13`)).
					WithArgs(uint8(1), uint8(1), addressBytes, uint64(crypto.IDSECP256K1), secp256k1test.AlicePublicKey.Bytes(), blockHash, blockHash, txHash, txHash, uint64(crypto.IDSECP256K1), secp256k1test.AlicePublicKey.Bytes(), blockHash, txHash).
					WillReturnResult(sqlmock.NewResult(1, 1))

				return mock{db, m}
			}(),
			false,
		},
		{
			"success-upsert",
			args{
				context.Background(),
				protocols.Ethereum,
				ethereum.Mainnet,
				addressBytes,
				&datastore.PublicKey{
					PublicKey: secp256k1test.AlicePublicKey,
					BlockHash: blockHash,
					TxHash:    txHash,
				},
			},
			func() mock {
				db, m, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}
				m.NewRows([]string{"protocol", "network", "address", "public_key_type", "public_key", "created_block_hash", "updated_block_hash", "created_tx_hash", "updated_tx_hash"}).
					AddRow(uint8(1), uint8(1), addressBytes, uint64(crypto.IDSECP256K1), secp256k1test.AlicePublicKey.Bytes(), blockHash, blockHash, txHash, txHash)
				m.ExpectExec(regexp.QuoteMeta(`INSERT INTO public_keys (protocol,network,address,public_key_type,public_key,created_block_hash,updated_block_hash,created_tx_hash,updated_tx_hash) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) ON CONFLICT (protocol, network, address) DO UPDATE SET public_key_type = $10, public_key = $11, updated_block_hash = $12, updated_tx_hash = $13`)).
					WithArgs(uint8(1), uint8(1), addressBytes, uint64(crypto.IDSECP256K1), secp256k1test.AlicePublicKey.Bytes(), blockHash, blockHash, txHash, txHash, uint64(crypto.IDSECP256K1), secp256k1test.AlicePublicKey.Bytes(), blockHash, txHash).
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
				&datastore.PublicKey{
					PublicKey: secp256k1test.AlicePublicKey,
					BlockHash: blockHash,
					TxHash:    txHash,
				},
			},
			func() mock {
				db, m, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}
				m.ExpectExec(regexp.QuoteMeta(`INSERT INTO public_keys (protocol,network,address,public_key_type,public_key,created_block_hash,updated_block_hash,created_tx_hash,updated_tx_hash) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) ON CONFLICT (protocol, network, address) DO UPDATE SET public_key_type = $10, public_key = $11, updated_block_hash = $12, updated_tx_hash = $13`)).
					WithArgs(uint8(1), uint8(1), addressBytes, uint64(crypto.IDSECP256K1), secp256k1test.AlicePublicKey.Bytes(), blockHash, blockHash, txHash, txHash, uint64(crypto.IDSECP256K1), secp256k1test.AlicePublicKey.Bytes(), blockHash, txHash).
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
				&datastore.PublicKey{
					PublicKey: secp256k1test.AlicePublicKey,
					BlockHash: blockHash,
					TxHash:    txHash,
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
		{
			"err-public-key-type",
			args{
				context.Background(),
				protocols.Ethereum,
				ethereum.Mainnet,
				addressBytes,
				func() *datastore.PublicKey {
					mockPublicKey := cryptotest.NewMockPublicKey(mockCtrl)
					mockPublicKey.EXPECT().Kind().Return("unknown").Times(1)
					return &datastore.PublicKey{
						PublicKey: mockPublicKey,
						BlockHash: blockHash,
						TxHash:    txHash,
					}
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
		pubKey  *datastore.PublicKey
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
				m.ExpectQuery(regexp.QuoteMeta(`SELECT public_key_type, public_key, updated_block_hash, updated_tx_hash FROM public_keys WHERE protocol = $1 AND network = $2 AND address = $3`)).
					WithArgs(uint8(1), uint8(1), addressBytes).
					WillReturnRows(
						sqlmock.NewRows([]string{"public_key_type", "public_key", "updated_block_hash", "updated_tx_hash"}).
							AddRow(uint8(crypto.IDSECP256K1), secp256k1test.AlicePublicKey.Bytes(), blockHash, txHash))

				return mock{db, m}
			}(),
			result{
				&datastore.PublicKey{
					PublicKey: secp256k1test.AlicePublicKey,
					BlockHash: blockHash,
					TxHash:    txHash,
				},
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
				m.ExpectQuery(regexp.QuoteMeta(`SELECT public_key_type, public_key, updated_block_hash, updated_tx_hash FROM public_keys WHERE protocol = $1 AND network = $2 AND address = $3`)).
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
				m.ExpectQuery(regexp.QuoteMeta(`SELECT public_key_type, public_key, updated_block_hash, updated_tx_hash FROM public_keys WHERE protocol = $1 AND network = $2 AND address = $3`)).
					WithArgs(uint8(1), uint8(1), addressBytes).
					WillReturnRows(
						sqlmock.NewRows([]string{"public_key_type", "public_key", "updated_block_hash", "updated_tx_hash"}).
							AddRow(uint8(0), secp256k1test.AlicePublicKey.Bytes(), blockHash, txHash))

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
				m.ExpectQuery(regexp.QuoteMeta(`SELECT public_key_type, public_key, updated_block_hash, updated_tx_hash FROM public_keys WHERE protocol = $1 AND network = $2 AND address = $3`)).
					WithArgs(uint8(1), uint8(1), addressBytes).
					WillReturnRows(
						sqlmock.NewRows([]string{"public_key_type", "public_key", "updated_block_hash", "updated_tx_hash"}).
							AddRow(uint8(crypto.IDSECP256K1), []byte("unknown public key"), blockHash, txHash))

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

			if !tt.result.wantErr && !assert.Equal(t, tt.result.pubKey.PublicKey.Bytes(), publicKey.PublicKey.Bytes()) {
				t.Errorf("PublicKeyStore.GetPublicKey() = %v, want %v", publicKey.PublicKey.Bytes(), tt.result.pubKey.PublicKey.Bytes())
			}

			if err := tt.mock.sqlmock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
