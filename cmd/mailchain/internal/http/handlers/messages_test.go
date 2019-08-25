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
	"reflect"
	"testing"

	"github.com/gorilla/mux"
)

func Test_parseGetMessagesRequest(t *testing.T) {
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name    string
		args    args
		want    *GetMessagesRequest
		wantErr bool
	}{
		{
			"success",
			args{
				func() *http.Request {
					req := httptest.NewRequest("GET", "/", nil)
					req = mux.SetURLVars(req, map[string]string{
						"address":  "0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761",
						"network":  "mainnet",
						"protocol": "ethereum",
					})
					return req
				}(),
			},
			&GetMessagesRequest{Address: "0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761", Network: "mainnet"},
			false,
		},
		{
			"err-empty-address",
			args{
				func() *http.Request {
					req := httptest.NewRequest("GET", "/", nil)
					req = mux.SetURLVars(req, map[string]string{
						"address":  "",
						"network":  "mainnet",
						"protocol": "ethereum",
					})
					return req
				}(),
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseGetMessagesRequest(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseGetMessagesRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseGetMessagesRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}
