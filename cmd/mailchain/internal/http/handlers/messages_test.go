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
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parseGetMessagesRequest(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		queryParams map[string]string
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
				map[string]string{
					"address":  "0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761",
					"network":  "mainnet",
					"protocol": "ethereum",
				},
			},
			&GetMessagesRequest{
				Address:      "0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761",
				Network:      "mainnet",
				Protocol:     "ethereum",
				addressBytes: []uint8{0x56, 0x2, 0xea, 0x95, 0x54, 0xb, 0xee, 0x46, 0xd0, 0x3b, 0xa3, 0x35, 0xee, 0xd6, 0xf4, 0x9d, 0x11, 0x7e, 0xab, 0x95, 0xc8, 0xab, 0x8b, 0x71, 0xba, 0xe2, 0xcd, 0xd1, 0xe5, 0x64, 0xa7, 0x61},
			},
			false,
		},
		{
			"err-empty-address",
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
			"err-empty-protocol",
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
			"err-empty-network",
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
			got, err := parseGetMessagesRequest(req)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseGetMessagesRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(tt.want, got) {
				t.Errorf("parseGetMessagesRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}
