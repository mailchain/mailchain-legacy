package processor

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mailchain/mailchain/cmd/indexer/internal/actions"
	"github.com/mailchain/mailchain/cmd/indexer/internal/actions/actionstest"
	"github.com/mailchain/mailchain/cmd/indexer/internal/clients"
	"github.com/mailchain/mailchain/cmd/indexer/internal/clients/clientstest"
	"github.com/mailchain/mailchain/cmd/internal/datastore"
	"github.com/mailchain/mailchain/cmd/internal/datastore/datastoretest"
	"github.com/mailchain/mailchain/internal/protocols"
	"github.com/mailchain/mailchain/internal/protocols/ethereum"
	"github.com/pkg/errors"
)

func TestSequential_NextBlock(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	type args struct {
		ctx      context.Context
		protocol string
		network  string
		store    datastore.SyncStore
		proc     actions.Block
		client   clients.BlockByNumber
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"success",
			args{
				context.Background(),
				protocols.Ethereum,
				ethereum.Mainnet,
				func() datastore.SyncStore {
					store := datastoretest.NewMockSyncStore(mockCtrl)
					store.EXPECT().GetBlockNumber(context.Background(), protocols.Ethereum, ethereum.Mainnet).Return(uint64(1), nil).Times(1)
					store.EXPECT().PutBlockNumber(context.Background(), protocols.Ethereum, ethereum.Mainnet, uint64(2)).Return(nil).Times(1)
					return store
				}(),
				func() actions.Block {
					proc := actionstest.NewMockBlock(mockCtrl)
					proc.EXPECT().Run(context.Background(), protocols.Ethereum, ethereum.Mainnet, "block").Return(nil).Times(1)
					return proc
				}(),
				func() clients.BlockByNumber {
					client := clientstest.NewMockBlockByNumber(mockCtrl)
					client.EXPECT().Get(context.Background(), uint64(2)).Return("block", nil).Times(1)
					return client
				}(),
			},
			false,
		},
		{
			"err-syncStore-GetBlockNumber",
			args{
				context.Background(),
				protocols.Ethereum,
				ethereum.Mainnet,
				func() datastore.SyncStore {
					store := datastoretest.NewMockSyncStore(mockCtrl)
					store.EXPECT().GetBlockNumber(context.Background(), protocols.Ethereum, ethereum.Mainnet).Return(uint64(0), errors.Errorf("error getting block number")).Times(1)
					store.EXPECT().PutBlockNumber(context.Background(), protocols.Ethereum, ethereum.Mainnet, uint64(2)).Return(nil).Times(0)
					return store
				}(),
				func() actions.Block {
					proc := actionstest.NewMockBlock(mockCtrl)
					proc.EXPECT().Run(context.Background(), protocols.Ethereum, ethereum.Mainnet, "block").Return(nil).Times(0)
					return proc
				}(),
				func() clients.BlockByNumber {
					client := clientstest.NewMockBlockByNumber(mockCtrl)
					client.EXPECT().Get(context.Background(), uint64(2)).Return("block", nil).Times(0)
					return client
				}(),
			},
			true,
		},
		{
			"err-blockClient-Get",
			args{
				context.Background(),
				protocols.Ethereum,
				ethereum.Mainnet,
				func() datastore.SyncStore {
					store := datastoretest.NewMockSyncStore(mockCtrl)
					store.EXPECT().GetBlockNumber(context.Background(), protocols.Ethereum, ethereum.Mainnet).Return(uint64(1), nil).Times(1)
					store.EXPECT().PutBlockNumber(context.Background(), protocols.Ethereum, ethereum.Mainnet, uint64(2)).Return(nil).Times(0)
					return store
				}(),
				func() actions.Block {
					proc := actionstest.NewMockBlock(mockCtrl)
					proc.EXPECT().Run(context.Background(), protocols.Ethereum, ethereum.Mainnet, "block").Return(nil).Times(0)
					return proc
				}(),
				func() clients.BlockByNumber {
					client := clientstest.NewMockBlockByNumber(mockCtrl)
					client.EXPECT().Get(context.Background(), uint64(2)).Return(nil, errors.Errorf("error getting block")).Times(1)
					return client
				}(),
			},
			true,
		},
		{
			"err-blockProcessor-Run",
			args{
				context.Background(),
				protocols.Ethereum,
				ethereum.Mainnet,
				func() datastore.SyncStore {
					store := datastoretest.NewMockSyncStore(mockCtrl)
					store.EXPECT().GetBlockNumber(context.Background(), protocols.Ethereum, ethereum.Mainnet).Return(uint64(1), nil).Times(1)
					store.EXPECT().PutBlockNumber(context.Background(), protocols.Ethereum, ethereum.Mainnet, uint64(2)).Return(nil).Times(0)
					return store
				}(),
				func() actions.Block {
					proc := actionstest.NewMockBlock(mockCtrl)
					proc.EXPECT().Run(context.Background(), protocols.Ethereum, ethereum.Mainnet, "block").Return(errors.Errorf("error running processor")).Times(1)
					return proc
				}(),
				func() clients.BlockByNumber {
					client := clientstest.NewMockBlockByNumber(mockCtrl)
					client.EXPECT().Get(context.Background(), uint64(2)).Return("block", nil).Times(1)
					return client
				}(),
			},
			true,
		},
		{
			"err-syncStore-PutBlockNumber",
			args{
				context.Background(),
				protocols.Ethereum,
				ethereum.Mainnet,
				func() datastore.SyncStore {
					store := datastoretest.NewMockSyncStore(mockCtrl)
					store.EXPECT().GetBlockNumber(context.Background(), protocols.Ethereum, ethereum.Mainnet).Return(uint64(1), nil).Times(1)
					store.EXPECT().PutBlockNumber(context.Background(), protocols.Ethereum, ethereum.Mainnet, uint64(2)).Return(errors.Errorf("error writing blck number")).Times(1)
					return store
				}(),
				func() actions.Block {
					proc := actionstest.NewMockBlock(mockCtrl)
					proc.EXPECT().Run(context.Background(), protocols.Ethereum, ethereum.Mainnet, "block").Return(nil).Times(1)
					return proc
				}(),
				func() clients.BlockByNumber {
					client := clientstest.NewMockBlockByNumber(mockCtrl)
					client.EXPECT().Get(context.Background(), uint64(2)).Return("block", nil).Times(1)
					return client
				}(),
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSequential(tt.args.protocol, tt.args.network, tt.args.store, tt.args.proc, tt.args.client)
			err := s.NextBlock(tt.args.ctx)

			if (err != nil) != tt.wantErr {
				t.Errorf("NewSequential.NextBlock() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
