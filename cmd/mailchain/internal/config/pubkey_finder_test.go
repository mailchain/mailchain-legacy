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
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/imdario/mergo"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/config/configtest"
	"github.com/mailchain/mailchain/internal/clients/etherscan"
	"github.com/mailchain/mailchain/internal/mailbox"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestPubKeyFinder_getFinder(t *testing.T) {
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type fields struct {
		viper        *viper.Viper
		clientGetter ClientsGetter
		clientSetter ClientsSetter
		mapMerge     func(dst interface{}, src interface{}, opts ...func(*mergo.Config)) error
	}
	type args struct {
		chain   string
		network string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    mailbox.PubKeyFinder
		wantErr bool
	}{
		{
			"etherscan",
			fields{
				func() *viper.Viper {
					v := viper.New()
					v.Set("chains.ethereum.networks.mainnet.pubkey-finder", "etherscan")
					v.Set("clients.etherscan.api-key", "api-key-value")
					return v
				}(),
				func() ClientsGetter {
					g := configtest.NewMockClientsGetter(mockCtrl)
					g.EXPECT().GetEtherscanClient().Return(etherscan.NewAPIClient("api-key-value"))
					return g
				}(),
				nil,
				nil,
			},
			args{
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
			"etherscan-no-auth",
			fields{
				func() *viper.Viper {
					v := viper.New()
					v.Set("chains.ethereum.networks.mainnet.pubkey-finder", "etherscan-no-auth")
					return v
				}(),
				func() ClientsGetter {
					g := configtest.NewMockClientsGetter(mockCtrl)
					g.EXPECT().GetEtherscanNoAuthClient().Return(etherscan.NewAPIClient(""))
					return g
				}(),
				nil,
				nil,
			},
			args{
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
			"error",
			fields{
				func() *viper.Viper {
					v := viper.New()
					v.Set("chains.ethereum.networks.mainnet.pubkey-finder", "unknown")
					return v
				}(),
				func() ClientsGetter {
					g := configtest.NewMockClientsGetter(mockCtrl)
					return g
				}(),
				nil,
				nil,
			},
			args{
				"ethereum",
				"mainnet",
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := PubKeyFinder{
				viper:        tt.fields.viper,
				clientGetter: tt.fields.clientGetter,
				clientSetter: tt.fields.clientSetter,
				mapMerge:     tt.fields.mapMerge,
			}
			got, err := p.getFinder(tt.args.chain, tt.args.network)
			if (err != nil) != tt.wantErr {
				t.Errorf("PubKeyFinder.getFinder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(tt.want, got) {
				t.Errorf("PubKeyFinder.getFinder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPubKeyFinder_getChainFinders(t *testing.T) {
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type fields struct {
		viper        *viper.Viper
		clientGetter ClientsGetter
		clientSetter ClientsSetter
		mapMerge     func(dst interface{}, src interface{}, opts ...func(*mergo.Config)) error
	}
	type args struct {
		chain string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]mailbox.PubKeyFinder
		wantErr bool
	}{
		{
			"empty",
			fields{
				viper.New(),
				func() ClientsGetter {
					g := configtest.NewMockClientsGetter(mockCtrl)
					return g
				}(),
				nil,
				nil,
			},
			args{
				"ethereum",
			},
			make(map[string]mailbox.PubKeyFinder),
			false,
		},
		{
			"single",
			fields{
				func() *viper.Viper {
					v := viper.New()
					v.Set("chains.ethereum.networks.mainnet.pubkey-finder", "etherscan-no-auth")
					return v
				}(),
				func() ClientsGetter {
					g := configtest.NewMockClientsGetter(mockCtrl)
					g.EXPECT().GetEtherscanNoAuthClient().Return(etherscan.NewAPIClient(""))
					return g
				}(),
				nil,
				nil,
			},
			args{
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
			fields{
				func() *viper.Viper {
					v := viper.New()
					v.Set("chains.ethereum.networks.mainnet.pubkey-finder", "etherscan-no-auth")
					v.Set("chains.ethereum.networks.ropsten.pubkey-finder", "etherscan")
					v.Set("clients.etherscan.api-key", "api-key-value")
					return v
				}(),
				func() ClientsGetter {
					g := configtest.NewMockClientsGetter(mockCtrl)
					g.EXPECT().GetEtherscanNoAuthClient().Return(etherscan.NewAPIClient(""))
					g.EXPECT().GetEtherscanClient().Return(etherscan.NewAPIClient("api-key-value"))
					return g
				}(),
				nil,
				nil,
			},
			args{
				"ethereum",
			},
			func() map[string]mailbox.PubKeyFinder {
				cNoAuth, _ := etherscan.NewAPIClient("")
				cAuth, _ := etherscan.NewAPIClient("api-key-value")
				return map[string]mailbox.PubKeyFinder{"ethereum.mainnet": cNoAuth, "ethereum.ropsten": cAuth}
			}(),
			false,
		},
		{
			"err",
			fields{
				func() *viper.Viper {
					v := viper.New()
					// v.Set("chains.ethereum.networks.mainnet.pubkey-finder", "etherscan-no-auth")
					v.Set("chains.ethereum.networks.ropsten.pubkey-finder", "unknown")
					return v
				}(),
				func() ClientsGetter {
					g := configtest.NewMockClientsGetter(mockCtrl)
					// g.EXPECT().GetEtherscanNoAuthClient().Return(etherscan.NewAPIClient(""))
					return g
				}(),
				nil,
				nil,
			},
			args{
				"ethereum",
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := PubKeyFinder{
				viper:        tt.fields.viper,
				clientGetter: tt.fields.clientGetter,
				clientSetter: tt.fields.clientSetter,
				mapMerge:     tt.fields.mapMerge,
			}
			got, err := p.getChainFinders(tt.args.chain)
			if (err != nil) != tt.wantErr {
				t.Errorf("PubKeyFinder.getChainFinders() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(tt.want, got) {
				t.Errorf("PubKeyFinder.getChainFinders() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPubKeyFinder_GetPublicKeyFinders(t *testing.T) {
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type fields struct {
		viper        *viper.Viper
		clientGetter ClientsGetter
		clientSetter ClientsSetter
		mapMerge     func(dst interface{}, src interface{}, opts ...func(*mergo.Config)) error
	}
	tests := []struct {
		name    string
		fields  fields
		want    map[string]mailbox.PubKeyFinder
		wantErr bool
	}{
		{
			"empty",
			fields{
				func() *viper.Viper {
					v := viper.New()
					return v
				}(),
				func() ClientsGetter {
					g := configtest.NewMockClientsGetter(mockCtrl)
					return g
				}(),
				nil,
				nil,
			},
			make(map[string]mailbox.PubKeyFinder),
			false,
		},
		{
			"single",
			fields{
				func() *viper.Viper {
					v := viper.New()
					v.Set("chains.ethereum.networks.mainnet.pubkey-finder", "etherscan-no-auth")
					return v
				}(),
				func() ClientsGetter {
					g := configtest.NewMockClientsGetter(mockCtrl)
					g.EXPECT().GetEtherscanNoAuthClient().Return(etherscan.NewAPIClient(""))
					return g
				}(),
				nil,
				mergo.Merge,
			},
			func() map[string]mailbox.PubKeyFinder {
				c, _ := etherscan.NewAPIClient("")
				return map[string]mailbox.PubKeyFinder{"ethereum.mainnet": c}
			}(),
			false,
		},
		{
			"multi",
			fields{
				func() *viper.Viper {
					v := viper.New()
					v.Set("chains.ethereum.networks.mainnet.pubkey-finder", "etherscan-no-auth")
					v.Set("chains.ethereum.networks.ropsten.pubkey-finder", "etherscan-no-auth")
					return v
				}(),
				func() ClientsGetter {
					g := configtest.NewMockClientsGetter(mockCtrl)
					g.EXPECT().GetEtherscanNoAuthClient().Return(etherscan.NewAPIClient("")).Times(2)
					return g
				}(),
				nil,
				mergo.Merge,
			},
			func() map[string]mailbox.PubKeyFinder {
				c, _ := etherscan.NewAPIClient("")
				return map[string]mailbox.PubKeyFinder{"ethereum.mainnet": c, "ethereum.ropsten": c}
			}(),
			false,
		},
		{
			"err-invalid-finder",
			fields{
				func() *viper.Viper {
					v := viper.New()
					v.Set("chains.ethereum.networks.mainnet.pubkey-finder", "etherscan-no-auth")
					v.Set("chains.ethereum.networks.ropsten.pubkey-finder", "unknown")
					return v
				}(),
				func() ClientsGetter {
					g := configtest.NewMockClientsGetter(mockCtrl)
					g.EXPECT().GetEtherscanNoAuthClient().Return(etherscan.NewAPIClient(""))
					return g
				}(),
				nil,
				mergo.Merge,
			},
			nil,
			true,
		},
		{
			"err-merge",
			fields{
				func() *viper.Viper {
					v := viper.New()
					v.Set("chains.ethereum.networks.mainnet.pubkey-finder", "etherscan-no-auth")
					v.Set("chains.ethereum.networks.ropsten.pubkey-finder", "etherscan-no-auth")
					return v
				}(),
				func() ClientsGetter {
					g := configtest.NewMockClientsGetter(mockCtrl)
					g.EXPECT().GetEtherscanNoAuthClient().Return(etherscan.NewAPIClient("")).Times(2)
					return g
				}(),
				nil,
				func(dst interface{}, src interface{}, opts ...func(*mergo.Config)) error {
					return errors.Errorf("merge failed")
				},
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := PubKeyFinder{
				viper:        tt.fields.viper,
				clientGetter: tt.fields.clientGetter,
				clientSetter: tt.fields.clientSetter,
				mapMerge:     tt.fields.mapMerge,
			}
			got, err := p.GetFinders()
			if (err != nil) != tt.wantErr {
				t.Errorf("PubKeyFinder.GetFinders() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(tt.want, got) {
				t.Errorf("PubKeyFinder.GetFinders() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPubKeyFinder_Set(t *testing.T) {
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type fields struct {
		viper        *viper.Viper
		clientGetter ClientsGetter
		clientSetter ClientsSetter
		mapMerge     func(dst interface{}, src interface{}, opts ...func(*mergo.Config)) error
	}
	type args struct {
		chain        string
		network      string
		pubkeyFinder string
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantErr      bool
		wantSettings map[string]interface{}
	}{
		{
			"success",
			fields{
				func() *viper.Viper {
					v := viper.New()
					return v
				}(),
				nil,
				func() ClientsSetter {
					g := configtest.NewMockClientsSetter(mockCtrl)
					g.EXPECT().SetClient("etherscan-no-auth", "mainnet").Return(nil)
					return g
				}(),
				nil,
			},
			args{
				"ethereum",
				"mainnet",
				"etherscan-no-auth",
			},
			false,
			map[string]interface{}{"pubkey-finder": "etherscan-no-auth"},
		},
		{
			"err-set-client",
			fields{
				func() *viper.Viper {
					v := viper.New()
					return v
				}(),
				nil,
				func() ClientsSetter {
					g := configtest.NewMockClientsSetter(mockCtrl)
					g.EXPECT().SetClient("etherscan-no-auth", "mainnet").Return(errors.Errorf("failed to set"))
					return g
				}(),
				nil,
			},
			args{
				"ethereum",
				"mainnet",
				"etherscan-no-auth",
			},
			true,
			map[string]interface{}{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := PubKeyFinder{
				viper:        tt.fields.viper,
				clientGetter: tt.fields.clientGetter,
				clientSetter: tt.fields.clientSetter,
				mapMerge:     tt.fields.mapMerge,
			}
			if err := p.Set(tt.args.chain, tt.args.network, tt.args.pubkeyFinder); (err != nil) != tt.wantErr {
				t.Errorf("PubKeyFinder.Set() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !assert.EqualValues(tt.wantSettings, tt.fields.viper.GetStringMap("chains.ethereum.networks.mainnet")) {
				t.Errorf("wantSettings = expected")
			}
		})
	}
}
