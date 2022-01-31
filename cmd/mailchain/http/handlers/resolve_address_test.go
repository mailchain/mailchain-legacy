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

package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/mailchain/mailchain/encoding/encodingtest"
	"github.com/mailchain/mailchain/nameservice"
	"github.com/mailchain/mailchain/nameservice/nameservicetest"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func Test_parseGetResolveAddressRequest(t *testing.T) {
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name         string
		args         args
		wantProtocol string
		wantNetwork  string
		wantAddress  []byte
		wantErr      bool
	}{
		{
			"success-ethereum",
			args{
				func() *http.Request {
					req := httptest.NewRequest("GET", "/?network=mainnet&protocol=ethereum", nil)
					req = mux.SetURLVars(req, map[string]string{
						"address": "0x4ad2b251246aafc2f3bdf3b690de3bf906622c51",
					})
					return req
				}(),
			},
			"ethereum",
			"mainnet",
			encodingtest.MustDecodeHex("4ad2b251246aafc2f3bdf3b690de3bf906622c51"),
			false,
		},
		{
			"success-substrate",
			args{
				func() *http.Request {
					req := httptest.NewRequest("GET", "/?network=beresheet&protocol=substrate", nil)
					req = mux.SetURLVars(req, map[string]string{
						"address": "5CaLgJUDdDRxw6KQXJY2f5hFkMEEGHvtUPQYDWdSbku42Dv2",
					})
					return req
				}(),
			},
			"substrate",
			"beresheet",
			[]byte{0x2a, 0x16, 0x9a, 0x11, 0x72, 0x18, 0x51, 0xf5, 0xdf, 0xf3, 0x54, 0x1d, 0xd5, 0xc4, 0xb0, 0xb4, 0x78, 0xac, 0x1c, 0xd0, 0x92, 0xc9, 0xd5, 0x97, 0x6e, 0x83, 0xda, 0xa0, 0xd0, 0x3f, 0x26, 0x62, 0xc, 0x46, 0x4b},
			false,
		},
		{
			"err-invalid-address",
			args{
				func() *http.Request {
					req := httptest.NewRequest("GET", "/?network=mainnet&protocol=ethereum", nil)
					req = mux.SetURLVars(req, map[string]string{
						"address": "0x4",
					})
					return req
				}(),
			},
			"",
			"",
			nil,
			true,
		},
		{
			"err-protocol",
			args{
				func() *http.Request {
					req := httptest.NewRequest("GET", "/?network=mainnet", nil)
					req = mux.SetURLVars(req, map[string]string{
						"domain-name": "address.ens",
					})
					return req
				}(),
			},
			"",
			"",
			nil,
			true,
		},
		{
			"err-network",
			args{
				func() *http.Request {
					req := httptest.NewRequest("GET", "/?protocol=ethereum", nil)
					req = mux.SetURLVars(req, map[string]string{
						"domain-name": "address.ens",
					})
					return req
				}(),
			},
			"",
			"",
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotProtocol, gotNetwork, gotAddress, err := parseGetResolveAddressRequest(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseGetResolveAddressRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotProtocol != tt.wantProtocol {
				t.Errorf("parseGetResolveAddressRequest() gotProtocol = %v, want %v", gotProtocol, tt.wantProtocol)
			}
			if gotNetwork != tt.wantNetwork {
				t.Errorf("parseGetResolveAddressRequest() gotNetwork = %v, want %v", gotNetwork, tt.wantNetwork)
			}
			if !assert.Equal(t, tt.wantAddress, gotAddress) {
				t.Errorf("parseGetResolveAddressRequest() gotAddress = %v, want %v", gotAddress, tt.wantAddress)
			}
		})
	}
}

