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

package setup

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mailchain/mailchain"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/prompts/promptstest"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/setup/setuptest"
	"github.com/mailchain/mailchain/internal/chains/ethereum"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func TestNetwork_networkFromCLI(t *testing.T) {
	type fields struct {
		receiverSelector     ChainNetworkExistingSelector
		senderSelector       ChainNetworkExistingSelector
		pubKeyFinderSelector ChainNetworkExistingSelector
		selectItem           func(label string, items []string) (string, error)
	}
	type args struct {
		cmd  *cobra.Command
		args []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			"nil-cmd",
			fields{
				nil, nil, nil, nil,
			},
			args{
				nil, nil,
			},
			"",
		},
		{
			"cmd-has-flag-value",
			fields{
				nil, nil, nil, nil,
			},
			args{
				func() *cobra.Command {
					cmd := &cobra.Command{}
					cmd.Flags().String("network", "cmd-value", "")
					return cmd
				}(), nil,
			},
			"cmd-value",
		},
		{
			"arg-value",
			fields{
				nil, nil, nil, nil,
			},
			args{
				func() *cobra.Command {
					cmd := &cobra.Command{}
					return cmd
				}(), []string{"arg-value"},
			},
			"arg-value",
		},
		{
			"no-cmd-args",
			fields{
				nil, nil, nil, nil,
			},
			args{
				func() *cobra.Command {
					cmd := &cobra.Command{}
					return cmd
				}(), []string{},
			},
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := Network{
				receiverSelector:     tt.fields.receiverSelector,
				senderSelector:       tt.fields.senderSelector,
				pubKeyFinderSelector: tt.fields.pubKeyFinderSelector,
				selectItem:           tt.fields.selectItem,
			}
			if got := n.networkFromCLI(tt.args.cmd, tt.args.args); got != tt.want {
				t.Errorf("Network.networkFromCLI() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNetwork_selectNetwork(t *testing.T) {
	type fields struct {
		receiverSelector     ChainNetworkExistingSelector
		senderSelector       ChainNetworkExistingSelector
		pubKeyFinderSelector ChainNetworkExistingSelector
		selectItem           func(label string, items []string) (string, error)
	}
	type args struct {
		cmd             *cobra.Command
		args            []string
		existingNetwork string
		networks        []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			"from-prompt",
			fields{
				nil,
				nil,
				nil,
				promptstest.MockSelectItem(t, []string{"mainnet", "testnet"}, "mainnet", nil),
			},
			args{
				nil,
				nil,
				mailchain.RequiresValue,
				[]string{"mainnet", "testnet"},
			},
			"mainnet",
			false,
		},
		{
			"from-existing-value",
			fields{
				nil,
				nil,
				nil,
				promptstest.MockSelectItem(t, []string{"mainnet", "testnet"}, "mainnet", nil),
			},
			args{
				nil,
				nil,
				"existing-value",
				[]string{"mainnet", "testnet"},
			},
			"existing-value",
			false,
		},
		{
			"from-cli-not-empty",
			fields{
				nil,
				nil,
				nil,
				promptstest.MockSelectItem(t, []string{"mainnet", "testnet"}, "mainnet", nil),
			},
			args{
				&cobra.Command{},
				[]string{"arg-value"},
				mailchain.RequiresValue,
				[]string{"mainnet", "testnet"},
			},
			"arg-value",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := Network{
				receiverSelector:     tt.fields.receiverSelector,
				senderSelector:       tt.fields.senderSelector,
				pubKeyFinderSelector: tt.fields.pubKeyFinderSelector,
				selectItem:           tt.fields.selectItem,
			}
			got, err := n.selectNetwork(tt.args.cmd, tt.args.args, tt.args.existingNetwork, tt.args.networks)
			if (err != nil) != tt.wantErr {
				t.Errorf("Network.selectNetwork() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Network.selectNetwork() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNetwork_Select(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type fields struct {
		receiverSelector     ChainNetworkExistingSelector
		senderSelector       ChainNetworkExistingSelector
		pubKeyFinderSelector ChainNetworkExistingSelector
		selectItem           func(label string, items []string) (string, error)
	}
	type args struct {
		cmd     *cobra.Command
		args    []string
		chain   string
		network string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			"success",
			fields{
				func() ChainNetworkExistingSelector {
					selector := setuptest.NewMockChainNetworkExistingSelector(mockCtrl)
					selector.EXPECT().Select("ethereum", "mainnet", "etherscan-no-auth").Return("etherscan-no-auth", nil)
					return selector
				}(),
				func() ChainNetworkExistingSelector {
					selector := setuptest.NewMockChainNetworkExistingSelector(mockCtrl)
					selector.EXPECT().Select("ethereum", "mainnet", mailchain.Relay).Return("etherscan-no-auth", nil)
					return selector
				}(),
				func() ChainNetworkExistingSelector {
					selector := setuptest.NewMockChainNetworkExistingSelector(mockCtrl)
					selector.EXPECT().Select("ethereum", "mainnet", "etherscan-no-auth").Return("etherscan-no-auth", nil)
					return selector
				}(),
				nil,
			},
			args{
				&cobra.Command{},
				[]string{"mainnet"},
				"ethereum",
				"mainnet",
			},
			"mainnet",
			false,
		},
		{
			"err-no-networks",
			fields{
				func() ChainNetworkExistingSelector {
					selector := setuptest.NewMockChainNetworkExistingSelector(mockCtrl)
					return selector
				}(),
				func() ChainNetworkExistingSelector {
					selector := setuptest.NewMockChainNetworkExistingSelector(mockCtrl)
					return selector
				}(),
				func() ChainNetworkExistingSelector {
					selector := setuptest.NewMockChainNetworkExistingSelector(mockCtrl)
					return selector
				}(),
				nil,
			},
			args{
				nil,
				nil,
				"unknown",
				"mainnet",
			},
			"",
			true,
		},
		{
			"err-select-network",
			fields{
				func() ChainNetworkExistingSelector {
					selector := setuptest.NewMockChainNetworkExistingSelector(mockCtrl)
					return selector
				}(),
				func() ChainNetworkExistingSelector {
					selector := setuptest.NewMockChainNetworkExistingSelector(mockCtrl)
					return selector
				}(),
				func() ChainNetworkExistingSelector {
					selector := setuptest.NewMockChainNetworkExistingSelector(mockCtrl)
					return selector
				}(),
				promptstest.MockSelectItem(t, ethereum.Networks(), "", errors.Errorf("failed")),
			},
			args{
				nil,
				nil,
				"ethereum",
				mailchain.RequiresValue,
			},
			"",
			true,
		},

		{
			"err-receiver",
			fields{
				func() ChainNetworkExistingSelector {
					selector := setuptest.NewMockChainNetworkExistingSelector(mockCtrl)
					selector.EXPECT().Select("ethereum", "mainnet", "etherscan-no-auth").Return("", errors.Errorf("failed"))
					return selector
				}(),
				func() ChainNetworkExistingSelector {
					selector := setuptest.NewMockChainNetworkExistingSelector(mockCtrl)
					return selector
				}(),
				func() ChainNetworkExistingSelector {
					selector := setuptest.NewMockChainNetworkExistingSelector(mockCtrl)
					return selector
				}(),
				nil,
			},
			args{
				&cobra.Command{},
				[]string{"mainnet"},
				"ethereum",
				"mainnet",
			},
			"",
			true,
		},
		{
			"err-sender",
			fields{
				func() ChainNetworkExistingSelector {
					selector := setuptest.NewMockChainNetworkExistingSelector(mockCtrl)
					selector.EXPECT().Select("ethereum", "mainnet", "etherscan-no-auth").Return("etherscan-no-auth", nil)
					return selector
				}(),
				func() ChainNetworkExistingSelector {
					selector := setuptest.NewMockChainNetworkExistingSelector(mockCtrl)
					selector.EXPECT().Select("ethereum", "mainnet", mailchain.Relay).Return("", errors.Errorf("failed"))
					return selector
				}(),
				func() ChainNetworkExistingSelector {
					selector := setuptest.NewMockChainNetworkExistingSelector(mockCtrl)
					return selector
				}(),
				nil,
			},
			args{
				&cobra.Command{},
				[]string{"mainnet"},
				"ethereum",
				"mainnet",
			},
			"",
			true,
		},
		{
			"err-pubkey-finder",
			fields{
				func() ChainNetworkExistingSelector {
					selector := setuptest.NewMockChainNetworkExistingSelector(mockCtrl)
					selector.EXPECT().Select("ethereum", "mainnet", "etherscan-no-auth").Return("etherscan-no-auth", nil)
					return selector
				}(),
				func() ChainNetworkExistingSelector {
					selector := setuptest.NewMockChainNetworkExistingSelector(mockCtrl)
					selector.EXPECT().Select("ethereum", "mainnet", mailchain.Relay).Return("etherscan-no-auth", nil)
					return selector
				}(),
				func() ChainNetworkExistingSelector {
					selector := setuptest.NewMockChainNetworkExistingSelector(mockCtrl)
					selector.EXPECT().Select("ethereum", "mainnet", "etherscan-no-auth").Return("", errors.Errorf("failed"))
					return selector
				}(),
				nil,
			},
			args{
				&cobra.Command{},
				[]string{"mainnet"},
				"ethereum",
				"mainnet",
			},
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := Network{
				receiverSelector:     tt.fields.receiverSelector,
				senderSelector:       tt.fields.senderSelector,
				pubKeyFinderSelector: tt.fields.pubKeyFinderSelector,
				selectItem:           tt.fields.selectItem,
			}
			got, err := n.Select(tt.args.cmd, tt.args.args, tt.args.chain, tt.args.network)
			if (err != nil) != tt.wantErr {
				t.Errorf("Network.Select() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Network.Select() = %v, want %v", got, tt.want)
			}
		})
	}
}
