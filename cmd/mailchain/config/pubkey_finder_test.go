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

func Test_getFinder(t *testing.T) {
	type args struct {
		vpr     *viper.Viper
		chain   string
		network string
	}
	tests := []struct {
		name    string
		args    args
		want    mailbox.PubKeyFinder
		wantErr bool
	}{
		{
			"etherscan",
			args{
				func() *viper.Viper {
					v := viper.New()
					v.Set("chains.ethereum.networks.mainnet.pubkey-finder", "etherscan")
					v.Set("clients.etherscan.api-key", "api-key-value")
					return v
				}(),
				"ethereum",
				"mainnet",
			},
			func() mailbox.PubKeyFinder {
				r, _ := etherscan.NewAPIClient("api-key-value")
				return r
			}(),
			false,
		},
		{
			"etherscan-auth",
			args{
				func() *viper.Viper {
					v := viper.New()
					v.Set("chains.ethereum.networks.mainnet.pubkey-finder", "etherscan-no-auth")
					return v
				}(),
				"ethereum",
				"mainnet",
			},
			func() mailbox.PubKeyFinder {
				r, _ := etherscan.NewAPIClient("")
				return r
			}(),
			false,
		},
		{
			"etherscan-auth",
			args{
				func() *viper.Viper {
					v := viper.New()
					v.Set("chains.ethereum.networks.mainnet.pubkey-finder", "unknown")
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
			got, err := getFinder(tt.args.vpr, tt.args.chain, tt.args.network)
			if (err != nil) != tt.wantErr {
				t.Errorf("getFinder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getFinder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getChainFinders(t *testing.T) {
	assert := assert.New(t)
	// is := is.New(t)
	type args struct {
		vpr   *viper.Viper
		chain string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]mailbox.PubKeyFinder
		wantErr bool
	}{
		{
			"empty",
			args{
				viper.New(),
				"ethereum",
			},
			make(map[string]mailbox.PubKeyFinder),
			false,
		},
		{
			"single",
			args{
				func() *viper.Viper {
					v := viper.New()
					v.Set("chains.ethereum.networks.mainnet.pubkey-finder", "etherscan-no-auth")
					return v
				}(),
				"ethereum",
			},
			func() map[string]mailbox.PubKeyFinder {
				c, _ := etherscan.NewAPIClient("")
				return map[string]mailbox.PubKeyFinder{"ethereum.mainnet": c}
			}(),
			false,
		},
		{
			"multi",
			args{
				func() *viper.Viper {
					v := viper.New()
					v.Set("chains.ethereum.networks.mainnet.pubkey-finder", "etherscan-no-auth")
					v.Set("chains.ethereum.networks.ropsten.pubkey-finder", "etherscan-no-auth")
					return v
				}(),
				"ethereum",
			},
			func() map[string]mailbox.PubKeyFinder {
				c, _ := etherscan.NewAPIClient("")
				return map[string]mailbox.PubKeyFinder{"ethereum.mainnet": c, "ethereum.ropsten": c}
			}(),
			false,
		},
		{
			"err",
			args{
				func() *viper.Viper {
					v := viper.New()
					v.Set("chains.ethereum.networks.mainnet.pubkey-finder", "etherscan-no-auth")
					v.Set("chains.ethereum.networks.ropsten.pubkey-finder", "unknown")
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
			got, err := getChainFinders(tt.args.vpr, tt.args.chain)
			if (err != nil) != tt.wantErr {
				t.Errorf("getChainFinders() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(got, tt.want) {
				t.Errorf("getChainFinders() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetPublicKeyFinders(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		vpr *viper.Viper
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]mailbox.PubKeyFinder
		wantErr bool
	}{
		{
			"empty",
			args{
				viper.New(),
			},
			make(map[string]mailbox.PubKeyFinder),
			false,
		},
		{
			"single",
			args{
				func() *viper.Viper {
					v := viper.New()
					v.Set("chains.ethereum.networks.mainnet.pubkey-finder", "etherscan-no-auth")
					return v
				}(),
			},
			func() map[string]mailbox.PubKeyFinder {
				c, _ := etherscan.NewAPIClient("")
				return map[string]mailbox.PubKeyFinder{"ethereum.mainnet": c}
			}(),
			false,
		},
		{
			"multi",
			args{
				func() *viper.Viper {
					v := viper.New()
					v.Set("chains.ethereum.networks.mainnet.pubkey-finder", "etherscan-no-auth")
					v.Set("chains.ethereum.networks.ropsten.pubkey-finder", "etherscan-no-auth")
					return v
				}(),
			},
			func() map[string]mailbox.PubKeyFinder {
				c, _ := etherscan.NewAPIClient("")
				return map[string]mailbox.PubKeyFinder{"ethereum.mainnet": c, "ethereum.ropsten": c}
			}(),
			false,
		},
		{
			"err",
			args{
				func() *viper.Viper {
					v := viper.New()
					v.Set("chains.ethereum.networks.mainnet.pubkey-finder", "etherscan-no-auth")
					v.Set("chains.ethereum.networks.ropsten.pubkey-finder", "unknown")
					return v
				}(),
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetPublicKeyFinders(tt.args.vpr)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPublicKeyFinders() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(tt.want, got) {
				t.Errorf("GetPublicKeyFinders() = %v, want %v", got, tt.want)
			}
		})
	}
}
