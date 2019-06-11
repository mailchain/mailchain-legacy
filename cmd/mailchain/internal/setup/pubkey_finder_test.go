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
package setup

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mailchain/mailchain"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/config"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/config/configtest"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/prompts/promptstest"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func TestPubKeyFinder_selectPubKeyFinder(t *testing.T) {
	type fields struct {
		setter             config.PubKeyFinderSetter
		viper              *viper.Viper
		selectItemSkipable func(label string, items []string, skipable bool) (selected string, skipped bool, err error)
	}
	type args struct {
		chain            string
		network          string
		existingPKFinder string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			"success-already-set",
			fields{
				nil,
				nil,
				nil,
			},
			args{
				"ethereum",
				"mainnet",
				"value-already-set",
			},
			"value-already-set",
			false,
		},
		{
			"success-skipped",
			fields{
				nil,
				func() *viper.Viper {
					v := viper.New()
					v.Set("chains.ethereum.networks.mainnet.pubkey-finder", "already-set")
					return v
				}(),
				promptstest.MockSelectItemSkipable(t, []string{mailchain.ClientEtherscanNoAuth, mailchain.ClientEtherscan}, "already-set", true, nil),
			},
			args{
				"ethereum",
				"mainnet",
				mailchain.RequiresValue,
			},
			"",
			false,
		},
		{
			"success-not-skipped",
			fields{
				nil,
				func() *viper.Viper {
					v := viper.New()
					v.Set("chains.ethereum.networks.mainnet.pubkey-finder", "already-set")
					return v
				}(),
				promptstest.MockSelectItemSkipable(t, []string{mailchain.ClientEtherscanNoAuth, mailchain.ClientEtherscan}, "new-value", false, nil),
			},
			args{
				"ethereum",
				"mainnet",
				mailchain.RequiresValue,
			},
			"new-value",
			false,
		},
		{
			"err-skipped",
			fields{
				nil,
				func() *viper.Viper {
					v := viper.New()
					v.Set("chains.ethereum.networks.mainnet.pubkey-finder", "already-set")
					return v
				}(),
				promptstest.MockSelectItemSkipable(t, []string{mailchain.ClientEtherscanNoAuth, mailchain.ClientEtherscan}, "", false, errors.Errorf("failed to select")),
			},
			args{
				"ethereum",
				"mainnet",
				mailchain.RequiresValue,
			},
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := PubKeyFinder{
				setter:             tt.fields.setter,
				viper:              tt.fields.viper,
				selectItemSkipable: tt.fields.selectItemSkipable,
			}
			got, err := f.selectPubKeyFinder(tt.args.chain, tt.args.network, tt.args.existingPKFinder)
			if (err != nil) != tt.wantErr {
				t.Errorf("PubKeyFinder.selectPubKeyFinder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PubKeyFinder.selectPubKeyFinder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPubKeyFinder_Select(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type fields struct {
		setter             config.PubKeyFinderSetter
		viper              *viper.Viper
		selectItemSkipable func(label string, items []string, skipable bool) (selected string, skipped bool, err error)
	}
	type args struct {
		chain            string
		network          string
		existingPKFinder string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			"success-set",
			fields{
				func() config.PubKeyFinderSetter {
					setter := configtest.NewMockPubKeyFinderSetter(mockCtrl)
					setter.EXPECT().Set("ethereum", "mainnet", "new-set-value").Return(nil)
					return setter
				}(),
				nil,
				nil,
			},
			args{
				"ethereum",
				"mainnet",
				"new-set-value",
			},
			"new-set-value",
			false,
		},
		{
			"success-skipped",
			fields{
				func() config.PubKeyFinderSetter {
					setter := configtest.NewMockPubKeyFinderSetter(mockCtrl)
					return setter
				}(),
				func() *viper.Viper {
					v := viper.New()
					v.Set("chains.ethereum.networks.mainnet.pubkey-finder", "already-set")
					return v
				}(),
				promptstest.MockSelectItemSkipable(t, []string{mailchain.ClientEtherscanNoAuth, mailchain.ClientEtherscan}, "already-set", true, nil),
			},
			args{
				"ethereum",
				"mainnet",
				mailchain.RequiresValue,
			},
			"",
			false,
		},
		{
			"err-select-failed",
			fields{
				func() config.PubKeyFinderSetter {
					setter := configtest.NewMockPubKeyFinderSetter(mockCtrl)
					return setter
				}(),
				func() *viper.Viper {
					v := viper.New()
					return v
				}(),
				promptstest.MockSelectItemSkipable(t, []string{mailchain.ClientEtherscanNoAuth, mailchain.ClientEtherscan}, "", true, errors.Errorf("failed to skip")),
			},
			args{
				"ethereum",
				"mainnet",
				mailchain.RequiresValue,
			},
			"",
			true,
		},
		{
			"err-setter-failed",
			fields{
				func() config.PubKeyFinderSetter {
					setter := configtest.NewMockPubKeyFinderSetter(mockCtrl)
					setter.EXPECT().Set("ethereum", "mainnet", "new-setting").Return(errors.Errorf("failed to error"))
					return setter
				}(),
				func() *viper.Viper {
					v := viper.New()
					return v
				}(),
				promptstest.MockSelectItemSkipable(t, []string{mailchain.ClientEtherscanNoAuth, mailchain.ClientEtherscan}, "new-setting", false, nil),
			},
			args{
				"ethereum",
				"mainnet",
				mailchain.RequiresValue,
			},
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := PubKeyFinder{
				setter:             tt.fields.setter,
				viper:              tt.fields.viper,
				selectItemSkipable: tt.fields.selectItemSkipable,
			}
			got, err := f.Select(tt.args.chain, tt.args.network, tt.args.existingPKFinder)
			if (err != nil) != tt.wantErr {
				t.Errorf("PubKeyFinder.Select() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PubKeyFinder.Select() = %v, want %v", got, tt.want)
			}
		})
	}
}
