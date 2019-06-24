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
	"github.com/mailchain/mailchain"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/config/defaults"
	"github.com/mailchain/mailchain/internal/chains"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// QUESTION: can this method be removed?
func (n Network) Select(cmd *cobra.Command, args []string, chain, network string) (string, error) {
	networkNames := chains.NetworkNames(chain)
	if len(networkNames) == 0 {
		return "", errors.Errorf("no network found for chain")
	}
	selectedNetwork, err := n.selectNetwork(cmd, args, network, networkNames)
	if err != nil {
		return "", err
	}
	if _, err := n.receiverSelector.Select(chain, selectedNetwork, mailchain.ClientEtherscanNoAuth); err != nil {
		return "", err
	}
	if _, err := n.senderSelector.Select(chain, selectedNetwork, defaults.Sender); err != nil {
		return "", err
	}
	if _, err := n.pubKeyFinderSelector.Select(chain, selectedNetwork, mailchain.ClientEtherscanNoAuth); err != nil {
		return "", err
	}
	return selectedNetwork, nil
}

func (n Network) selectNetwork(cmd *cobra.Command, args []string, existingNetwork string, networks []string) (string, error) {
	if existingNetwork != mailchain.RequiresValue {
		return existingNetwork, nil
	}
	if networkFromCommand := n.networkFromCLI(cmd, args); networkFromCommand != "" {
		return networkFromCommand, nil
	}

	return n.selectItem("Network", networks)
}

func (n Network) networkFromCLI(cmd *cobra.Command, args []string) string {
	if cmd == nil {
		return ""
	}
	flg, _ := cmd.Flags().GetString("network")
	if flg != "" {
		return flg
	}
	if len(args) == 1 {
		return args[0]
	}
	return ""
}
