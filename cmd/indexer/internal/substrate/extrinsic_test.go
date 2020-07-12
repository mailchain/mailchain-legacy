package substrate_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mailchain/go-substrate-rpc-client/types"
	"github.com/mailchain/mailchain/cmd/indexer/internal/actions"
	target "github.com/mailchain/mailchain/cmd/indexer/internal/substrate"
	"github.com/mailchain/mailchain/cmd/internal/datastore"
	"github.com/mailchain/mailchain/cmd/internal/datastore/datastoretest"
	"github.com/mailchain/mailchain/internal/protocols"
	networks "github.com/mailchain/mailchain/internal/protocols/substrate"
	"github.com/stretchr/testify/assert"
)

func TestExtrinsic_Run(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	type fields struct {
		txStore    datastore.TransactionStore
		rawTxStore datastore.RawTransactionStore
		pkStore    datastore.PublicKeyStore
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
					m.EXPECT().PutTransaction(context.Background(), protocols.Substrate, networks.EdgewareBerlin, gomock.Any(), gomock.Any())
					return m
				}(),
				func() datastore.RawTransactionStore {
					m := datastoretest.NewMockRawTransactionStore(mockCtrl)
					m.EXPECT().PutRawTransaction(context.Background(), protocols.Substrate, networks.EdgewareBerlin, gomock.Any(), gomock.Any())
					return m
				}(),
				func() datastore.PublicKeyStore {
					m := datastoretest.NewMockPublicKeyStore(mockCtrl)
					return m
				}(),
			},
			args{
				context.Background(),
				"substrate",
				"edgeware-berlin",
				func() *types.Extrinsic {
					b := getBlock(t, "0x7ce3d93396dac53f1ae5fba268afcaa623e224e359507d581439dea791bab971")
					e := b.Extrinsics[2]
					return &e
				}(),
				&target.TxOptions{},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			extrinsic := target.NewExtrinsicProcessor(
				tt.fields.txStore,
				tt.fields.rawTxStore,
				tt.fields.pkStore,
			)

			if err := extrinsic.Run(tt.args.ctx, tt.args.protocol, tt.args.network, tt.args.tx, tt.args.txOpts); (err != nil) != tt.wantErr {
				t.Errorf("Extrinsic.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_ExtrinsicHash(t *testing.T) {
	type args struct {
		ex *types.Extrinsic
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			"success-timestamp-set",
			args{
				func() *types.Extrinsic {
					b := getBlock(t, "0x7ce3d93396dac53f1ae5fba268afcaa623e224e359507d581439dea791bab971")
					e := b.Extrinsics[0]
					return &e
				}(),
			},
			[]byte{0x94, 0x6c, 0x96, 0xef, 0x50, 0x2b, 0x4f, 0x4, 0x6e, 0x4f, 0xae, 0x52, 0x9, 0x52, 0x3a, 0xe3, 0xaa, 0x59, 0xa5, 0x66, 0x49, 0x2d, 0xf6, 0x4e, 0xdf, 0x38, 0x10, 0xd0, 0x35, 0x7c, 0x35, 0x51},
			false,
		},
		{
			"success-final-hint",
			args{
				func() *types.Extrinsic {
					b := getBlock(t, "0x7ce3d93396dac53f1ae5fba268afcaa623e224e359507d581439dea791bab971")
					e := b.Extrinsics[1]
					return &e
				}(),
			},
			[]byte{0x49, 0xf0, 0xb1, 0x75, 0x20, 0xe9, 0x77, 0x4f, 0x8d, 0x86, 0xf4, 0xba, 0xae, 0x6, 0x15, 0xe9, 0x2c, 0x54, 0xbd, 0x9f, 0x87, 0x35, 0xd6, 0xe8, 0xa7, 0xab, 0xeb, 0x8c, 0x76, 0x2a, 0x67, 0xd4},
			false,
		},
		{
			"success-contract",
			args{
				func() *types.Extrinsic {
					b := getBlock(t, "0x7ce3d93396dac53f1ae5fba268afcaa623e224e359507d581439dea791bab971")
					e := b.Extrinsics[2]
					return &e
				}(),
			},
			[]byte{0xdc, 0xd9, 0x48, 0x65, 0x81, 0x56, 0x13, 0x1f, 0x43, 0xb7, 0x8a, 0x4c, 0x85, 0xe1, 0x63, 0xa8, 0xa1, 0xaa, 0x33, 0x65, 0xa3, 0x92, 0xf4, 0x9e, 0x51, 0x43, 0x65, 0x20, 0x57, 0xfc, 0xf5, 0xe9},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := target.ExtrinsicHash(tt.args.ex)
			if (err != nil) != tt.wantErr {
				t.Errorf("extrinsicHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("extrinsicHash() = %v, want %v", got, tt.want)
			}
		})
	}
}
