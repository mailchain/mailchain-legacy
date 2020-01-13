package actions

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/golang/mock/gomock"
	"github.com/mailchain/mailchain/cmd/internal/datastore"
	"github.com/mailchain/mailchain/cmd/internal/datastore/datastoretest"
	"github.com/mailchain/mailchain/internal/address/addresstest"
	"github.com/mailchain/mailchain/internal/protocols"
	"github.com/mailchain/mailchain/internal/protocols/ethereum"
	"github.com/pkg/errors"
)

var txHash = []byte("0x98beb27135aa0a25650557005ad962919d6a278c4b3dde7f4f6a3a1e65aa746c")
var blockHash = []byte("0x373d339e45a701447367d7b9c7cef84aab79c2b2714271b908cda0ab3ad0849b")
var txn = &datastore.Transaction{
	From:      addresstest.EthereumCharlotte,
	BlockHash: blockHash,
	Hash:      txHash,
	Data:      hexutil.MustDecode("0x6d61696c636861696e383162336636383539326431393338396439656432346664636338316331666630323835383962653535303436303532366631633961613436623864333739346337653032616565363563386631373733376361366637333564393565303965366131396636303838366638313239326535373835373133343562386531653466393238326531306433396637316238636639653731613231656336393939333637346634616261643231623831393531646565346665643565666465663334643131303264346333336538626662613330623461343730646162643434653938653262363439346136653862363963393336353864393631393639356633313561356266356262313865363265336266623237363463363335323631616366363730303862353761316262333838353164396132656635353730323861336166373839646537396234346662346130336137653637393037343030376531623237"),
	To:        addresstest.EthereumSofia,
	Value:     *big.NewInt(int64(0)),
	GasUsed:   *new(big.Int).SetUint64(122207),
	GasPrice:  *new(big.Int).SetUint64(50000000000),
}

func TestTransaction_StoreTransaction(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	type args struct {
		ctx        context.Context
		protocol   string
		network    string
		txStore    datastore.TransactionStore
		rawTxStore datastore.RawTransactionStore
		tx         *datastore.Transaction
		rawTx      interface{}
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
				func() datastore.TransactionStore {
					store := datastoretest.NewMockTransactionStore(mockCtrl)
					store.EXPECT().PutTransaction(context.Background(), protocols.Ethereum, ethereum.Mainnet, txHash, txn).Return(nil).Times(1)
					return store
				}(),
				func() datastore.RawTransactionStore {
					store := datastoretest.NewMockRawTransactionStore(mockCtrl)
					store.EXPECT().PutRawTransaction(context.Background(), protocols.Ethereum, ethereum.Mainnet, txHash, "rawTx").Return(nil).Times(1)
					return store
				}(),
				txn,
				"rawTx",
			},
			false,
		},
		{
			"success-invalid-prefix-encoding",
			args{
				context.Background(),
				protocols.Ethereum,
				ethereum.Mainnet,
				func() datastore.TransactionStore {
					store := datastoretest.NewMockTransactionStore(mockCtrl)
					store.EXPECT().PutTransaction(context.Background(), protocols.Ethereum, ethereum.Mainnet, txHash, txn).Return(nil).Times(0)
					return store
				}(),
				func() datastore.RawTransactionStore {
					store := datastoretest.NewMockRawTransactionStore(mockCtrl)
					store.EXPECT().PutRawTransaction(context.Background(), protocols.Ethereum, ethereum.Mainnet, txHash, "rawTx").Return(nil).Times(0)
					return store
				}(),
				&datastore.Transaction{
					From:      addresstest.EthereumCharlotte,
					BlockHash: blockHash,
					Hash:      txHash,
					Data:      hexutil.MustDecode("0xaa61696c636861696e38316233663638"),
					To:        addresstest.EthereumSofia,
					Value:     *big.NewInt(int64(0)),
					GasUsed:   *new(big.Int).SetUint64(122207),
					GasPrice:  *new(big.Int).SetUint64(50000000000),
				},
				"rawTx",
			},
			false,
		},
		{
			"err-PutRawTransaction",
			args{
				context.Background(),
				protocols.Ethereum,
				ethereum.Mainnet,
				func() datastore.TransactionStore {
					store := datastoretest.NewMockTransactionStore(mockCtrl)
					store.EXPECT().PutTransaction(context.Background(), protocols.Ethereum, ethereum.Mainnet, txHash, txn).Return(nil).Times(0)
					return store
				}(),
				func() datastore.RawTransactionStore {
					store := datastoretest.NewMockRawTransactionStore(mockCtrl)
					store.EXPECT().PutRawTransaction(context.Background(), protocols.Ethereum, ethereum.Mainnet, txHash, "rawTx").Return(errors.Errorf("error writing to store")).Times(1)
					return store
				}(),
				txn,
				"rawTx",
			},
			true,
		},
		{
			"err-PutTransaction",
			args{
				context.Background(),
				protocols.Ethereum,
				ethereum.Mainnet,
				func() datastore.TransactionStore {
					store := datastoretest.NewMockTransactionStore(mockCtrl)
					store.EXPECT().PutTransaction(context.Background(), protocols.Ethereum, ethereum.Mainnet, txHash, txn).Return(errors.Errorf("error writing to store")).Times(1)
					return store
				}(),
				func() datastore.RawTransactionStore {
					store := datastoretest.NewMockRawTransactionStore(mockCtrl)
					store.EXPECT().PutRawTransaction(context.Background(), protocols.Ethereum, ethereum.Mainnet, txHash, "rawTx").Return(nil).Times(1)
					return store
				}(),
				txn,
				"rawTx",
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := StoreTransaction(tt.args.ctx, tt.args.txStore, tt.args.rawTxStore, tt.args.protocol, tt.args.network, tt.args.tx, tt.args.rawTx)

			if (err != nil) != tt.wantErr {
				t.Errorf("StoreTransaction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
