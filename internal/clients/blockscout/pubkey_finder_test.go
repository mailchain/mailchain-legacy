// Copyright 2022 Mailchain Ltd.
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

package blockscout

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/secp256k1"
	"github.com/mailchain/mailchain/encoding/encodingtest"
	"github.com/stretchr/testify/assert"
)

func TestGetFromResultHash(t *testing.T) {

	mockTxResult := []txResult{
		{From: "address1", Hash: "aaa111"},
		{From: "address2", Hash: "bbb222"},
		{From: "address3", Hash: "ccc333"},
	}

	type args struct {
		address string
		txList  *txList
	}

	testCases := []struct {
		name    string
		args    args
		want    common.Hash
		wantErr bool
	}{
		{
			"err-empty-transaction-list",
			args{
				"should-not-matter",
				&txList{
					Status:  "",
					Message: "",
					Result:  []txResult{},
				},
			},
			common.Hash{},
			true,
		},
		{
			"match-transaction-1",
			args{
				"address1",
				&txList{
					Status:  "",
					Message: "",
					Result:  mockTxResult,
				},
			},
			common.HexToHash("aaa111"),
			false,
		},
		{
			"match-transaction-3",
			args{
				"address3",
				&txList{
					Status:  "",
					Message: "",
					Result:  mockTxResult,
				},
			},
			common.HexToHash("ccc333"),
			false,
		},
		{
			"err-no-matching-transactions",
			args{
				"address11",
				&txList{
					Status:  "",
					Message: "",
					Result:  mockTxResult,
				},
			},
			common.Hash{},
			true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			hash, err := getFromResultHash(testCase.args.address, testCase.args.txList)

			if (err != nil) != testCase.wantErr {
				t.Errorf("getFromResultHash() error = %v, wantErr %v", err, testCase.wantErr)
				return
			}
			if !assert.Equal(t, testCase.want, hash) {
				t.Errorf("getFromResultHash() = %v, want %v", hash, testCase.want)
			}
		})
	}
}

func TestAPIClient_PublicKeyFromAddress(t *testing.T) {
	type args struct {
		ctx      context.Context
		protocol string
		network  string
		address  []byte
	}
	tests := []struct {
		name    string
		args    args
		want    crypto.PublicKey
		wantErr bool
	}{
		{
			"success",
			args{
				context.Background(),
				"ethereum",
				"mainnet",
				encodingtest.MustDecodeHexZeroX("0x92D8f10248C6a3953CC3692A894655ad05D61Efb"),
			},
			func() crypto.PublicKey {
				k, _ := secp256k1.PublicKeyFromBytes(encodingtest.MustDecodeHexZeroX("0xbdf6fb97c97c126b492186a4d5b28f34f0671a5aacc974da3bde0be93e45a1c50f89ceff72bd04ac9e25a04a1a6cb010aedaf65f91cec8ebe75901c49b63355d"))
				return k
			}(),
			false,
		},
		{
			"err-invalid-public-key",
			args{
				context.Background(),
				"ethereum",
				"mainnet",
				encodingtest.MustDecodeHexZeroX("0x92D8f10248C6a3953CC3692A894655ad05D61Efb"),
			},
			nil,
			true,
		},
		{
			"err-get-transaction-by-hash",
			args{
				context.Background(),
				"ethereum",
				"mainnet",
				encodingtest.MustDecodeHexZeroX("0x92D8f10248C6a3953CC3692A894655ad05D61Efb"),
			},
			nil,
			true,
		},
		{
			"err-get-result-from-hash",
			args{
				context.Background(),
				"ethereum",
				"mainnet",
				encodingtest.MustDecodeHexZeroX("0x92D8f10248C6a3953CC3692A894655ad05D61Efb"),
			},
			nil,
			true,
		},
		{
			"err-get-transactions-by-address",
			args{
				context.Background(),
				"ethereum",
				"mainnet",
				encodingtest.MustDecodeHexZeroX("0x92D8f10248C6a3953CC3692A894655ad05D61Efb"),
			},
			nil,
			true,
		},
		{
			"err-invalid-network",
			args{
				context.Background(),
				"ethereum",
				"invalid",
				encodingtest.MustDecodeHexZeroX("0x92D8f10248C6a3953CC3692A894655ad05D61Efb"),
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		testName := t.Name()
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					golden, err := ioutil.ReadFile(fmt.Sprintf("./testdata/%s/%s-txlist.json", testName, tt.name))
					if err != nil {
						t.Log(r.URL.String())
						assert.FailNow(t, err.Error())
					}
					w.Write([]byte(golden))
				}),
			)
			defer server.Close()
			rpcServer := httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					golden, err := ioutil.ReadFile(fmt.Sprintf("./testdata/%s/%s-rpc.json", testName, tt.name))
					if err != nil {
						t.Log(r.URL.String())
						assert.FailNow(t, err.Error())
					}
					w.Write([]byte(golden))
				}),
			)
			c := APIClient{
				networkConfigs: map[string]networkConfig{
					"mainnet": networkConfig{
						url:    server.URL,
						rpcURL: rpcServer.URL,
					},
				},
			}
			got, err := c.PublicKeyFromAddress(tt.args.ctx, tt.args.protocol, tt.args.network, tt.args.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("APIClient.PublicKeyFromAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("APIClient.PublicKeyFromAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}
