// Copyright 2019 Finobo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ethrpc2

import (
	"context"
	"math/big"
	"testing"

	geth "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/golang/mock/gomock"
	"github.com/mailchain/mailchain/internal/protocols/ethereum"
	"github.com/mailchain/mailchain/internal/mailbox/signer"
	"github.com/mailchain/mailchain/internal/mailbox/signer/signertest"
	"github.com/mailchain/mailchain/internal/testutil"
	"github.com/mailchain/mailchain/sender"
	"github.com/mailchain/mailchain/sender/ethrpc2/ethrpc2test"
	"github.com/pkg/errors"
)

func TestEthRPC2_Send(t *testing.T) {
	value := big.NewInt(0)
	nonce := uint64(12)
	gasPrice := big.NewInt(int64(45))
	gas := uint64(42345)
	networkID := big.NewInt(12)
	to := testutil.MustHexDecodeString("5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761")
	from := testutil.MustHexDecodeString("4cb0a77b76667dac586c40cc9523ace73b5d772bd503c63ed0ca596eae1658b2")
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type signerOpts struct {
		arg1 string
	}
	type fields struct {
		client Client
	}
	type args struct {
		ctx     context.Context
		network string
		to      []byte
		from    []byte
		data    []byte
		signer  signer.Signer
		opts    sender.SendOpts
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
				func() Client {
					m := ethrpc2test.NewMockClient(mockCtrl)
					m.EXPECT().NetworkID(context.Background()).Return(networkID, nil)
					m.EXPECT().SuggestGasPrice(context.Background()).Return(gasPrice, nil)
					to := common.BytesToAddress(to)
					m.EXPECT().EstimateGas(context.Background(), geth.CallMsg{
						To:       &to,
						From:     common.BytesToAddress(from),
						Value:    value,
						Data:     []byte("transactionDataValue"),
						GasPrice: gasPrice,
					}).Return(gas, nil)
					m.EXPECT().NonceAt(context.Background(), common.BytesToAddress(from), nil).Return(uint64(12), nil)
					m.EXPECT().SendTransaction(context.Background(), gomock.AssignableToTypeOf(&types.Transaction{})).Return(nil)
					return m
				}(),
			},
			args{
				context.Background(),
				ethereum.Mainnet,
				to,
				from,
				[]byte("transactionDataValue"),
				func() signer.Signer {
					m := signertest.NewMockSigner(mockCtrl)
					m.EXPECT().Sign(ethereum.SignerOptions{
						Tx:      types.NewTransaction(nonce, common.BytesToAddress(to), value, gas, gasPrice, []byte("transactionDataValue")),
						ChainID: networkID,
					}).Return(&types.Transaction{}, nil)
					return m
				}(),
				signerOpts{"value1"},
			},
			false,
		},
		{
			"err-send-transaction",
			fields{
				func() Client {
					m := ethrpc2test.NewMockClient(mockCtrl)
					m.EXPECT().NetworkID(context.Background()).Return(networkID, nil)
					m.EXPECT().SuggestGasPrice(context.Background()).Return(gasPrice, nil)
					to := common.BytesToAddress(to)
					m.EXPECT().EstimateGas(context.Background(), geth.CallMsg{
						To:       &to,
						From:     common.BytesToAddress(from),
						Value:    value,
						Data:     []byte("transactionDataValue"),
						GasPrice: gasPrice,
					}).Return(gas, nil)
					m.EXPECT().NonceAt(context.Background(), common.BytesToAddress(from), nil).Return(uint64(12), nil)
					m.EXPECT().SendTransaction(context.Background(), gomock.AssignableToTypeOf(&types.Transaction{})).Return(errors.Errorf("failed"))
					return m
				}(),
			},
			args{
				context.Background(),
				ethereum.Mainnet,
				to,
				from,
				[]byte("transactionDataValue"),
				func() signer.Signer {
					m := signertest.NewMockSigner(mockCtrl)
					m.EXPECT().Sign(ethereum.SignerOptions{
						Tx:      types.NewTransaction(nonce, common.BytesToAddress(to), value, gas, gasPrice, []byte("transactionDataValue")),
						ChainID: networkID,
					}).Return(&types.Transaction{}, nil)
					return m
				}(),
				signerOpts{"value1"},
			},
			true,
		},
		{
			"err-invalid-transaction",
			fields{
				func() Client {
					m := ethrpc2test.NewMockClient(mockCtrl)
					m.EXPECT().NetworkID(context.Background()).Return(networkID, nil)
					m.EXPECT().SuggestGasPrice(context.Background()).Return(gasPrice, nil)
					to := common.BytesToAddress(to)
					m.EXPECT().EstimateGas(context.Background(), geth.CallMsg{
						To:       &to,
						From:     common.BytesToAddress(from),
						Value:    value,
						Data:     []byte("transactionDataValue"),
						GasPrice: gasPrice,
					}).Return(gas, nil)
					m.EXPECT().NonceAt(context.Background(), common.BytesToAddress(from), nil).Return(uint64(12), nil)
					return m
				}(),
			},
			args{
				context.Background(),
				ethereum.Mainnet,
				to,
				from,
				[]byte("transactionDataValue"),
				func() signer.Signer {
					m := signertest.NewMockSigner(mockCtrl)
					m.EXPECT().Sign(ethereum.SignerOptions{
						Tx:      types.NewTransaction(nonce, common.BytesToAddress(to), value, gas, gasPrice, []byte("transactionDataValue")),
						ChainID: networkID,
					}).Return(&types.Block{}, nil)
					return m
				}(),
				signerOpts{"value1"},
			},
			true,
		},
		{
			"err-sign",
			fields{
				func() Client {
					m := ethrpc2test.NewMockClient(mockCtrl)
					m.EXPECT().NetworkID(context.Background()).Return(networkID, nil)
					m.EXPECT().SuggestGasPrice(context.Background()).Return(gasPrice, nil)
					to := common.BytesToAddress(to)
					m.EXPECT().EstimateGas(context.Background(), geth.CallMsg{
						To:       &to,
						From:     common.BytesToAddress(from),
						Value:    value,
						Data:     []byte("transactionDataValue"),
						GasPrice: gasPrice,
					}).Return(gas, nil)
					m.EXPECT().NonceAt(context.Background(), common.BytesToAddress(from), nil).Return(uint64(12), nil)
					return m
				}(),
			},
			args{
				context.Background(),
				ethereum.Mainnet,
				to,
				from,
				[]byte("transactionDataValue"),
				func() signer.Signer {
					m := signertest.NewMockSigner(mockCtrl)
					m.EXPECT().Sign(ethereum.SignerOptions{
						Tx:      types.NewTransaction(nonce, common.BytesToAddress(to), value, gas, gasPrice, []byte("transactionDataValue")),
						ChainID: networkID,
					}).Return(nil, errors.Errorf("failed"))
					return m
				}(),
				signerOpts{"value1"},
			},
			true,
		},
		{
			"err-nonce-at",
			fields{
				func() Client {
					m := ethrpc2test.NewMockClient(mockCtrl)
					m.EXPECT().NetworkID(context.Background()).Return(networkID, nil)
					m.EXPECT().SuggestGasPrice(context.Background()).Return(gasPrice, nil)
					to := common.BytesToAddress(to)
					m.EXPECT().EstimateGas(context.Background(), geth.CallMsg{
						To:       &to,
						From:     common.BytesToAddress(from),
						Value:    value,
						Data:     []byte("transactionDataValue"),
						GasPrice: gasPrice,
					}).Return(gas, nil)
					m.EXPECT().NonceAt(context.Background(), common.BytesToAddress(from), nil).Return(uint64(0), errors.New("failed"))
					return m
				}(),
			},
			args{
				context.Background(),
				ethereum.Mainnet,
				to,
				from,
				[]byte("transactionDataValue"),
				func() signer.Signer {
					m := signertest.NewMockSigner(mockCtrl)
					return m
				}(),
				signerOpts{"value1"},
			},
			true,
		},
		{
			"err-estimate-gas",
			fields{
				func() Client {
					m := ethrpc2test.NewMockClient(mockCtrl)
					m.EXPECT().NetworkID(context.Background()).Return(networkID, nil)
					m.EXPECT().SuggestGasPrice(context.Background()).Return(gasPrice, nil)
					to := common.BytesToAddress(to)
					m.EXPECT().EstimateGas(context.Background(), geth.CallMsg{
						To:       &to,
						From:     common.BytesToAddress(from),
						Value:    value,
						Data:     []byte("transactionDataValue"),
						GasPrice: gasPrice,
					}).Return(uint64(0), errors.New("failed"))
					return m
				}(),
			},
			args{
				context.Background(),
				ethereum.Mainnet,
				to,
				from,
				[]byte("transactionDataValue"),
				func() signer.Signer {
					m := signertest.NewMockSigner(mockCtrl)
					return m
				}(),
				signerOpts{"value1"},
			},
			true,
		},
		{
			"err-suggest-gas-price",
			fields{
				func() Client {
					m := ethrpc2test.NewMockClient(mockCtrl)
					m.EXPECT().NetworkID(context.Background()).Return(networkID, nil)
					m.EXPECT().SuggestGasPrice(context.Background()).Return(nil, errors.New("failed"))
					return m
				}(),
			},
			args{
				context.Background(),
				ethereum.Mainnet,
				to,
				from,
				[]byte("transactionDataValue"),
				func() signer.Signer {
					m := signertest.NewMockSigner(mockCtrl)
					return m
				}(),
				signerOpts{"value1"},
			},
			true,
		},
		{
			"err-network-id",
			fields{
				func() Client {
					m := ethrpc2test.NewMockClient(mockCtrl)
					m.EXPECT().NetworkID(context.Background()).Return(nil, errors.New("failed"))
					return m
				}(),
			},
			args{
				context.Background(),
				ethereum.Mainnet,
				to,
				from,
				[]byte("transactionDataValue"),
				func() signer.Signer {
					m := signertest.NewMockSigner(mockCtrl)
					return m
				}(),
				signerOpts{"value1"},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := EthRPC2{
				client: tt.fields.client,
			}
			if err := e.Send(tt.args.ctx, tt.args.network, tt.args.to, tt.args.from, tt.args.data, tt.args.signer, tt.args.opts); (err != nil) != tt.wantErr {
				t.Errorf("EthRPC2.Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
