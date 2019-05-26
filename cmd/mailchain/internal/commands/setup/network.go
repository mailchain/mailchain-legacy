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
	"github.com/mailchain/mailchain/cmd/mailchain/config/names"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/commands/prompts"
	"github.com/mailchain/mailchain/internal/pkg/encoding"
	"github.com/spf13/cobra"
)

func Network(cmd *cobra.Command, args []string, chain, network string) (string, error) {
	network, err := selectNetwork(cmd, args, network, chainNetworks(chain))
	if err != nil {
		return "", err
	}
	if _, err := Receiver(cmd, chain, network, names.EtherscanNoAuth); err != nil {
		return "", err
	}
	if _, err := Sender(cmd, chain, network, names.Empty); err != nil {
		return "", err
	}
	if _, err := PublicKeyFinder(cmd, chain, network, names.EtherscanNoAuth); err != nil {
		return "", err
	}
	return network, nil
}
func chainNetworks(chain string) []string {
	switch chain {
	case encoding.Ethereum:
		return encoding.EthereumNetworks()
	default:
		return nil
	}
}
func selectNetwork(cmd *cobra.Command, args []string, network string, networks []string) (string, error) {
	if network != names.Empty {
		return network, nil
	}
	flg, _ := cmd.Flags().GetString("network")
	if flg != "" {
		return flg, nil
	}
	if len(args) == 1 {
		return args[0], nil
	}
	return prompts.SelectItem("Network", networks)
}
