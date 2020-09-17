package substrate_test

import (
	"context"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/types"
	"github.com/golang/mock/gomock"
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
					m.EXPECT().PutTransaction(context.Background(), protocols.Substrate, networks.EdgewareBeresheet, gomock.Any(), gomock.Any())
					return m
				}(),
				func() datastore.RawTransactionStore {
					m := datastoretest.NewMockRawTransactionStore(mockCtrl)
					m.EXPECT().PutRawTransaction(context.Background(), protocols.Substrate, networks.EdgewareBeresheet, gomock.Any(), gomock.Any())
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
				"edgeware-beresheet",
				func() *types.Extrinsic {
					b := getBlock(t, "0x61c83ac28eec9ac530df6f117edc761e3c3a0861f73dc8534f2e0682f1d9ef75")
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
					b := getBlock(t, "0x61c83ac28eec9ac530df6f117edc761e3c3a0861f73dc8534f2e0682f1d9ef75")
					e := b.Extrinsics[0]
					return &e
				}(),
			},
			[]byte{0xc7, 0x6e, 0xcb, 0x62, 0x2e, 0xca, 0x6e, 0x4d, 0xdb, 0xab, 0xd3, 0xfc, 0x50, 0xe, 0x59, 0x9c, 0x7e, 0xce, 0xb, 0x87, 0x16, 0x9d, 0x39, 0x2f, 0xa1, 0xab, 0xc6, 0xce, 0x8f, 0xca, 0x20, 0x32},
			false,
		},
		{
			"success-final-hint",
			args{
				func() *types.Extrinsic {
					b := getBlock(t, "0x61c83ac28eec9ac530df6f117edc761e3c3a0861f73dc8534f2e0682f1d9ef75")
					e := b.Extrinsics[1]
					return &e
				}(),
			},
			[]byte{0x18, 0xd0, 0x2a, 0x17, 0x69, 0x9d, 0x50, 0xba, 0xe, 0x42, 0xa2, 0x7a, 0xa3, 0x52, 0x6d, 0x36, 0x4a, 0x9, 0xbe, 0x30, 0x54, 0xe8, 0x49, 0xa5, 0x74, 0x47, 0x8d, 0xa5, 0x50, 0x5a, 0x4c, 0xc2},
			false,
		},
		{
			"success-contract",
			args{
				func() *types.Extrinsic {
					b := getBlock(t, "0x61c83ac28eec9ac530df6f117edc761e3c3a0861f73dc8534f2e0682f1d9ef75")
					e := b.Extrinsics[2]
					return &e
				}(),
			},
			[]byte{0x3f, 0xcb, 0x9c, 0x48, 0x7d, 0x1a, 0xc, 0x28, 0xfe, 0x27, 0xa7, 0x73, 0xcd, 0x9d, 0x20, 0x12, 0x48, 0x25, 0x33, 0xd9, 0x6, 0x4a, 0xed, 0x91, 0xb2, 0x98, 0x3, 0xaa, 0xcc, 0x3d, 0x4a, 0x9e},
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