func TestGetResolveAddress(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type args struct {
		resolvers map[string]nameservice.ReverseLookup
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
				req = mux.SetURLVars(req, map[string]string{})
				return req
			}(),
			"{\"code\":422,\"message\":\"'protocol' must be specified exactly once\"}\n",
			http.StatusUnprocessableEntity,
		},
		{
			"no-network-finder",
			args{
				func() map[string]nameservice.ReverseLookup {
					m := nameservicetest.NewMockReverseLookup(mockCtrl)
					return map[string]nameservice.ReverseLookup{"ethereum.no-network": m}
				}(),
			},
			func() *http.Request {
				req := httptest.NewRequest("GET", "/?network=mainnet&protocol=ethereum", nil)
				req = mux.SetURLVars(req, map[string]string{
					"address": "0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761",
				})
				return req
			}(),
			"{\"code\":422,\"message\":\"nameserver not supported on \\\"ethereum/mainnet\\\"\"}\n",
			http.StatusUnprocessableEntity,
		},
		{
			"nil-network-finder",
			args{
				map[string]nameservice.ReverseLookup{"ethereum/mainnet": nil},
			},
			func() *http.Request {
				req := httptest.NewRequest("GET", "/?network=mainnet&protocol=ethereum", nil)
				req = mux.SetURLVars(req, map[string]string{
					"address": "0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761",
				})
				return req
			}(),
			"{\"code\":422,\"message\":\"no nameserver configured for \\\"ethereum/mainnet\\\"\"}\n",
			http.StatusUnprocessableEntity,
		},
		{
			"networkNotSupportedError",
			args{
				func() map[string]nameservice.ReverseLookup {
					m := nameservicetest.NewMockReverseLookup(mockCtrl)
					m.EXPECT().ResolveAddress(gomock.Any(), "ethereum", "mainnet", encodingtest.MustDecodeHex("5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761")).Return("", errors.New("network not supported")).Times(1)
					return map[string]nameservice.ReverseLookup{"ethereum/mainnet": m}
				}(),
			},
			func() *http.Request {
				req := httptest.NewRequest("GET", "/?network=mainnet&protocol=ethereum", nil)
				req = mux.SetURLVars(req, map[string]string{
					"address": "0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761",
				})
				return req
			}(),
			"{\"code\":406,\"message\":\"\\\"ethereum/mainnet\\\" not supported\"}\n",
			http.StatusNotAcceptable,
		},
		{
			"err-resolve-name",
			args{
				func() map[string]nameservice.ReverseLookup {
					m := nameservicetest.NewMockReverseLookup(mockCtrl)
					m.EXPECT().ResolveAddress(gomock.Any(), "ethereum", "mainnet", encodingtest.MustDecodeHex("5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761")).Return("", errors.New("error")).Times(1)
					return map[string]nameservice.ReverseLookup{"ethereum/mainnet": m}
				}(),
			},
			func() *http.Request {
				req := httptest.NewRequest("GET", "/?network=mainnet&protocol=ethereum", nil)
				req = mux.SetURLVars(req, map[string]string{
					"address": "0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761",
				})
				return req
			}(),
			"{\"code\":500,\"message\":\"error\"}\n",
			http.StatusInternalServerError,
		},
		{
			"err-invalid-address",
			args{
				func() map[string]nameservice.ReverseLookup {
					m := nameservicetest.NewMockReverseLookup(mockCtrl)
					m.EXPECT().ResolveAddress(gomock.Any(), "ethereum", "mainnet", encodingtest.MustDecodeHex("5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761")).Return("", nameservice.ErrFormat).Times(1)
					return map[string]nameservice.ReverseLookup{"ethereum/mainnet": m}
				}(),
			},
			func() *http.Request {
				req := httptest.NewRequest("GET", "/?network=mainnet&protocol=ethereum", nil)
				req = mux.SetURLVars(req, map[string]string{
					"address": "0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761",
				})
				return req
			}(),
			"{\"name\":\"\",\"status\":1}\n",
			http.StatusOK,
		},
		{
			"err-nx-domain",
			args{
				func() map[string]nameservice.ReverseLookup {
					m := nameservicetest.NewMockReverseLookup(mockCtrl)
					m.EXPECT().ResolveAddress(gomock.Any(), "ethereum", "mainnet", encodingtest.MustDecodeHex("5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761")).Return("", nameservice.ErrNXDomain).Times(1)
					return map[string]nameservice.ReverseLookup{"ethereum/mainnet": m}
				}(),
			},
			func() *http.Request {
				req := httptest.NewRequest("GET", "/?network=mainnet&protocol=ethereum", nil)
				req = mux.SetURLVars(req, map[string]string{
					"address": "0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761",
				})
				return req
			}(),
			"{\"name\":\"\",\"status\":3}\n",
			http.StatusOK,
		},
		{
			"success",
			args{
				func() map[string]nameservice.ReverseLookup {
					m := nameservicetest.NewMockReverseLookup(mockCtrl)
					m.EXPECT().ResolveAddress(gomock.Any(), "ethereum", "mainnet", encodingtest.MustDecodeHex("5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761")).Return("name.ens", nil).Times(1)
					return map[string]nameservice.ReverseLookup{"ethereum/mainnet": m}
				}(),
			},
			func() *http.Request {
				req := httptest.NewRequest("GET", "/?network=mainnet&protocol=ethereum", nil)
				req = mux.SetURLVars(req, map[string]string{
					"address": "0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761",
				})
				return req
			}(),
			"{\"name\":\"name.ens\",\"status\":0}\n",
			http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(GetResolveAddress(tt.args.resolvers))

			// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
			// directly and pass in our Request and ResponseRecorder.
			handler.ServeHTTP(rr, tt.req)

			// Check the status code is what we expect.
			if !assert.Equal(t, tt.wantStatus, rr.Code) {
				t.Errorf("handler returned wrong status code: got %v want %v",
					rr.Code, tt.wantStatus)
			}
			if !assert.Equal(t, tt.wantBody, rr.Body.String()) {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), tt.wantBody)
			}
		})
	}
}
