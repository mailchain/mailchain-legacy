package substrate_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mailchain/go-substrate-rpc-client/types"
	"github.com/mailchain/mailchain/cmd/indexer/internal/actions"
	"github.com/mailchain/mailchain/cmd/indexer/internal/actions/actionstest"
	"github.com/mailchain/mailchain/cmd/indexer/internal/substrate"
	"github.com/mailchain/mailchain/internal/protocols"
	networks "github.com/mailchain/mailchain/internal/protocols/substrate"
)

func TestBlock_Run(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type fields struct {
		txProcessor actions.Transaction
	}
	type args struct {
		ctx      context.Context
		protocol string
		network  string
		blk      interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"err-arg",
			fields{
				func() actions.Transaction {
					m := actionstest.NewMockTransaction(mockCtrl)
					return m
				}(),
			},
			args{
				context.Background(),
				protocols.Substrate,
				networks.EdgewareBeresheet,
				&types.Extrinsic{},
			},
			true,
		},
		{
			"err-run",
			fields{
				func() actions.Transaction {
					m := actionstest.NewMockTransaction(mockCtrl)
					m.EXPECT().Run(context.Background(), protocols.Substrate, networks.EdgewareBeresheet, gomock.Any(), gomock.Any()).Return(errors.New("error"))
					return m
				}(),
			},
			args{
				context.Background(),
				protocols.Substrate,
				networks.EdgewareBeresheet,
				&types.Block{
					Header: types.Header{},
					Extrinsics: []types.Extrinsic{
						types.Extrinsic{},
						types.Extrinsic{
							Signature: types.ExtrinsicSignatureV4{
								Signer: types.Address{
									IsAccountID: true,
								},
							},
						},
					},
				},
			},
			true,
		},
		{
			"success",
			fields{
				func() actions.Transaction {
					m := actionstest.NewMockTransaction(mockCtrl)
					m.EXPECT().Run(context.Background(), protocols.Substrate, networks.EdgewareBeresheet, gomock.Any(), gomock.Any()).Return(nil)
					return m
				}(),
			},
			args{
				context.Background(),
				protocols.Substrate,
				networks.EdgewareBeresheet,
				&types.Block{
					Header: types.Header{},
					Extrinsics: []types.Extrinsic{
						types.Extrinsic{},
						types.Extrinsic{
							Signature: types.ExtrinsicSignatureV4{
								Signer: types.Address{
									IsAccountID: true,
								},
							},
						},
					},
				},
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := substrate.NewBlockProcessor(tt.fields.txProcessor)
			if err := b.Run(tt.args.ctx, tt.args.protocol, tt.args.network, tt.args.blk); (err != nil) != tt.wantErr {
				t.Errorf("Block.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
