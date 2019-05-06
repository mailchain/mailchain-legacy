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
	"testing"

	"github.com/mailchain/mailchain/internal/pkg/clients/etherscan"
	"github.com/matryer/is"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func Test_setEtherscan(t *testing.T) {
	is := is.New(t)
	type args struct {
		vpr           *viper.Viper
		requiredInput func(label string) (string, error)
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		expected map[string]interface{}
	}{
		{
			"set-api-key",
			args{
				viper.New(),
				func(label string) (string, error) {
					return "api-key-value", nil
				},
			},
			false,
			map[string]interface{}{"api-key": "api-key-value"},
		},
		{
			"error",
			args{
				viper.New(),
				func(label string) (string, error) {
					return "", errors.Errorf("failed")
				},
			},
			true,
			nil,
		},
		{
			"already-specified",
			args{
				func() *viper.Viper {
					v := viper.New()
					v.Set("clients.etherscan.api-key", "already-set")
					return v
				}(),
				func(label string) (string, error) {
					return "api-key-value", nil
				},
			},
			false,
			map[string]interface{}{"api-key": "already-set"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := setEtherscan(tt.args.vpr, tt.args.requiredInput); (err != nil) != tt.wantErr {
				t.Errorf("setEtherscan() error = %v, wantErr %v", err, tt.wantErr)
			}

			is.Equal(tt.expected, tt.args.vpr.Get("clients.etherscan"))
		})
	}
}

func Test_setEthRPC(t *testing.T) {
	is := is.New(t)
	type args struct {
		vpr           *viper.Viper
		requiredInput func(label string) (string, error)
		network       string
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		expected map[string]interface{}
	}{
		{
			"set-address",
			args{
				viper.New(),
				func(label string) (string, error) {
					return "address-value", nil
				},
				"mainnet",
			},
			false,
			map[string]interface{}{"address": "address-value"},
		},
		{
			"error",
			args{
				viper.New(),
				func(label string) (string, error) {
					return "", errors.Errorf("failed")
				},
				"mainnet",
			},
			true,
			nil,
		},
		{
			"already-specified",
			args{
				func() *viper.Viper {
					v := viper.New()
					v.Set("clients.ethereum-rpc2.mainnet.address", "already-set")
					return v
				}(),
				func(label string) (string, error) {
					return "api-key-value", nil
				},
				"mainnet",
			},
			false,
			map[string]interface{}{"address": "already-set"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := setEthRPC(tt.args.vpr, tt.args.requiredInput, tt.args.network); (err != nil) != tt.wantErr {
				t.Errorf("setEthRPC() error = %v, wantErr %v", err, tt.wantErr)
			}

			is.Equal(tt.expected, tt.args.vpr.Get("clients.ethereum-rpc2.mainnet"))
		})
	}
}

func Test_getEtherscanClient(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		vpr *viper.Viper
	}
	tests := []struct {
		name    string
		args    args
		want    *etherscan.APIClient
		wantErr bool
	}{
		{
			"success",
			args{
				func() *viper.Viper {
					v := viper.New()
					v.Set("clients.etherscan.api-key", "apikey-value")
					return v
				}(),
			},
			func() *etherscan.APIClient {
				v, _ := etherscan.NewAPIClient("apikey-value")
				return v
			}(),
			false,
		},
		{
			"error",
			args{
				viper.New(),
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getEtherscanClient(tt.args.vpr)
			if (err != nil) != tt.wantErr {
				t.Errorf("getEtherscanClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(tt.want, got) {
				t.Errorf("getEtherscanClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getEtherRPC2Client(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		vpr     *viper.Viper
		network string
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
				func() *viper.Viper {
					v := viper.New()
					v.Set("clients.ethereum-rpc2.mainnet.address", "http://localhost:123423")
					return v
				}(),
				"mainnet",
			},
			false,
			false,
		},
		{
			"error",
			args{
				viper.New(),
				"mainnet",
			},
			true,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getEtherRPC2Client(tt.args.vpr, tt.args.network)
			if (err != nil) != tt.wantErr {
				t.Errorf("getEtherRPC2Client() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(tt.wantNil, (got == nil)) {
				t.Errorf("getEtherRPC2Client() = %v, want %v", got, tt.wantNil)
			}
		})
	}
}

func Test_getEtherscanNoAuthClient(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name    string
		want    *etherscan.APIClient
		wantErr bool
	}{
		{
			"success",
			func() *etherscan.APIClient {
				r, _ := etherscan.NewAPIClient("")
				return r
			}(),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getEtherscanNoAuthClient()
			if (err != nil) != tt.wantErr {
				t.Errorf("getEtherscanNoAuthClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(tt.want, got) {
				t.Errorf("getEtherscanNoAuthClient() = %v, want %v", got, tt.want)
			}
		})
	}
}
