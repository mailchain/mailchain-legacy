package ethereum_test

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/golang/mock/gomock"
	"github.com/mailchain/mailchain/cmd/indexer/internal/actions"
	"github.com/mailchain/mailchain/cmd/indexer/internal/ethereum"
	"github.com/mailchain/mailchain/cmd/internal/datastore"
	"github.com/mailchain/mailchain/cmd/internal/datastore/datastoretest"
	"github.com/mailchain/mailchain/encoding/encodingtest"
	"github.com/stretchr/testify/assert"
)

func TestTransaction_ToTransaction(t *testing.T) {
	type args struct {
		blk *types.Block
		tx  *types.Transaction
	}
	tests := []struct {
		name    string
		args    args
		want    *datastore.Transaction
		wantErr bool
	}{
		{
			"success",
			args{
				types.NewBlock(&types.Header{}, []*types.Transaction{getTx(t, "0xcd3ccc846c566fbf76f38b1184ba02261a821c6942d18146b89b9d285aa29b9c")}, nil, nil),
				getTx(t, "0xcd3ccc846c566fbf76f38b1184ba02261a821c6942d18146b89b9d285aa29b9c"),
			},
			&datastore.Transaction{
				From:      []uint8{0x11, 0x90, 0x58, 0xdc, 0x2c, 0x57, 0x7e, 0x9c, 0x4b, 0xa6, 0x91, 0x46, 0x78, 0xaa, 0x9d, 0xb5, 0x65, 0x30, 0xf, 0xfe},
				To:        []uint8{0xec, 0xa5, 0x6d, 0x4, 0x54, 0x6a, 0xff, 0xce, 0xc0, 0xb3, 0xce, 0x61, 0x97, 0x11, 0x36, 0xf4, 0x97, 0x86, 0x6a, 0x3b},
				Data:      []uint8{0xb6, 0x1d, 0x27, 0xf6, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x8a, 0x6, 0x6f, 0x5f, 0xde, 0xbe, 0x40, 0x9f, 0xc1, 0x93, 0x9f, 0x51, 0x5e, 0x19, 0xa2, 0x81, 0x63, 0x88, 0x91, 0xfb, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xa, 0x60, 0x4b, 0x9a, 0x42, 0xdf, 0x9c, 0xa0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x60, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
				BlockHash: []uint8{0x52, 0xa7, 0x3a, 0x5c, 0x6d, 0xfc, 0xea, 0x77, 0xa0, 0x32, 0x85, 0x68, 0xad, 0x26, 0x7d, 0xdb, 0xfe, 0xcb, 0x2e, 0xcb, 0x75, 0xa2, 0xfe, 0x5e, 0xed, 0xad, 0xec, 0x58, 0xdb, 0xc, 0xd, 0x93},
				Hash:      []uint8{0xcd, 0x3c, 0xcc, 0x84, 0x6c, 0x56, 0x6f, 0xbf, 0x76, 0xf3, 0x8b, 0x11, 0x84, 0xba, 0x2, 0x26, 0x1a, 0x82, 0x1c, 0x69, 0x42, 0xd1, 0x81, 0x46, 0xb8, 0x9b, 0x9d, 0x28, 0x5a, 0xa2, 0x9b, 0x9c},
				Value:     *big.NewInt(0),
				GasUsed:   *big.NewInt(147670),
				GasPrice:  *big.NewInt(50000000000),
			},
			false,
		},
		{
			"err-to-nil",
			args{
				types.NewBlock(&types.Header{}, []*types.Transaction{getTx(t, "0x157cdd4029ce9bbf319b9d7cc27fd9989e9be63d0427a193f2297e2dee7b4a2e")}, nil, nil),
				getTx(t, "0x157cdd4029ce9bbf319b9d7cc27fd9989e9be63d0427a193f2297e2dee7b4a2e"),
			},
			nil,
			true,
		},
		{
			"err-tx-not-in-block",
			args{
				types.NewBlock(&types.Header{}, []*types.Transaction{}, nil, nil),
				getTx(t, "0xcd3ccc846c566fbf76f38b1184ba02261a821c6942d18146b89b9d285aa29b9c"),
			},
			nil,
			true,
		},
		{
			"err-invalid-from",
			args{
				types.NewBlock(&types.Header{}, []*types.Transaction{getTx(t, "invalid-from")}, nil, nil),
				getTx(t, "invalid-from"),
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transaction := &ethereum.Transaction{}

			got, err := transaction.ToTransaction(tt.args.blk, tt.args.tx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Transaction.toTransaction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("Transaction.toTransaction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransaction_From(t *testing.T) {
	type args struct {
		blockNo *big.Int
		tx      *types.Transaction
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			"success",
			args{
				big.NewInt(500000),
				getTx(t, "0xcd3ccc846c566fbf76f38b1184ba02261a821c6942d18146b89b9d285aa29b9c"),
			},
			[]byte{0x11, 0x90, 0x58, 0xdc, 0x2c, 0x57, 0x7e, 0x9c, 0x4b, 0xa6, 0x91, 0x46, 0x78, 0xaa, 0x9d, 0xb5, 0x65, 0x30, 0xf, 0xfe},
			false,
		},
		{
			"err-invalid",
			args{
				big.NewInt(500000),
				getTx(t, "invalid-from"),
			},
			[]byte(nil),
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transaction := &ethereum.Transaction{}
			got, err := transaction.From(tt.args.blockNo, tt.args.tx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Transaction.From() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("Transaction.From() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransaction_Run(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type fields struct {
		txStore    datastore.TransactionStore
		rawTxStore datastore.RawTransactionStore
		pkStore    datastore.PublicKeyStore
		networkID  *big.Int
	}
	type args struct {
		ctx      context.Context
		protocol string
		network  string
		tx       interface{}
		txOpts   actions.TransactionOptions
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"success",
			fields{
				func() datastore.TransactionStore {
					m := datastoretest.NewMockTransactionStore(mockCtrl)
					return m
				}(),
				func() datastore.RawTransactionStore {
					m := datastoretest.NewMockRawTransactionStore(mockCtrl)
					return m
				}(),
				func() datastore.PublicKeyStore {
					m := datastoretest.NewMockPublicKeyStore(mockCtrl)
					m.EXPECT().PutPublicKey(context.Background(), "ethereum", "mainnet", encodingtest.MustDecodeHexZeroX("0x119058dc2c577e9c4ba6914678aa9db565300ffe"), gomock.Any()).Return(nil)
					return m
				}(),
				big.NewInt(1),
			},
			args{
				context.Background(),
				"ethereum",
				"mainnet",
				getTx(t, "0xcd3ccc846c566fbf76f38b1184ba02261a821c6942d18146b89b9d285aa29b9c"),
				&ethereum.TxOptions{types.NewBlock(&types.Header{}, []*types.Transaction{getTx(t, "0xcd3ccc846c566fbf76f38b1184ba02261a821c6942d18146b89b9d285aa29b9c")}, nil, nil)},
			},
			false,
		},
		{
			"err-to-transation",
			fields{
				func() datastore.TransactionStore {
					m := datastoretest.NewMockTransactionStore(mockCtrl)
					return m
				}(),
				func() datastore.RawTransactionStore {
					m := datastoretest.NewMockRawTransactionStore(mockCtrl)
					return m
				}(),
				func() datastore.PublicKeyStore {
					m := datastoretest.NewMockPublicKeyStore(mockCtrl)
					m.EXPECT().PutPublicKey(context.Background(), "ethereum", "mainnet", encodingtest.MustDecodeHexZeroX("0x119058dc2c577e9c4ba6914678aa9db565300ffe"), gomock.Any()).Return(nil)
					return m
				}(),
				big.NewInt(1),
			},
			args{
				context.Background(),
				"ethereum",
				"mainnet",
				getTx(t, "0x157cdd4029ce9bbf319b9d7cc27fd9989e9be63d0427a193f2297e2dee7b4a2e"),
				&ethereum.TxOptions{types.NewBlock(&types.Header{}, []*types.Transaction{getTx(t, "0x157cdd4029ce9bbf319b9d7cc27fd9989e9be63d0427a193f2297e2dee7b4a2e")}, nil, nil)},
			},
			true,
		},
		{
			"err-put-public-key",
			fields{
				func() datastore.TransactionStore {
					m := datastoretest.NewMockTransactionStore(mockCtrl)
					return m
				}(),
				func() datastore.RawTransactionStore {
					m := datastoretest.NewMockRawTransactionStore(mockCtrl)
					return m
				}(),
				func() datastore.PublicKeyStore {
					m := datastoretest.NewMockPublicKeyStore(mockCtrl)
					m.EXPECT().PutPublicKey(context.Background(), "ethereum", "mainnet", encodingtest.MustDecodeHexZeroX("0x119058dc2c577e9c4ba6914678aa9db565300ffe"), gomock.Any()).Return(errors.New("error"))
					return m
				}(),
				big.NewInt(1),
			},
			args{
				context.Background(),
				"ethereum",
				"mainnet",
				getTx(t, "0xcd3ccc846c566fbf76f38b1184ba02261a821c6942d18146b89b9d285aa29b9c"),
				&ethereum.TxOptions{types.NewBlock(&types.Header{}, []*types.Transaction{getTx(t, "0xcd3ccc846c566fbf76f38b1184ba02261a821c6942d18146b89b9d285aa29b9c")}, nil, nil)},
			},
			true,
		},
		{
			"err-get-public-key",
			fields{
				func() datastore.TransactionStore {
					m := datastoretest.NewMockTransactionStore(mockCtrl)
					return m
				}(),
				func() datastore.RawTransactionStore {
					m := datastoretest.NewMockRawTransactionStore(mockCtrl)
					return m
				}(),
				func() datastore.PublicKeyStore {
					m := datastoretest.NewMockPublicKeyStore(mockCtrl)
					return m
				}(),
				big.NewInt(1),
			},
			args{
				context.Background(),
				"ethereum",
				"mainnet",
				getTx(t, "invalid-from"),
				&ethereum.TxOptions{types.NewBlock(&types.Header{}, []*types.Transaction{getTx(t, "invalid-from")}, nil, nil)},
			},
			true,
		},
		{
			"err-tx-opts",
			fields{
				func() datastore.TransactionStore {
					m := datastoretest.NewMockTransactionStore(mockCtrl)
					return m
				}(),
				func() datastore.RawTransactionStore {
					m := datastoretest.NewMockRawTransactionStore(mockCtrl)
					return m
				}(),
				func() datastore.PublicKeyStore {
					m := datastoretest.NewMockPublicKeyStore(mockCtrl)
					return m
				}(),
				big.NewInt(1),
			},
			args{
				context.Background(),
				"ethereum",
				"mainnet",
				getTx(t, "0xcd3ccc846c566fbf76f38b1184ba02261a821c6942d18146b89b9d285aa29b9c"),
				0,
			},
			true,
		},
		{
			"err-tx",
			fields{
				func() datastore.TransactionStore {
					m := datastoretest.NewMockTransactionStore(mockCtrl)
					return m
				}(),
				func() datastore.RawTransactionStore {
					m := datastoretest.NewMockRawTransactionStore(mockCtrl)
					return m
				}(),
				func() datastore.PublicKeyStore {
					m := datastoretest.NewMockPublicKeyStore(mockCtrl)
					return m
				}(),
				big.NewInt(1),
			},
			args{
				context.Background(),
				"ethereum",
				"mainnet",
				0,
				&ethereum.TxOptions{types.NewBlock(&types.Header{}, []*types.Transaction{getTx(t, "0xcd3ccc846c566fbf76f38b1184ba02261a821c6942d18146b89b9d285aa29b9c")}, nil, nil)},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transaction := ethereum.NewTransactionProcessor(
				tt.fields.txStore,
				tt.fields.rawTxStore,
				tt.fields.pkStore,
				tt.fields.networkID,
			)
			if err := transaction.Run(tt.args.ctx, tt.args.protocol, tt.args.network, tt.args.tx, tt.args.txOpts); (err != nil) != tt.wantErr {
				t.Errorf("Transaction.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
