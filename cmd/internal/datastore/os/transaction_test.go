package os

import (
	"context"
	"math/big"
	"testing"

	"github.com/mailchain/mailchain/cmd/internal/datastore"
	"github.com/mailchain/mailchain/encoding/encodingtest"
	"github.com/mailchain/mailchain/internal/address/addresstest"
	"github.com/mailchain/mailchain/internal/protocols"
	"github.com/mailchain/mailchain/internal/protocols/ethereum"
	"github.com/spf13/afero"
)

func TestRawTransactionStore_PutRawTransaction(t *testing.T) {

	var txHash = encodingtest.MustDecodeHexZeroX("0x98beb27135aa0a25650557005ad962919d6a278c4b3dde7f4f6a3a1e65aa746c")
	var blockHash = encodingtest.MustDecodeHexZeroX("0x373d339e45a701447367d7b9c7cef84aab79c2b2714271b908cda0ab3ad0849b")
	var txn = &datastore.Transaction{
		From:      addresstest.EthereumCharlotte,
		BlockHash: blockHash,
		Hash:      txHash,
		Data:      encodingtest.MustDecodeHexZeroX("0x6d61696c636861696e383162336636383539326431393338396439656432346664636338316331666630323835383962653535303436303532366631633961613436623864333739346337653032616565363563386631373733376361366637333564393565303965366131396636303838366638313239326535373835373133343562386531653466393238326531306433396637316238636639653731613231656336393939333637346634616261643231623831393531646565346665643565666465663334643131303264346333336538626662613330623461343730646162643434653938653262363439346136653862363963393336353864393631393639356633313561356266356262313865363265336266623237363463363335323631616366363730303862353761316262333838353164396132656635353730323861336166373839646537396234346662346130336137653637393037343030376531623237"),
		To:        addresstest.EthereumSofia,
		Value:     *big.NewInt(int64(0)),
		GasUsed:   *new(big.Int).SetUint64(122207),
		GasPrice:  *new(big.Int).SetUint64(50000000000),
	}

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
				txn.Hash,
				txn,
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
				txn.Hash,
				txn,
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
