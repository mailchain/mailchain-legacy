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
	"reflect"
	"testing"

	"github.com/mailchain/mailchain/internal/clients/etherscan"
	"github.com/mailchain/mailchain/internal/clients/ethrpc"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestClients_setEtherscan(t *testing.T) {
	assert := assert.New(t)
	type fields struct {
		viper         *viper.Viper
		requiredInput func(label string) (string, error)
	}
	tests := []struct {
		name         string
		fields       fields
		wantErr      bool
		wantSettings map[string]interface{}
	}{
		{
			"set-api-key",
			fields{
				viper.New(),
				func(label string) (string, error) {
					return "api-key-value", nil
				},
			},
			false,
			map[string]interface{}{
				"clients": map[string]interface{}{
					"etherscan": map[string]interface{}{
						"api-key": "api-key-value"}}},
		},
		{
			"error",
			fields{
				viper.New(),
				func(label string) (string, error) {
					return "", errors.Errorf("failed")
				},
			},
			true,
			map[string]interface{}{},
		},
		{
			"already-specified",
			fields{
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
			map[string]interface{}{
				"clients": map[string]interface{}{
					"etherscan": map[string]interface{}{
						"api-key": "already-set"}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Clients{
				viper:         tt.fields.viper,
				requiredInput: tt.fields.requiredInput,
			}
			if err := c.setEtherscan(); (err != nil) != tt.wantErr {
				t.Errorf("Clients.setEtherscan() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !assert.Equal(tt.wantSettings, tt.fields.viper.AllSettings()) {
				t.Errorf("settings = %v, wantSettings %v", tt.fields.viper.AllSettings(), tt.wantSettings)
			}
		})
	}
}

func TestClients_setEthRPC(t *testing.T) {
	assert := assert.New(t)
	type fields struct {
		viper         *viper.Viper
		requiredInput func(label string) (string, error)
	}
	type args struct {
		network string
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantErr      bool
		wantSettings map[string]interface{}
	}{
		{
			"set-address",
			fields{
				viper.New(),
				func(label string) (string, error) {
					return "address-value", nil
				},
			},
			args{
				"mainnet",
			},
			false,
			map[string]interface{}{
				"clients": map[string]interface{}{
					"ethereum-rpc2": map[string]interface{}{
						"mainnet": map[string]interface{}{
							"address": "address-value"}}}},
		},
		{
			"error",
			fields{
				viper.New(),
				func(label string) (string, error) {
					return "", errors.Errorf("failed")
				},
			},
			args{
				"mainnet",
			},
			true,
			map[string]interface{}{},
		},
		{
			"already-specified",
			fields{
				func() *viper.Viper {
					v := viper.New()
					v.Set("clients.ethereum-rpc2.mainnet.address", "already-set")
					return v
				}(),
				func(label string) (string, error) {
					return "api-key-value", nil
				},
			},
			args{
				"mainnet",
			},
			false,
			map[string]interface{}{
				"clients": map[string]interface{}{
					"ethereum-rpc2": map[string]interface{}{
						"mainnet": map[string]interface{}{
							"address": "already-set"}}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Clients{
				viper:         tt.fields.viper,
				requiredInput: tt.fields.requiredInput,
			}
			if err := c.setEthRPC(tt.args.network); (err != nil) != tt.wantErr {
				t.Errorf("Clients.setEthRPC() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !assert.Equal(tt.wantSettings, tt.fields.viper.AllSettings()) {
				t.Errorf("settings = %v, wantSettings %v", tt.fields.viper.AllSettings(), tt.wantSettings)
			}
		})
	}
}

func TestClients_SetClient(t *testing.T) {
	assert := assert.New(t)
	type fields struct {
		viper         *viper.Viper
		requiredInput func(label string) (string, error)
	}
	type args struct {
		client  string
		network string
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantErr      bool
		wantSettings map[string]interface{}
	}{
		{
			"ethereum-rpc2",
			fields{
				viper.New(),
				func(label string) (string, error) {
					return "address-value", nil
				},
			},
			args{
				"ethereum-rpc2",
				"mainnet",
			},
			false,
			map[string]interface{}{
				"clients": map[string]interface{}{
					"ethereum-rpc2": map[string]interface{}{
						"mainnet": map[string]interface{}{
							"address": "address-value"}}}},
		},
		{
			"etherscan",
			fields{
				viper.New(),
				func(label string) (string, error) {
					return "api-key-value", nil
				},
			},
			args{
				"etherscan",
				"mainnet",
			},
			false,
			map[string]interface{}{
				"clients": map[string]interface{}{
					"etherscan": map[string]interface{}{
						"api-key": "api-key-value"}}},
		},
		{
			"etherscan-no-auth",
			fields{
				viper.New(),
				func(label string) (string, error) {
					return "api-key-value", nil
				},
			},
			args{
				"etherscan-no-auth",
				"mainnet",
			},
			false,
			map[string]interface{}{},
		},
		{
			"err-unknown",
			fields{
				viper.New(),
				func(label string) (string, error) {
					return "api-key-value", nil
				},
			},
			args{
				"unknown",
				"mainnet",
			},
			true,
			map[string]interface{}{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Clients{
				viper:         tt.fields.viper,
				requiredInput: tt.fields.requiredInput,
			}
			if err := c.SetClient(tt.args.client, tt.args.network); (err != nil) != tt.wantErr {
				t.Errorf("Clients.SetClient() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !assert.Equal(tt.wantSettings, tt.fields.viper.AllSettings()) {
				t.Errorf("settings = %v, wantSettings %v", tt.fields.viper.AllSettings(), tt.wantSettings)
			}
		})
	}
}

func TestClients_GetEtherscanNoAuthClient(t *testing.T) {
	type fields struct {
		viper         *viper.Viper
		requiredInput func(label string) (string, error)
	}
	tests := []struct {
		name    string
		fields  fields
		want    *etherscan.APIClient
		wantErr bool
	}{
		{
			"success",
			fields{
				nil,
				nil,
			},
			func() *etherscan.APIClient {
				r, _ := etherscan.NewAPIClient("")
				return r
			}(),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Clients{
				viper:         tt.fields.viper,
				requiredInput: tt.fields.requiredInput,
			}
			got, err := c.GetEtherscanNoAuthClient()
			if (err != nil) != tt.wantErr {
				t.Errorf("Clients.GetEtherscanNoAuthClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Clients.GetEtherscanNoAuthClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClients_GetEtherRPC2Client(t *testing.T) {
	assert := assert.New(t)
	type fields struct {
		viper         *viper.Viper
		requiredInput func(label string) (string, error)
	}
	type args struct {
		network string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *ethrpc.EthRPC2
		wantErr bool
	}{
		{
			"success",
			fields{
				func() *viper.Viper {
					v := viper.New()
					v.Set("clients.ethereum-rpc2.mainnet.address", "http://localhost:123423")
					return v
				}(),
				nil,
			},
			args{
				"mainnet",
			},
			&ethrpc.EthRPC2{},
			false,
		},
		{
			"error",
			fields{
				viper.New(),
				nil,
			},
			args{
				"mainnet",
			},
			&ethrpc.EthRPC2{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Clients{
				viper:         tt.fields.viper,
				requiredInput: tt.fields.requiredInput,
			}
			got, err := c.GetEtherRPC2Client(tt.args.network)
			if (err != nil) != tt.wantErr {
				t.Errorf("Clients.GetEtherRPC2Client() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.IsType(tt.want, got) {
				t.Errorf("Clients.GetEtherRPC2Client() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClients_GetEtherscanClient(t *testing.T) {
	type fields struct {
		viper         *viper.Viper
		requiredInput func(label string) (string, error)
	}
	tests := []struct {
		name    string
		fields  fields
		want    *etherscan.APIClient
		wantErr bool
	}{
		{
			"success",
			fields{
				func() *viper.Viper {
					v := viper.New()
					v.Set("clients.etherscan.api-key", "apikey-value")
					return v
				}(),
				nil,
			},
			func() *etherscan.APIClient {
				v, _ := etherscan.NewAPIClient("apikey-value")
				return v
			}(),
			false,
		},
		{
			"error",
			fields{
				viper.New(),
				nil,
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Clients{
				viper:         tt.fields.viper,
				requiredInput: tt.fields.requiredInput,
			}
			got, err := c.GetEtherscanClient()
			if (err != nil) != tt.wantErr {
				t.Errorf("Clients.GetEtherscanClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Clients.GetEtherscanClient() = %v, want %v", got, tt.want)
			}
		})
	}
}
