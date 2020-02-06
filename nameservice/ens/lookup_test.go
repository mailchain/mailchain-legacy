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

package ens

import (
	"context"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/mailchain/mailchain/encoding/encodingtest"
	"github.com/stretchr/testify/assert"
)

func TestNewLookupService(t *testing.T) {
	type args struct {
		clientURL string
	}
	tests := []struct {
		name    string
		args    args
		wantNil bool
		wantErr bool
	}{
		{
			"success",
			args{
				"https://client.url",
			},
			false,
			false,
		},
		{
			"err",
			args{
				"/client.url",
			},
			true,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewLookupService(tt.args.clientURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewLookupService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got == nil) != tt.wantNil {
				t.Errorf("NewLookupService() = %v, want %v", got, tt.wantNil)
			}
		})
	}
}

func TestLookupService_ResolveName(t *testing.T) {
	server := httptest.NewServer(nil)
	defer server.Close()
	type fields struct {
		client *ethclient.Client
	}
	type args struct {
		ctx        context.Context
		protocol   string
		network    string
		domainName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		{
			"err-client",
			fields{
				func() *ethclient.Client {
					c, _ := ethclient.Dial(server.URL)
					return c
				}(),
			},
			args{
				context.Background(),
				"ethereum",
				"mainnet",
				"test.eth",
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := LookupService{
				client: tt.fields.client,
			}
			got, err := s.ResolveName(tt.args.ctx, tt.args.protocol, tt.args.network, tt.args.domainName)
			if (err != nil) != tt.wantErr {
				t.Errorf("LookupService.ResolveName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LookupService.ResolveName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLookupService_ResolveAddress(t *testing.T) {
	server := httptest.NewServer(nil)
	defer server.Close()
	type fields struct {
		client *ethclient.Client
	}
	type args struct {
		ctx      context.Context
		protocol string
		network  string
		address  []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			"err-client",
			fields{
				func() *ethclient.Client {
					c, _ := ethclient.Dial(server.URL)
					return c
				}(),
			},
			args{
				context.Background(),
				"ethereum",
				"mainnet",
				encodingtest.MustDecodeHex("5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761"),
			},
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := LookupService{
				client: tt.fields.client,
			}
			got, err := s.ResolveAddress(tt.args.ctx, tt.args.protocol, tt.args.network, tt.args.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("LookupService.ResolveName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("LookupService.ResolveName() = %v, want %v", got, tt.want)
			}
		})
	}
}
