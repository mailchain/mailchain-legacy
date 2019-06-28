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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/mailchain/mailchain/internal/mailbox"
	"github.com/mailchain/mailchain/internal/mailbox/mailboxtest"
	"github.com/mailchain/mailchain/internal/testutil"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func Test_parseGetPublicKey(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name        string
		args        args
		wantAddress []byte
		wantNetwork string
		wantErr     bool
	}{
		{
			"success",
			args{
				func() *http.Request {
					req := httptest.NewRequest("GET", "/", nil)
					req = mux.SetURLVars(req, map[string]string{
						"address": "5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761",
						"network": "ethereum",
					})
					return req
				}(),
			},
			[]byte{0xee, 0xd6, 0xf4, 0x9d, 0x11, 0x7e, 0xab, 0x95, 0xc8, 0xab, 0x8b, 0x71, 0xba, 0xe2, 0xcd, 0xd1, 0xe5, 0x64, 0xa7, 0x61},
			"ethereum",
			false,
		},
		{
			"err_address",
			args{
				func() *http.Request {
					req := httptest.NewRequest("GET", "/", nil)
					req = mux.SetURLVars(req, map[string]string{
						"network": "ethereum",
					})
					return req
				}(),
			},
			nil,
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAddress, gotNetwork, err := parseGetPublicKey(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseGetPublicKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(tt.wantAddress, gotAddress) {
				t.Errorf("parseGetPublicKey() gotAddress = %v, want %v", gotAddress, tt.wantAddress)
			}
			if gotNetwork != tt.wantNetwork {
				t.Errorf("parseGetPublicKey() gotNetwork = %v, want %v", gotNetwork, tt.wantNetwork)
			}
		})
	}
}

func TestGetPublicKey(t *testing.T) {
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type args struct {
		finders map[string]mailbox.PubKeyFinder
	}
	tests := []struct {
		name       string
		args       args
		req        *http.Request
		wantBody   string
		wantStatus int
	}{
		{
			"err-invalid-request",
			args{
				nil,
			},
			func() *http.Request {
				req := httptest.NewRequest("GET", "/", nil)
				req = mux.SetURLVars(req, map[string]string{
					"address": "",
				})
				return req
			}(),
			"{\"code\":422,\"message\":\"'address' must not be empty\"}\n",
			http.StatusUnprocessableEntity,
		},
		{
			"no-network-finder",
			args{
				func() map[string]mailbox.PubKeyFinder {
					finder := mailboxtest.NewMockPubKeyFinder(mockCtrl)
					return map[string]mailbox.PubKeyFinder{"ethereum.no-network": finder}
				}(),
			},
			func() *http.Request {
				req := httptest.NewRequest("GET", "/", nil)
				req = mux.SetURLVars(req, map[string]string{
					"address": "5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761",
				})
				return req
			}(),
			"{\"code\":422,\"message\":\"no public key finder for chain.network configured\"}\n",
			http.StatusUnprocessableEntity,
		},
		{
			"networkNotSupportedError",
			args{
				func() map[string]mailbox.PubKeyFinder {
					finder := mailboxtest.NewMockPubKeyFinder(mockCtrl)
					finder.EXPECT().PublicKeyFromAddress(gomock.Any(), "mainnet", []byte{0xee, 0xd6, 0xf4, 0x9d, 0x11, 0x7e, 0xab, 0x95, 0xc8, 0xab, 0x8b, 0x71, 0xba, 0xe2, 0xcd, 0xd1, 0xe5, 0x64, 0xa7, 0x61}).Return(nil, errors.New("network not supported")).Times(1)
					return map[string]mailbox.PubKeyFinder{"ethereum/mainnet": finder}
				}(),
			},
			func() *http.Request {
				req := httptest.NewRequest("GET", "/", nil)
				req = mux.SetURLVars(req, map[string]string{
					"address": "5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761",
					"network": "mainnet",
				})
				return req
			}(),
			"{\"code\":406,\"message\":\"network \\\"mainnet\\\" not supported\"}\n",
			http.StatusNotAcceptable,
		},
		{
			"PublicKeyFromAddress_error",
			args{
				func() map[string]mailbox.PubKeyFinder {
					finder := mailboxtest.NewMockPubKeyFinder(mockCtrl)
					finder.EXPECT().PublicKeyFromAddress(gomock.Any(), "mainnet", []byte{0xee, 0xd6, 0xf4, 0x9d, 0x11, 0x7e, 0xab, 0x95, 0xc8, 0xab, 0x8b, 0x71, 0xba, 0xe2, 0xcd, 0xd1, 0xe5, 0x64, 0xa7, 0x61}).Return(nil, errors.New("error")).Times(1)
					return map[string]mailbox.PubKeyFinder{"ethereum/mainnet": finder}
				}(),
			},
			func() *http.Request {
				req := httptest.NewRequest("GET", "/", nil)
				req = mux.SetURLVars(req, map[string]string{
					"address": "5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761",
					"network": "mainnet",
				})
				return req
			}(),
			"{\"code\":500,\"message\":\"error\"}\n",
			http.StatusInternalServerError,
		},
		{
			"success",
			args{
				func() map[string]mailbox.PubKeyFinder {
					finder := mailboxtest.NewMockPubKeyFinder(mockCtrl)
					finder.EXPECT().PublicKeyFromAddress(gomock.Any(), "mainnet", []byte{0xee, 0xd6, 0xf4, 0x9d, 0x11, 0x7e, 0xab, 0x95, 0xc8, 0xab, 0x8b, 0x71, 0xba, 0xe2, 0xcd, 0xd1, 0xe5, 0x64, 0xa7, 0x61}).Return(testutil.MustHexDecodeString("3ada323710def1e02f3586710ae3624ceefba1638e9d9894f724a5401997cd792933ddfd0687874e515a8ab479a38646e6db9f3d8b74d27c4e4eae5a116f9f1400"), nil).Times(1)
					return map[string]mailbox.PubKeyFinder{"ethereum/mainnet": finder}
				}(),
			},
			func() *http.Request {
				req := httptest.NewRequest("GET", "/", nil)
				req = mux.SetURLVars(req, map[string]string{
					"address": "5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761",
					"network": "mainnet",
				})
				return req
			}(),
			"{\"public_key\":\"0x3ada323710def1e02f3586710ae3624ceefba1638e9d9894f724a5401997cd792933ddfd0687874e515a8ab479a38646e6db9f3d8b74d27c4e4eae5a116f9f1400\"}\n",
			http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(GetPublicKey(tt.args.finders))

			// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
			// directly and pass in our Request and ResponseRecorder.
			handler.ServeHTTP(rr, tt.req)

			// Check the status code is what we expect.
			if !assert.Equal(tt.wantStatus, rr.Code) {
				t.Errorf("handler returned wrong status code: got %v want %v",
					rr.Code, tt.wantStatus)
			}
			if !assert.Equal(tt.wantBody, rr.Body.String()) {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), tt.wantBody)
			}
		})
	}
}
