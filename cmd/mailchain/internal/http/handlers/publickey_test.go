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

	"github.com/gorilla/mux"
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
