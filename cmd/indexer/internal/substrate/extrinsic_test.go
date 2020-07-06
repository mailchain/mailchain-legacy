package substrate_test

import (
	"context"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/types"
	"github.com/golang/mock/gomock"
	"github.com/mailchain/mailchain/cmd/indexer/internal/actions"
	"github.com/mailchain/mailchain/cmd/indexer/internal/substrate"
	"github.com/mailchain/mailchain/cmd/internal/datastore"
	"github.com/mailchain/mailchain/cmd/internal/datastore/datastoretest"
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
			},
			args{
				context.Background(),
				"subsrate",
				"edgeware-berlin",
				func() *types.Extrinsic {
					b := getBlock(t, "0x7ce3d93396dac53f1ae5fba268afcaa623e224e359507d581439dea791bab971")
					e := b.Extrinsics[2]
					return &e
				}(),
				&substrate.TxOptions{},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			extrinsic := substrate.NewExtrinsicProcessor(
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
