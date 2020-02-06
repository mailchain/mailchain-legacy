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
	"github.com/mailchain/mailchain/encoding/encodingtest"
	"github.com/mailchain/mailchain/nameservice"
	"github.com/mailchain/mailchain/nameservice/nameservicetest"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func Test_parseGetResolveNameRequest(t *testing.T) {
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name         string
		args         args
		wantProtocol string
		wantNetwork  string
		wantDomain   string
		wantErr      bool
	}{
		{
			"success",
			args{
				func() *http.Request {
					req := httptest.NewRequest("GET", "/?network=mainnet&protocol=ethereum", nil)
					req = mux.SetURLVars(req, map[string]string{
						"domain-name": "address.ens",
					})
					return req
				}(),
			},
			"ethereum",
			"mainnet",
			"address.ens",
			false,
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
			"",
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
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotProtocol, gotNetwork, gotDomain, err := parseGetResolveNameRequest(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseGetResolveNameRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotProtocol != tt.wantProtocol {
				t.Errorf("parseGetResolveNameRequest() gotProtocol = %v, want %v", gotProtocol, tt.wantProtocol)
			}
			if gotNetwork != tt.wantNetwork {
				t.Errorf("parseGetResolveNameRequest() gotNetwork = %v, want %v", gotNetwork, tt.wantNetwork)
			}
			if gotDomain != tt.wantDomain {
				t.Errorf("parseGetResolveNameRequest() gotDomain = %v, want %v", gotDomain, tt.wantDomain)
			}
		})
	}
}

func TestGetResolveName(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type args struct {
		resolvers map[string]nameservice.ForwardLookup
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
				func() map[string]nameservice.ForwardLookup {
					m := nameservicetest.NewMockForwardLookup(mockCtrl)
					return map[string]nameservice.ForwardLookup{"ethereum.no-network": m}
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
				map[string]nameservice.ForwardLookup{"ethereum/mainnet": nil},
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
				func() map[string]nameservice.ForwardLookup {
					m := nameservicetest.NewMockForwardLookup(mockCtrl)
					m.EXPECT().ResolveName(gomock.Any(), "ethereum", "mainnet", "name.ens").Return(nil, errors.New("network not supported")).Times(1)
					return map[string]nameservice.ForwardLookup{"ethereum/mainnet": m}
				}(),
			},
			func() *http.Request {
				req := httptest.NewRequest("GET", "/?network=mainnet&protocol=ethereum", nil)
				req = mux.SetURLVars(req, map[string]string{
					"domain-name": "name.ens",
				})
				return req
			}(),
			"{\"code\":406,\"message\":\"\\\"ethereum/mainnet\\\" not supported\"}\n",
			http.StatusNotAcceptable,
		},
		{
			"err-resolve-name",
			args{
				func() map[string]nameservice.ForwardLookup {
					m := nameservicetest.NewMockForwardLookup(mockCtrl)
					m.EXPECT().ResolveName(gomock.Any(), "ethereum", "mainnet", "name.ens").Return(nil, errors.New("error")).Times(1)
					return map[string]nameservice.ForwardLookup{"ethereum/mainnet": m}
				}(),
			},
			func() *http.Request {
				req := httptest.NewRequest("GET", "/?network=mainnet&protocol=ethereum", nil)
				req = mux.SetURLVars(req, map[string]string{
					"domain-name": "name.ens",
				})
				return req
			}(),
			"{\"code\":500,\"message\":\"error\"}\n",
			http.StatusInternalServerError,
		},
		{
			"err-encoding-error",
			args{
				func() map[string]nameservice.ForwardLookup {
					m := nameservicetest.NewMockForwardLookup(mockCtrl)
					m.EXPECT().ResolveName(gomock.Any(), "invalid", "mainnet", "name.ens").Return(encodingtest.MustDecodeHex("5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761"), nil).Times(1)
					return map[string]nameservice.ForwardLookup{"invalid/mainnet": m}
				}(),
			},
			func() *http.Request {
				req := httptest.NewRequest("GET", "/?network=mainnet&protocol=invalid", nil)
				req = mux.SetURLVars(req, map[string]string{
					"domain-name": "name.ens",
				})
				return req
			}(),
			`{"code":500,"message":"failed to encode address: \"invalid\" unsupported protocol"}` + "\n",
			http.StatusInternalServerError,
		},
		{
			"err-nx-domain",
			args{
				func() map[string]nameservice.ForwardLookup {
					m := nameservicetest.NewMockForwardLookup(mockCtrl)
					m.EXPECT().ResolveName(gomock.Any(), "ethereum", "mainnet", "name.ens").Return(nil, nameservice.ErrNXDomain).Times(1)
					return map[string]nameservice.ForwardLookup{"ethereum/mainnet": m}
				}(),
			},
			func() *http.Request {
				req := httptest.NewRequest("GET", "/?network=mainnet&protocol=ethereum", nil)
				req = mux.SetURLVars(req, map[string]string{
					"domain-name": "name.ens",
				})
				return req
			}(),
			"{\"address\":\"\",\"status\":3}\n",
			http.StatusOK,
		},
		{
			"success",
			args{
				func() map[string]nameservice.ForwardLookup {
					m := nameservicetest.NewMockForwardLookup(mockCtrl)
					m.EXPECT().ResolveName(gomock.Any(), "ethereum", "mainnet", "name.ens").Return(encodingtest.MustDecodeHex("5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761"), nil).Times(1)
					return map[string]nameservice.ForwardLookup{"ethereum/mainnet": m}
				}(),
			},
			func() *http.Request {
				req := httptest.NewRequest("GET", "/?network=mainnet&protocol=ethereum", nil)
				req = mux.SetURLVars(req, map[string]string{
					"domain-name": "name.ens",
				})
				return req
			}(),
			"{\"address\":\"0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761\",\"status\":0}\n",
			http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(GetResolveName(tt.args.resolvers))

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
