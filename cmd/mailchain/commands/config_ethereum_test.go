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

package commands

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mailchain/mailchain/cmd/mailchain/commands/commandstest"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/setup"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/setup/setuptest"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func Test_configChainEthereumNetwork(t *testing.T) {
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type args struct {
		chain                string
		network              string
		receiverSelector     setup.ChainNetworkExistingSelector
		senderSelector       setup.ChainNetworkExistingSelector
		pubKeyFinderSelector setup.ChainNetworkExistingSelector
	}
	tests := []struct {
		name        string
		args        args
		cmdArgs     []string
		cmdFlags    map[string]string
		wantOutput  string
		wantExecErr bool
	}{
		{
			"success",
			args{
				"ethereum",
				"mainnet",
				func() setup.ChainNetworkExistingSelector {
					g := setuptest.NewMockChainNetworkExistingSelector(mockCtrl)
					g.EXPECT().Select("ethereum", "mainnet", "supplied-receiver").Return("selected-receiver", nil)
					return g
				}(),
				func() setup.ChainNetworkExistingSelector {
					g := setuptest.NewMockChainNetworkExistingSelector(mockCtrl)
					g.EXPECT().Select("ethereum", "mainnet", "supplied-sender").Return("selected-sender", nil)
					return g
				}(),
				func() setup.ChainNetworkExistingSelector {
					g := setuptest.NewMockChainNetworkExistingSelector(mockCtrl)
					g.EXPECT().Select("ethereum", "mainnet", "supplied-pkf").Return("selected-pkf", nil)
					return g
				}(),
			},
			nil,
			map[string]string{
				"receiver":          "supplied-receiver",
				"sender":            "supplied-sender",
				"public-key-finder": "supplied-pkf",
			},
			"ethereum mainnet configured using:\n- supplied-sender: messages sent from key owner\n- supplied-receiver: messages sent to key owner\n- supplied-pkf: looking up addresses\n",
			false,
		},
		{
			"err-receiver",
			args{
				"ethereum",
				"mainnet",
				func() setup.ChainNetworkExistingSelector {
					g := setuptest.NewMockChainNetworkExistingSelector(mockCtrl)
					g.EXPECT().Select("ethereum", "mainnet", "supplied-receiver").Return("", errors.Errorf("failed"))
					return g
				}(),
				func() setup.ChainNetworkExistingSelector {
					g := setuptest.NewMockChainNetworkExistingSelector(mockCtrl)
					return g
				}(),
				func() setup.ChainNetworkExistingSelector {
					g := setuptest.NewMockChainNetworkExistingSelector(mockCtrl)
					return g
				}(),
			},
			nil,
			map[string]string{
				"receiver":          "supplied-receiver",
				"sender":            "supplied-sender",
				"public-key-finder": "supplied-pkf",
			},
			"Error: failed",
			true,
		},
		{
			"err-sender",
			args{
				"ethereum",
				"mainnet",
				func() setup.ChainNetworkExistingSelector {
					g := setuptest.NewMockChainNetworkExistingSelector(mockCtrl)
					g.EXPECT().Select("ethereum", "mainnet", "supplied-receiver").Return("selected-receiver", nil)
					return g
				}(),
				func() setup.ChainNetworkExistingSelector {
					g := setuptest.NewMockChainNetworkExistingSelector(mockCtrl)
					g.EXPECT().Select("ethereum", "mainnet", "supplied-sender").Return("", errors.Errorf("failed"))
					return g
				}(),
				func() setup.ChainNetworkExistingSelector {
					g := setuptest.NewMockChainNetworkExistingSelector(mockCtrl)
					return g
				}(),
			},
			nil,
			map[string]string{
				"receiver":          "supplied-receiver",
				"sender":            "supplied-sender",
				"public-key-finder": "supplied-pkf",
			},
			"Error: failed",
			true,
		},
		{
			"err-public-key-finder",
			args{
				"ethereum",
				"mainnet",
				func() setup.ChainNetworkExistingSelector {
					g := setuptest.NewMockChainNetworkExistingSelector(mockCtrl)
					g.EXPECT().Select("ethereum", "mainnet", "supplied-receiver").Return("selected-receiver", nil)
					return g
				}(),
				func() setup.ChainNetworkExistingSelector {
					g := setuptest.NewMockChainNetworkExistingSelector(mockCtrl)
					g.EXPECT().Select("ethereum", "mainnet", "supplied-sender").Return("selected-sender", nil)
					return g
				}(),
				func() setup.ChainNetworkExistingSelector {
					g := setuptest.NewMockChainNetworkExistingSelector(mockCtrl)
					g.EXPECT().Select("ethereum", "mainnet", "supplied-pkf").Return("", errors.Errorf("failed"))
					return g
				}(),
			},
			nil,
			map[string]string{
				"receiver":          "supplied-receiver",
				"sender":            "supplied-sender",
				"public-key-finder": "supplied-pkf",
			},
			"Error: failed",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := configChainEthereumNetwork(tt.args.chain, tt.args.network, tt.args.receiverSelector, tt.args.senderSelector, tt.args.pubKeyFinderSelector)
			if !assert.NotNil(got) {
				t.Error("configChainEthereumNetwork() is nil")
			}
			_, out, err := commandstest.ExecuteCommandC(got, tt.cmdArgs, tt.cmdFlags)
			if (err != nil) != tt.wantExecErr {
				t.Errorf("configChainEthereumNetwork().execute() error = %v, wantExecErr %v", err, tt.wantExecErr)
				return
			}
			if !commandstest.AssertCommandOutput(t, got, err, out, tt.wantOutput) {
				t.Errorf("configChainEthereumNetwork().Execute().out != %v", tt.wantOutput)
			}
		})
	}
}

func Test_configChainEthereum(t *testing.T) {
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type args struct {
		receiverSelector     setup.ChainNetworkExistingSelector
		senderSelector       setup.ChainNetworkExistingSelector
		pubKeyFinderSelector setup.ChainNetworkExistingSelector
	}
	tests := []struct {
		name        string
		args        args
		cmdArgs     []string
		cmdFlags    map[string]string
		wantExecErr bool
	}{
		{
			"success",
			args{
				func() setup.ChainNetworkExistingSelector {
					g := setuptest.NewMockChainNetworkExistingSelector(mockCtrl)
					return g
				}(),
				func() setup.ChainNetworkExistingSelector {
					g := setuptest.NewMockChainNetworkExistingSelector(mockCtrl)
					return g
				}(),
				func() setup.ChainNetworkExistingSelector {
					g := setuptest.NewMockChainNetworkExistingSelector(mockCtrl)
					return g
				}(),
			},
			nil,
			map[string]string{
				"receiver":          "supplied-receiver",
				"sender":            "supplied-sender",
				"public-key-finder": "supplied-pkf",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := configChainEthereum(tt.args.receiverSelector, tt.args.senderSelector, tt.args.pubKeyFinderSelector)
			if !assert.NotNil(got) {
				t.Error("configChainEthereum() is nil")
			}
			_, out, err := commandstest.ExecuteCommandC(got, tt.cmdArgs, tt.cmdFlags)
			if (err != nil) != tt.wantExecErr {
				t.Errorf("configChainEthereum().execute() error = %v, wantExecErr %v", err, tt.wantExecErr)
				return
			}
			if !assert.Equal(got.UsageString(), out) {
				t.Errorf("configChainEthereum().Execute().out != %v", got.Usage())
			}
		})
	}
}
