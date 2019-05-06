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

// nolint: dupl
package config

import (
	"reflect"
	"testing"

	"github.com/mailchain/mailchain/internal/pkg/clients/etherscan"
	"github.com/mailchain/mailchain/internal/pkg/mailbox"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func Test_getReceiver(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		vpr     *viper.Viper
		chain   string
		network string
	}
	tests := []struct {
		name    string
		args    args
		want    mailbox.Receiver
		wantErr bool
	}{
		{
			"etherscan",
			args{
				func() *viper.Viper {
					v := viper.New()
					v.Set("chains.ethereum.networks.mainnet.receiver", "etherscan")
					v.Set("clients.etherscan.api-key", "api-key-value")
					return v
				}(),
				"ethereum",
				"mainnet",
			},
			func() mailbox.Receiver {
				r, _ := etherscan.NewAPIClient("api-key-value")
				return r
			}(),
			false,
		},
		{
			"etherscan-no-auth",
			args{
				func() *viper.Viper {
					v := viper.New()
					v.Set("chains.ethereum.networks.mainnet.receiver", "etherscan-no-auth")
					return v
				}(),
				"ethereum",
				"mainnet",
			},
			func() mailbox.Receiver {
				r, _ := etherscan.NewAPIClient("")
				return r
			}(),
			false,
		},
		{
			"error",
			args{
				func() *viper.Viper {
					v := viper.New()
					v.Set("chains.ethereum.networks.mainnet.receiver", "unknown")
					return v
				}(),
				"ethereum",
				"mainnet",
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getReceiver(tt.args.vpr, tt.args.chain, tt.args.network)
			if (err != nil) != tt.wantErr {
				t.Errorf("getReceiver() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(tt.want, got) {
				t.Errorf("getReceiver() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getChainReceivers(t *testing.T) {
	type args struct {
		vpr   *viper.Viper
		chain string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]mailbox.Receiver
		wantErr bool
	}{
		{
			"empty",
			args{
				viper.New(),
				"ethereum",
			},
			make(map[string]mailbox.Receiver),
			false,
		},
		{
			"single",
			args{
				func() *viper.Viper {
					v := viper.New()
					v.Set("chains.ethereum.networks.mainnet.receiver", "etherscan-no-auth")
					return v
				}(),
				"ethereum",
			},
			func() map[string]mailbox.Receiver {
				c, _ := etherscan.NewAPIClient("")
				return map[string]mailbox.Receiver{"ethereum.mainnet": c}
			}(),
			false,
		},
		{
			"multi",
			args{
				func() *viper.Viper {
					v := viper.New()
					v.Set("chains.ethereum.networks.mainnet.receiver", "etherscan-no-auth")
					v.Set("chains.ethereum.networks.ropsten.receiver", "etherscan-no-auth")
					return v
				}(),
				"ethereum",
			},
			func() map[string]mailbox.Receiver {
				c, _ := etherscan.NewAPIClient("")
				return map[string]mailbox.Receiver{"ethereum.mainnet": c, "ethereum.ropsten": c}
			}(),
			false,
		},
		{
			"err",
			args{
				func() *viper.Viper {
					v := viper.New()
					v.Set("chains.ethereum.networks.mainnet.receiver", "etherscan-no-auth")
					v.Set("chains.ethereum.networks.ropsten.receiver", "unknown")
					return v
				}(),
				"ethereum",
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getChainReceivers(tt.args.vpr, tt.args.chain)
			if (err != nil) != tt.wantErr {
				t.Errorf("getChainReceivers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getChainReceivers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetReceivers(t *testing.T) {
	type args struct {
		vpr *viper.Viper
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]mailbox.Receiver
		wantErr bool
	}{
		{
			"empty",
			args{
				viper.New(),
			},
			make(map[string]mailbox.Receiver),
			false,
		},
		{
			"single",
			args{
				func() *viper.Viper {
					v := viper.New()
					v.Set("chains.ethereum.networks.mainnet.receiver", "etherscan-no-auth")
					return v
				}(),
			},
			func() map[string]mailbox.Receiver {
				c, _ := etherscan.NewAPIClient("")
				return map[string]mailbox.Receiver{"ethereum.mainnet": c}
			}(),
			false,
		},
		{
			"multi",
			args{
				func() *viper.Viper {
					v := viper.New()
					v.Set("chains.ethereum.networks.mainnet.receiver", "etherscan-no-auth")
					v.Set("chains.ethereum.networks.ropsten.receiver", "etherscan-no-auth")
					return v
				}(),
			},
			func() map[string]mailbox.Receiver {
				c, _ := etherscan.NewAPIClient("")
				return map[string]mailbox.Receiver{"ethereum.mainnet": c, "ethereum.ropsten": c}
			}(),
			false,
		},
		{
			"err",
			args{
				func() *viper.Viper {
					v := viper.New()
					v.Set("chains.ethereum.networks.mainnet.receiver", "etherscan-no-auth")
					v.Set("chains.ethereum.networks.ropsten.receiver", "unknown")
					return v
				}(),
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetReceivers(tt.args.vpr)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetReceivers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetReceivers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetReceiver(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		v        *viper.Viper
		chain    string
		network  string
		receiver string
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
				func() *viper.Viper {
					v := viper.New()
					return v
				}(),
				"ethereum",
				"mainnet",
				"etherscan-no-auth",
			},
			false,
			map[string]interface{}{"receiver": "etherscan-no-auth"},
		},
		{
			"error",
			args{
				func() *viper.Viper {
					v := viper.New()
					return v
				}(),
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
			if err := SetReceiver(tt.args.v, tt.args.chain, tt.args.network, tt.args.receiver); (err != nil) != tt.wantErr {
				t.Errorf("SetReceiver() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !assert.EqualValues(tt.expected, tt.args.v.GetStringMap("chains.ethereum.networks.mainnet")) {
				t.Errorf("getChainFinders() = expected")
			}
		})
	}
}
