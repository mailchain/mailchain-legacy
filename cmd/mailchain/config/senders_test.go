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

package config

import (
	"net/http/httptest"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func Test_getSender(t *testing.T) {
	assert := assert.New(t)
	server := httptest.NewServer(nil)
	defer server.Close()
	type args struct {
		v       *viper.Viper
		chain   string
		network string
	}
	tests := []struct {
		name    string
		args    args
		wantNil bool
		wantErr bool
	}{
		{
			"etherscan",
			args{
				func() *viper.Viper {
					v := viper.New()
					v.Set("chains.ethereum.networks.mainnet.sender", "ethereum-rpc2")
					v.Set("clients.ethereum-rpc2.mainnet.address", server.URL)
					return v
				}(),
				"ethereum",
				"mainnet",
			},
			false,
			false,
		},
		{
			"error",
			args{
				func() *viper.Viper {
					v := viper.New()
					v.Set("chains.ethereum.networks.mainnet.sender", "unknown")
					return v
				}(),
				"ethereum",
				"mainnet",
			},
			true,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getSender(tt.args.v, tt.args.chain, tt.args.network)
			if (err != nil) != tt.wantErr {
				t.Errorf("getSender() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(tt.wantNil, got == nil) {
				t.Errorf("getSender() = %v, want %v", got, tt.wantNil)
			}
		})
	}
}

func Test_getChainSenders(t *testing.T) {
	assert := assert.New(t)
	server := httptest.NewServer(nil)
	type args struct {
		v     *viper.Viper
		chain string
	}
	tests := []struct {
		name    string
		args    args
		wantNil bool
		wantErr bool
	}{
		{
			"empty",
			args{
				viper.New(),
				"ethereum",
			},
			false,
			false,
		},
		{
			"single",
			args{
				func() *viper.Viper {
					v := viper.New()
					v.Set("chains.ethereum.networks.mainnet.sender", "ethereum-rpc2")
					v.Set("clients.ethereum-rpc2.mainnet.address", server.URL)
					return v
				}(),
				"ethereum",
			},
			false,
			false,
		},
		{
			"multi",
			args{
				func() *viper.Viper {
					v := viper.New()
					v.Set("chains.ethereum.networks.mainnet.sender", "ethereum-rpc2")
					v.Set("chains.ethereum.networks.ropsten.sender", "ethereum-rpc2")
					v.Set("clients.ethereum-rpc2.mainnet.address", server.URL)
					v.Set("clients.ethereum-rpc2.ropsten.address", server.URL)
					return v
				}(),
				"ethereum",
			},
			false,
			false,
		},
		{
			"err",
			args{
				func() *viper.Viper {
					v := viper.New()
					v.Set("chains.ethereum.networks.mainnet.sender", "ethereum-rpc2")
					v.Set("chains.ethereum.networks.ropsten.sender", "unknown")
					v.Set("clients.ethereum-rpc2.mainnet.address", server.URL)
					return v
				}(),
				"ethereum",
			},
			true,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getChainSenders(tt.args.v, tt.args.chain)
			if (err != nil) != tt.wantErr {
				t.Errorf("getChainSenders() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(tt.wantNil, got == nil) {
				t.Errorf("getChainSenders() = %v, want %v", got, tt.wantNil)
			}
		})
	}
}

func TestGetSenders(t *testing.T) {
	assert := assert.New(t)
	server := httptest.NewServer(nil)
	type args struct {
		v *viper.Viper
	}
	tests := []struct {
		name    string
		args    args
		wantNil bool
		wantErr bool
	}{
		{
			"empty",
			args{
				viper.New(),
			},
			false,
			false,
		},
		{
			"single",
			args{
				func() *viper.Viper {
					v := viper.New()
					v.Set("chains.ethereum.networks.mainnet.sender", "ethereum-rpc2")
					v.Set("clients.ethereum-rpc2.mainnet.address", server.URL)
					return v
				}(),
			},
			false,
			false,
		},
		{
			"multi",
			args{
				func() *viper.Viper {
					v := viper.New()
					v.Set("chains.ethereum.networks.mainnet.sender", "ethereum-rpc2")
					v.Set("chains.ethereum.networks.ropsten.sender", "ethereum-rpc2")
					v.Set("clients.ethereum-rpc2.mainnet.address", server.URL)
					v.Set("clients.ethereum-rpc2.ropsten.address", server.URL)
					return v
				}(),
			},
			false,
			false,
		},
		{
			"err",
			args{
				func() *viper.Viper {
					v := viper.New()
					v.Set("chains.ethereum.networks.mainnet.sender", "ethereum-rpc2")
					v.Set("chains.ethereum.networks.ropsten.sender", "unknown")
					v.Set("clients.ethereum-rpc2.mainnet.address", server.URL)
					return v
				}(),
			},
			true,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetSenders(tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSenders() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(tt.wantNil, got == nil) {
				t.Errorf("GetSenders() = %v, want %v", got, tt.wantNil)
			}
		})
	}
}

func TestSetSender(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		v       *viper.Viper
		chain   string
		network string
		sender  string
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		expected map[string]interface{}
	}{
		{
			"success",
			args{
				viper.New(),
				"ethereum",
				"mainnet",
				"etherscan-no-auth",
			},
			false,
			map[string]interface{}{"sender": "etherscan-no-auth"},
		},
		{
			"error",
			args{
				viper.New(),
				"ethereum",
				"mainnet",
				"invalid",
			},
			true,
			map[string]interface{}{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SetSender(tt.args.v, tt.args.chain, tt.args.network, tt.args.sender); (err != nil) != tt.wantErr {
				t.Errorf("SetSender() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !assert.EqualValues(tt.expected, tt.args.v.GetStringMap("chains.ethereum.networks.mainnet")) {
				t.Errorf("SetSender() = expected")
			}
		})
	}
}
