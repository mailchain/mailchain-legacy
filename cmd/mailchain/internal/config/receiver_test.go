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
	"net/http/httptest"
	"reflect"
	"sort"
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

func sortReceiverMapKeys(m map[string]mailbox.Receiver) []string {
	ss := reflectAsString(reflect.ValueOf(m).MapKeys())
	sort.Strings(ss)
	return ss
}

func TestReceiver_getReceiver(t *testing.T) {
	assert := assert.New(t)
	server := httptest.NewServer(nil)
	defer server.Close()
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
		want    mailbox.Receiver
		wantErr bool
	}{
		{
			"etherscan",
			fields{
				func() *viper.Viper {
					v := viper.New()
					v.Set("chains.ethereum.networks.mainnet.receiver", "etherscan")
					v.Set("clients.etherscan.mainnet.address", server.URL)
					return v
				}(),
				func() ClientsGetter {
					g := configtest.NewMockClientsGetter(mockCtrl)
					g.EXPECT().GetEtherscanClient().Return(etherscan.NewAPIClient(server.URL))
					return g
				}(),
				nil,
				nil,
			},
			args{
				"ethereum",
				"mainnet",
			},
			func() mailbox.Receiver {
				e, _ := etherscan.NewAPIClient(server.URL)
				return e
			}(),
			false,
		},
		{
			"etherscan-no-auth",
			fields{
				func() *viper.Viper {
					v := viper.New()
					v.Set("chains.ethereum.networks.mainnet.receiver", "etherscan-no-auth")
					v.Set("clients.etherscan.mainnet.address", server.URL)
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
			func() mailbox.Receiver {
				e, _ := etherscan.NewAPIClient("")
				return e
			}(),
			false,
		},
		{
			"error",
			fields{
				func() *viper.Viper {
					v := viper.New()
					v.Set("chains.ethereum.networks.mainnet.receiver", "unknown")
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
				"mainnet",
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Receiver{
				viper:        tt.fields.viper,
				clientGetter: tt.fields.clientGetter,
				clientSetter: tt.fields.clientSetter,
				mapMerge:     tt.fields.mapMerge,
			}
			got, err := s.getReceiver(tt.args.chain, tt.args.network)
			if (err != nil) != tt.wantErr {
				t.Errorf("Receiver.getReceiver() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.IsType(tt.want, got) {
				t.Errorf("Receiver.getReceiver() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReceiver_getChainReceivers(t *testing.T) {
	assert := assert.New(t)
	server := httptest.NewServer(nil)
	defer server.Close()
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
		want    map[string]mailbox.Receiver
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
			args{
				"ethereum",
			},
			nil,
			false,
		},
		{
			"single",
			fields{
				func() *viper.Viper {
					v := viper.New()
					v.Set("chains.ethereum.networks.mainnet.receiver", "etherscan-no-auth")
					v.Set("clients.etherscan.mainnet.address", server.URL)
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
			map[string]mailbox.Receiver{
				"ethereum.mainnet": func() mailbox.Receiver {
					e, _ := etherscan.NewAPIClient("")
					return e
				}(),
			},
			false,
		},
		{
			"multi",
			fields{
				func() *viper.Viper {
					v := viper.New()
					v.Set("chains.ethereum.networks.mainnet.receiver", "etherscan-no-auth")
					v.Set("chains.ethereum.networks.ropsten.receiver", "etherscan-no-auth")
					v.Set("clients.etherscan.mainnet.address", server.URL)
					v.Set("clients.etherscan.ropsten.address", server.URL)
					return v
				}(),
				func() ClientsGetter {
					g := configtest.NewMockClientsGetter(mockCtrl)
					g.EXPECT().GetEtherscanNoAuthClient().Return(etherscan.NewAPIClient(""))
					g.EXPECT().GetEtherscanNoAuthClient().Return(etherscan.NewAPIClient(""))
					return g
				}(),
				nil,
				nil,
			},
			args{
				"ethereum",
			},
			map[string]mailbox.Receiver{
				"ethereum.mainnet": func() mailbox.Receiver {
					e, _ := etherscan.NewAPIClient("")
					return e
				}(),
				"ethereum.ropsten": func() mailbox.Receiver {
					e, _ := etherscan.NewAPIClient("")
					return e
				}(),
			},
			false,
		},
		{
			"err",
			fields{
				func() *viper.Viper {
					v := viper.New()
					// v.Set("chains.ethereum.networks.mainnet.receiver", "etherscan")
					v.Set("chains.ethereum.networks.ropsten.receiver", "unknown")
					v.Set("clients.etherscan.mainnet.address", server.URL)
					v.Set("clients.etherscan.ropsten.address", server.URL)
					return v
				}(),
				func() ClientsGetter {
					g := configtest.NewMockClientsGetter(mockCtrl)
					// g.EXPECT().GetEtherscanNoAuthClient().Return(etherscan.NewAPIClient(""))
					// g.EXPECT().GetEtherscanNoAuthClient().Return(etherscan.NewAPIClient("")"failed"))
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
			s := Receiver{
				viper:        tt.fields.viper,
				clientGetter: tt.fields.clientGetter,
				clientSetter: tt.fields.clientSetter,
				mapMerge:     tt.fields.mapMerge,
			}
			got, err := s.getChainReceivers(tt.args.chain)
			if (err != nil) != tt.wantErr {
				t.Errorf("Receiver.getChainReceivers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !assert.EqualValues(sortReceiverMapKeys(tt.want), sortReceiverMapKeys(got)) {
				t.Errorf("Receiver.getChainReceivers() missing keys = %s, want %s", sortReceiverMapKeys(tt.want), sortReceiverMapKeys(got))
			}

			for x := range tt.want {
				if !assert.IsType(tt.want[x], got[x]) {
					t.Errorf("Receiver.getChainReceivers().[%s] = %v, want %v", x, got, tt.want)
				}
			}
		})
	}
}

func TestReceiver_GetReceiver(t *testing.T) {
	assert := assert.New(t)
	server := httptest.NewServer(nil)
	defer server.Close()
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
		want    map[string]mailbox.Receiver
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
				mergo.Merge,
			},
			nil,
			false,
		},
		{
			"single",
			fields{
				func() *viper.Viper {
					v := viper.New()
					v.Set("chains.ethereum.networks.mainnet.receiver", "etherscan-no-auth")
					v.Set("clients.etherscan.mainnet.address", server.URL)
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
			map[string]mailbox.Receiver{
				"ethereum.mainnet": func() mailbox.Receiver {
					e, _ := etherscan.NewAPIClient("")
					return e
				}(),
			},
			false,
		},
		{
			"multi",
			fields{
				func() *viper.Viper {
					v := viper.New()
					v.Set("chains.ethereum.networks.mainnet.receiver", "etherscan-no-auth")
					v.Set("chains.ethereum.networks.ropsten.receiver", "etherscan-no-auth")
					v.Set("clients.etherscan.mainnet.address", server.URL)
					v.Set("clients.etherscan.ropsten.address", server.URL)
					return v
				}(),
				func() ClientsGetter {
					g := configtest.NewMockClientsGetter(mockCtrl)
					g.EXPECT().GetEtherscanNoAuthClient().Return(etherscan.NewAPIClient(""))
					g.EXPECT().GetEtherscanNoAuthClient().Return(etherscan.NewAPIClient(""))
					return g
				}(),
				nil,
				mergo.Merge,
			},
			map[string]mailbox.Receiver{
				"ethereum.mainnet": func() mailbox.Receiver {
					e, _ := etherscan.NewAPIClient("")
					return e
				}(),
				"ethereum.ropsten": func() mailbox.Receiver {
					e, _ := etherscan.NewAPIClient("")
					return e
				}(),
			},
			false,
		},
		{
			"err-mergo",
			fields{
				func() *viper.Viper {
					v := viper.New()
					v.Set("chains.ethereum.networks.mainnet.receiver", "etherscan-no-auth")
					v.Set("chains.ethereum.networks.ropsten.receiver", "etherscan-no-auth")
					v.Set("clients.etherscan.mainnet.address", server.URL)
					v.Set("clients.etherscan.ropsten.address", server.URL)
					return v
				}(),
				func() ClientsGetter {
					g := configtest.NewMockClientsGetter(mockCtrl)
					g.EXPECT().GetEtherscanNoAuthClient().Return(etherscan.NewAPIClient(""))
					g.EXPECT().GetEtherscanNoAuthClient().Return(etherscan.NewAPIClient(""))
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
		{
			"err",
			fields{
				func() *viper.Viper {
					v := viper.New()
					v.Set("chains.ethereum.networks.mainnet.receiver", "etherscan-no-auth")
					v.Set("clients.etherscan.mainnet.address", server.URL)
					v.Set("clients.etherscan.ropsten.address", server.URL)
					return v
				}(),
				func() ClientsGetter {
					g := configtest.NewMockClientsGetter(mockCtrl)
					g.EXPECT().GetEtherscanNoAuthClient().Return(nil, errors.Errorf("failed"))
					return g
				}(),
				nil,
				mergo.Merge,
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Receiver{
				viper:        tt.fields.viper,
				clientGetter: tt.fields.clientGetter,
				clientSetter: tt.fields.clientSetter,
				mapMerge:     tt.fields.mapMerge,
			}
			got, err := s.GetReceivers()
			if (err != nil) != tt.wantErr {
				t.Errorf("Receiver.GetReceivers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !assert.EqualValues(sortReceiverMapKeys(tt.want), sortReceiverMapKeys(got)) {
				t.Errorf("Receiver.GetReceivers() missing keys = %s, want %s", sortReceiverMapKeys(tt.want), sortReceiverMapKeys(got))
			}

			for x := range tt.want {
				if !assert.IsType(tt.want[x], got[x]) {
					t.Errorf("Receiver.GetReceivers().[%s] = %v, want %v", x, got, tt.want)
				}
			}
		})
	}
}

func TestReceiver_Set(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type fields struct {
		viper        *viper.Viper
		clientGetter ClientsGetter
		clientSetter ClientsSetter
		mapMerge     func(dst interface{}, src interface{}, opts ...func(*mergo.Config)) error
	}
	type args struct {
		chain    string
		network  string
		receiver string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"success",
			fields{
				viper.New(),
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
		},
		{
			"error",
			fields{
				viper.New(),
				nil,
				func() ClientsSetter {
					g := configtest.NewMockClientsSetter(mockCtrl)
					g.EXPECT().SetClient("etherscan-no-auth", "mainnet").Return(errors.Errorf("failed"))
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
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Receiver{
				viper:        tt.fields.viper,
				clientGetter: tt.fields.clientGetter,
				clientSetter: tt.fields.clientSetter,
				mapMerge:     tt.fields.mapMerge,
			}
			if err := s.Set(tt.args.chain, tt.args.network, tt.args.receiver); (err != nil) != tt.wantErr {
				t.Errorf("Receiver.Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
