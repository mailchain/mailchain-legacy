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

package handlers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mailchain/mailchain/crypto/secp256k1/secp256k1test"
	"github.com/mailchain/mailchain/crypto/sr25519/sr25519test"
	"github.com/mailchain/mailchain/encoding/encodingtest"
	"github.com/mailchain/mailchain/internal/mailbox"
	"github.com/mailchain/mailchain/internal/mailbox/mailboxtest"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func Test_parseGetPublicKey(t *testing.T) {
	type args struct {
		queryParams map[string]string
	}
	tests := []struct {
		name    string
		args    args
		wantRes *GetPublicKeyRequest
		wantErr bool
	}{
		{
			"success",
			args{
				map[string]string{
					"address":  "0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761",
					"network":  "mainnet",
					"protocol": "ethereum",
				},
			},
			&GetPublicKeyRequest{
				Address:      "0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761",
				addressBytes: []byte{0x56, 0x2, 0xea, 0x95, 0x54, 0xb, 0xee, 0x46, 0xd0, 0x3b, 0xa3, 0x35, 0xee, 0xd6, 0xf4, 0x9d, 0x11, 0x7e, 0xab, 0x95, 0xc8, 0xab, 0x8b, 0x71, 0xba, 0xe2, 0xcd, 0xd1, 0xe5, 0x64, 0xa7, 0x61},
				Network:      "mainnet",
				Protocol:     "ethereum",
			},
			false,
		},
		{
			"err_empty_address",
			args{
				map[string]string{
					"address":  "",
					"network":  "mainnet",
					"protocol": "ethereum",
				},
			},
			nil,
			true,
		},
		{
			"err_address",
			args{
				map[string]string{
					"address":  "0x560",
					"network":  "mainnet",
					"protocol": "ethereum",
				},
			},
			nil,
			true,
		},
		{
			"err_protocol",
			args{
				map[string]string{
					"address":  "0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761",
					"network":  "mainnet",
					"protocol": "",
				},
			},
			nil,
			true,
		},
		{
			"err_network",
			args{
				map[string]string{
					"address":  "0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761",
					"network":  "",
					"protocol": "ethereum",
				},
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/", nil)
			q := req.URL.Query()
			for k, v := range tt.args.queryParams {
				q.Add(k, v)
			}
			req.URL.RawQuery = q.Encode()
			gotRes, err := parseGetPublicKey(req)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseGetPublicKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.wantRes, gotRes) {
				t.Errorf("parseGetPublicKey() gotRes = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestGetPublicKey(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type args struct {
		finders map[string]mailbox.PubKeyFinder
	}
	tests := []struct {
		name        string
		args        args
		queryParams map[string]string
		wantStatus  int
	}{
		{
			"422-invalid-request",
			args{
				nil,
			},
			map[string]string{},
			http.StatusUnprocessableEntity,
		},
		{
			"422-no-network-finder",
			args{
				func() map[string]mailbox.PubKeyFinder {
					finder := mailboxtest.NewMockPubKeyFinder(mockCtrl)
					return map[string]mailbox.PubKeyFinder{"ethereum.no-network": finder}
				}(),
			},
			map[string]string{
				"address":  "0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761",
				"network":  "mainnet",
				"protocol": "ethereum",
			},
			http.StatusUnprocessableEntity,
		},
		{
			"422-nil-network-finder",
			args{
				map[string]mailbox.PubKeyFinder{"ethereum/mainnet": nil},
			},
			map[string]string{
				"address":  "0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761",
				"network":  "mainnet",
				"protocol": "ethereum",
			},
			http.StatusUnprocessableEntity,
		},
		{
			"406-network-not-supported",
			args{
				func() map[string]mailbox.PubKeyFinder {
					finder := mailboxtest.NewMockPubKeyFinder(mockCtrl)
					finder.EXPECT().PublicKeyFromAddress(gomock.Any(), "ethereum", "mainnet", []byte{0x56, 0x2, 0xea, 0x95, 0x54, 0xb, 0xee, 0x46, 0xd0, 0x3b, 0xa3, 0x35, 0xee, 0xd6, 0xf4, 0x9d, 0x11, 0x7e, 0xab, 0x95, 0xc8, 0xab, 0x8b, 0x71, 0xba, 0xe2, 0xcd, 0xd1, 0xe5, 0x64, 0xa7, 0x61}).Return(nil, errors.New("network not supported")).Times(1)
					return map[string]mailbox.PubKeyFinder{"ethereum/mainnet": finder}
				}(),
			},
			map[string]string{
				"address":  "0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761",
				"network":  "mainnet",
				"protocol": "ethereum",
			},
			http.StatusNotAcceptable,
		},
		{
			"500-public-key-from-address-error",
			args{
				func() map[string]mailbox.PubKeyFinder {
					finder := mailboxtest.NewMockPubKeyFinder(mockCtrl)
					finder.EXPECT().PublicKeyFromAddress(gomock.Any(), "ethereum", "mainnet", []byte{0x56, 0x2, 0xea, 0x95, 0x54, 0xb, 0xee, 0x46, 0xd0, 0x3b, 0xa3, 0x35, 0xee, 0xd6, 0xf4, 0x9d, 0x11, 0x7e, 0xab, 0x95, 0xc8, 0xab, 0x8b, 0x71, 0xba, 0xe2, 0xcd, 0xd1, 0xe5, 0x64, 0xa7, 0x61}).Return(nil, errors.New("error: PublicKeyFromAddress")).Times(1)
					return map[string]mailbox.PubKeyFinder{"ethereum/mainnet": finder}
				}(),
			},
			map[string]string{
				"address":  "0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761",
				"network":  "mainnet",
				"protocol": "ethereum",
			},
			http.StatusInternalServerError,
		},
		{
			"200-sofia-secp256k1",
			args{
				func() map[string]mailbox.PubKeyFinder {
					finder := mailboxtest.NewMockPubKeyFinder(mockCtrl)
					finder.EXPECT().PublicKeyFromAddress(gomock.Any(), "ethereum", "mainnet", encodingtest.MustDecodeHexZeroX("0xD5ab4CE3605Cd590Db609b6b5C8901fdB2ef7FE6")).Return(secp256k1test.SofiaPublicKey, nil).Times(1)
					return map[string]mailbox.PubKeyFinder{"ethereum/mainnet": finder}
				}(),
			},
			map[string]string{
				"address":  "0xD5ab4CE3605Cd590Db609b6b5C8901fdB2ef7FE6",
				"network":  "mainnet",
				"protocol": "ethereum",
			},
			http.StatusOK,
		},
		{
			"200-charlotte-secp256k1",
			args{
				func() map[string]mailbox.PubKeyFinder {
					finder := mailboxtest.NewMockPubKeyFinder(mockCtrl)
					finder.EXPECT().PublicKeyFromAddress(gomock.Any(), "ethereum", "mainnet", encodingtest.MustDecodeHexZeroX("0xD5ab4CE3605Cd590Db609b6b5C8901fdB2ef7FE6")).Return(secp256k1test.CharlottePublicKey, nil).Times(1)
					return map[string]mailbox.PubKeyFinder{"ethereum/mainnet": finder}
				}(),
			},
			map[string]string{
				"address":  "0xD5ab4CE3605Cd590Db609b6b5C8901fdB2ef7FE6",
				"network":  "mainnet",
				"protocol": "ethereum",
			},
			http.StatusOK,
		},
		{
			"200-charlotte-sr25519",
			args{
				func() map[string]mailbox.PubKeyFinder {
					finder := mailboxtest.NewMockPubKeyFinder(mockCtrl)
					finder.EXPECT().PublicKeyFromAddress(gomock.Any(), "substrate", "edgeware-berlin", encodingtest.MustDecodeBase58("5CaLgJUDdDRxw6KQXJY2f5hFkMEEGHvtUPQYDWdSbku42Dv2")).Return(sr25519test.CharlottePublicKey, nil).Times(1)
					return map[string]mailbox.PubKeyFinder{"substrate/edgeware-berlin": finder}
				}(),
			},
			map[string]string{
				"address":  "5CaLgJUDdDRxw6KQXJY2f5hFkMEEGHvtUPQYDWdSbku42Dv2",
				"network":  "edgeware-berlin",
				"protocol": "substrate",
			},
			http.StatusOK,
		},
	}

	for _, tt := range tests {
		testName := t.Name()
		t.Run(tt.name, func(t *testing.T) {
			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			req, _ := http.NewRequest("GET", "/", nil)
			q := req.URL.Query()
			for k, v := range tt.queryParams {
				q.Add(k, v)
			}
			req.URL.RawQuery = q.Encode()

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(GetPublicKey(tt.args.finders))

			// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
			// directly and pass in our Request and ResponseRecorder.
			handler.ServeHTTP(rr, req)

			// Check the status code is what we expect.
			if !assert.Equal(t, tt.wantStatus, rr.Code) {
				t.Errorf("handler returned wrong status code: got %v want %v",
					rr.Code, tt.wantStatus)
			}
			golden, err := ioutil.ReadFile(fmt.Sprintf("./testdata/%s/response-%s.json", testName, tt.name))
			if err != nil {
				assert.FailNow(t, err.Error())
			}
			if !assert.JSONEq(t, string(golden), rr.Body.String()) {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), golden)
			}
		})
	}
}
