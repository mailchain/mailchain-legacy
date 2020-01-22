package os

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/mailchain/mailchain/internal/protocols"
	"github.com/mailchain/mailchain/internal/protocols/ethereum"
	"github.com/spf13/afero"
)

var tx = types.NewTransaction(
	uint64(21),
	common.HexToAddress("0xf02c1c8e6114b1dbe8937a39260b5b0a374432bb"),
	big.NewInt(4290000000000000),
	uint64(50000),
	big.NewInt(int64(20000000000)),
	[]byte("hello!"))

func TestRawTransactionStore_PutRawTransaction(t *testing.T) {
	type fields struct {
		fs afero.Fs
	}
	type args struct {
		ctx      context.Context
		protocol string
		network  string
		hash     []byte
		tx       interface{}
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
				afero.NewMemMapFs(),
			},
			args{
				context.Background(),
				protocols.Ethereum,
				ethereum.Mainnet,
				tx.Hash().Bytes(),
				tx,
			},
			false,
		},
		{
			"err-read-only-dir",
			fields{
				afero.NewReadOnlyFs(afero.NewMemMapFs()),
			},
			args{
				context.Background(),
				protocols.Ethereum,
				ethereum.Mainnet,
				tx.Hash().Bytes(),
				tx,
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := RawTransactionStore{fs: tt.fields.fs}
			if err := s.PutRawTransaction(tt.args.ctx, tt.args.protocol, tt.args.network, tt.args.hash, tt.args.tx); (err != nil) != tt.wantErr {
				t.Errorf("RawTransactionStore.PutRawTransaction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
