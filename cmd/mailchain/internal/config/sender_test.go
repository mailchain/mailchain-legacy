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
	"fmt"
	"net/http/httptest"
	"reflect"
	"sort"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/imdario/mergo"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/config/configtest"
	"github.com/mailchain/mailchain/sender/ethrpc2"
	"github.com/mailchain/mailchain/sender"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func reflectAsString(v []reflect.Value) []string {
	r := []string{}
	for x := range v {
		r = append(r, fmt.Sprintf("%s", v[x]))
	}
	return r
}

func sortSenderMapKeys(m map[string]sender.Message) []string {
	ss := reflectAsString(reflect.ValueOf(m).MapKeys())
	sort.Strings(ss)
	return ss
}

type ByReflect []reflect.Value

func (a ByReflect) Len() int           { return len(a) }
func (a ByReflect) Less(i, j int) bool { return fmt.Sprintf("%s", a[i]) < fmt.Sprintf("%s", a[j]) }
func (a ByReflect) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func TestSender_getSender(t *testing.T) {
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
		want    sender.Message
		wantErr bool
	}{
		{
			"etherscan",
			fields{
				func() *viper.Viper {
					v := viper.New()
					v.Set("chains.ethereum.networks.mainnet.sender", "ethereum-rpc2")
					v.Set("clients.ethereum-rpc2.mainnet.address", server.URL)
					return v
				}(),
				func() ClientsGetter {
					g := configtest.NewMockClientsGetter(mockCtrl)
					g.EXPECT().GetEtherRPC2Client("mainnet").Return(ethrpc2.New(server.URL))
					return g
				}(),
				nil,
				nil,
			},
			args{
				"ethereum",
				"mainnet",
			},
			&ethrpc2.EthRPC2{},
			false,
		},
		{
			"error",
			fields{
				func() *viper.Viper {
					v := viper.New()
					v.Set("chains.ethereum.networks.mainnet.sender", "unknown")
					return v
				}(),
				func() ClientsGetter {
					g := configtest.NewMockClientsGetter(mockCtrl)
					// g.EXPECT().GetEtherRPC2Client("mainnet").Return(ethrpc2.New(server.URL))
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
			s := Sender{
				viper:        tt.fields.viper,
				clientGetter: tt.fields.clientGetter,
				clientSetter: tt.fields.clientSetter,
				mapMerge:     tt.fields.mapMerge,
			}
			got, err := s.getSender(tt.args.chain, tt.args.network)
			if (err != nil) != tt.wantErr {
				t.Errorf("Sender.getSender() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.IsType(tt.want, got) {
				t.Errorf("Sender.getSender() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSender_getChainSenders(t *testing.T) {
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
		want    map[string]sender.Message
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
					v.Set("chains.ethereum.networks.mainnet.sender", "ethereum-rpc2")
					v.Set("clients.ethereum-rpc2.mainnet.address", server.URL)
					return v
				}(),
				func() ClientsGetter {
					g := configtest.NewMockClientsGetter(mockCtrl)
					g.EXPECT().GetEtherRPC2Client("mainnet").Return(ethrpc2.New(server.URL))
					return g
				}(),
				nil,
				nil,
			},
			args{
				"ethereum",
			},
			map[string]sender.Message{
				"ethereum.mainnet": &ethrpc2.EthRPC2{},
			},
			false,
		},
		{
			"multi",
			fields{
				func() *viper.Viper {
					v := viper.New()
					v.Set("chains.ethereum.networks.mainnet.sender", "ethereum-rpc2")
					v.Set("chains.ethereum.networks.ropsten.sender", "ethereum-rpc2")
					v.Set("clients.ethereum-rpc2.mainnet.address", server.URL)
					v.Set("clients.ethereum-rpc2.ropsten.address", server.URL)
					return v
				}(),
				func() ClientsGetter {
					g := configtest.NewMockClientsGetter(mockCtrl)
					g.EXPECT().GetEtherRPC2Client("mainnet").Return(ethrpc2.New(server.URL))
					g.EXPECT().GetEtherRPC2Client("ropsten").Return(ethrpc2.New(server.URL))
					return g
				}(),
				nil,
				nil,
			},
			args{
				"ethereum",
			},
			map[string]sender.Message{
				"ethereum.mainnet": &ethrpc2.EthRPC2{},
				"ethereum.ropsten": &ethrpc2.EthRPC2{},
			},
			false,
		},
		{
			"err",
			fields{
				func() *viper.Viper {
					v := viper.New()
					// v.Set("chains.ethereum.networks.mainnet.sender", "ethereum-rpc2")
					v.Set("chains.ethereum.networks.ropsten.sender", "unknown")
					v.Set("clients.ethereum-rpc2.mainnet.address", server.URL)
					v.Set("clients.ethereum-rpc2.ropsten.address", server.URL)
					return v
				}(),
				func() ClientsGetter {
					g := configtest.NewMockClientsGetter(mockCtrl)
					// g.EXPECT().GetEtherRPC2Client("mainnet").Return(ethrpc2.New(server.URL))
					// g.EXPECT().GetEtherRPC2Client("ropsten").Return(nil, errors.Errorf("failed"))
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
			s := Sender{
				viper:        tt.fields.viper,
				clientGetter: tt.fields.clientGetter,
				clientSetter: tt.fields.clientSetter,
				mapMerge:     tt.fields.mapMerge,
			}
			got, err := s.getChainSenders(tt.args.chain)
			if (err != nil) != tt.wantErr {
				t.Errorf("Sender.getChainSenders() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !assert.EqualValues(sortSenderMapKeys(tt.want), sortSenderMapKeys(got)) {
				t.Errorf("Sender.getChainSenders() missing keys = %s, want %s", sortSenderMapKeys(tt.want), sortSenderMapKeys(got))
			}

			for x := range tt.want {
				if !assert.IsType(tt.want[x], got[x]) {
					t.Errorf("Sender.getChainSenders().[%s] = %v, want %v", x, got, tt.want)
				}
			}
		})
	}
}

func TestSender_GetSenders(t *testing.T) {
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
		want    map[string]sender.Message
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
					v.Set("chains.ethereum.networks.mainnet.sender", "ethereum-rpc2")
					v.Set("clients.ethereum-rpc2.mainnet.address", server.URL)
					return v
				}(),
				func() ClientsGetter {
					g := configtest.NewMockClientsGetter(mockCtrl)
					g.EXPECT().GetEtherRPC2Client("mainnet").Return(ethrpc2.New(server.URL))
					return g
				}(),
				nil,
				mergo.Merge,
			},
			map[string]sender.Message{
				"ethereum.mainnet": &ethrpc2.EthRPC2{},
			},
			false,
		},
		{
			"multi",
			fields{
				func() *viper.Viper {
					v := viper.New()
					v.Set("chains.ethereum.networks.mainnet.sender", "ethereum-rpc2")
					v.Set("chains.ethereum.networks.ropsten.sender", "ethereum-rpc2")
					v.Set("clients.ethereum-rpc2.mainnet.address", server.URL)
					v.Set("clients.ethereum-rpc2.ropsten.address", server.URL)
					return v
				}(),
				func() ClientsGetter {
					g := configtest.NewMockClientsGetter(mockCtrl)
					g.EXPECT().GetEtherRPC2Client("mainnet").Return(ethrpc2.New(server.URL))
					g.EXPECT().GetEtherRPC2Client("ropsten").Return(ethrpc2.New(server.URL))
					return g
				}(),
				nil,
				mergo.Merge,
			},
			map[string]sender.Message{
				"ethereum.mainnet": &ethrpc2.EthRPC2{},
				"ethereum.ropsten": &ethrpc2.EthRPC2{},
			},
			false,
		},
		{
			"err-mergo",
			fields{
				func() *viper.Viper {
					v := viper.New()
					v.Set("chains.ethereum.networks.mainnet.sender", "ethereum-rpc2")
					v.Set("chains.ethereum.networks.ropsten.sender", "ethereum-rpc2")
					v.Set("clients.ethereum-rpc2.mainnet.address", server.URL)
					v.Set("clients.ethereum-rpc2.ropsten.address", server.URL)
					return v
				}(),
				func() ClientsGetter {
					g := configtest.NewMockClientsGetter(mockCtrl)
					g.EXPECT().GetEtherRPC2Client("mainnet").Return(ethrpc2.New(server.URL))
					g.EXPECT().GetEtherRPC2Client("ropsten").Return(ethrpc2.New(server.URL))
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
					v.Set("chains.ethereum.networks.mainnet.sender", "ethereum-rpc2")
					v.Set("clients.ethereum-rpc2.mainnet.address", server.URL)
					v.Set("clients.ethereum-rpc2.ropsten.address", server.URL)
					return v
				}(),
				func() ClientsGetter {
					g := configtest.NewMockClientsGetter(mockCtrl)
					g.EXPECT().GetEtherRPC2Client("mainnet").Return(nil, errors.Errorf("failed"))
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
			s := Sender{
				viper:        tt.fields.viper,
				clientGetter: tt.fields.clientGetter,
				clientSetter: tt.fields.clientSetter,
				mapMerge:     tt.fields.mapMerge,
			}
			got, err := s.GetSenders()
			if (err != nil) != tt.wantErr {
				t.Errorf("Sender.GetSenders() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !assert.EqualValues(sortSenderMapKeys(tt.want), sortSenderMapKeys(got)) {
				t.Errorf("Sender.GetSenders() missing keys = %s, want %s", sortSenderMapKeys(tt.want), sortSenderMapKeys(got))
			}

			for x := range tt.want {
				if !assert.IsType(tt.want[x], got[x]) {
					t.Errorf("Sender.GetSenders().[%s] = %v, want %v", x, got, tt.want)
				}
			}
		})
	}
}

func TestSender_Set(t *testing.T) {
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
		sender  string
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
			s := Sender{
				viper:        tt.fields.viper,
				clientGetter: tt.fields.clientGetter,
				clientSetter: tt.fields.clientSetter,
				mapMerge:     tt.fields.mapMerge,
			}
			if err := s.Set(tt.args.chain, tt.args.network, tt.args.sender); (err != nil) != tt.wantErr {
				t.Errorf("Sender.Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
