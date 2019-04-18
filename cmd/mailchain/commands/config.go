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
	"github.com/mailchain/mailchain/cmd/mailchain/config"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/prerun"
	"github.com/mailchain/mailchain/internal/pkg/encoding"
	"github.com/spf13/cobra"
)

func cfgCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "config",
		Short:             "config mailchain",
		Aliases:           []string{"cfg"},
		PersistentPreRunE: prerun.InitConfig,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Usage()
		},
	}
	cmd.AddCommand(cfgChainCmd())
	cmd.AddCommand(cfgKeystore())

	return cmd
}

func cfgChainCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "chain",
		Short: "setup chain",
		// Long:  ``,
		Example:           "",
		PersistentPreRunE: prerun.InitConfig,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Usage()
		},
	}
	cmd.AddCommand(cfgChainEthereum())
	return cmd
}

func cfgKeystore() *cobra.Command {
	return &cobra.Command{
		Use:   "keystore",
		Short: "setup keystore",
		// Long:  ``,
		PersistentPreRunE: prerun.InitConfig,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Usage()
		},
	}
}

func cfgChainEthereum() *cobra.Command {
	return &cobra.Command{
		Use:      "ethereum",
		Short:    "setup ethereum",
		PostRunE: config.WriteConfig,
		RunE: func(cmd *cobra.Command, args []string) error {
			network, err := selectNetwork(cmd, args, encoding.EthereumNetworks())
			if err != nil {
				return err
			}
			if err := config.SetReceiver(network); err != nil {
				return err
			}
			if err := config.SetSender(network); err != nil {
				return err
			}
			if err := config.SetPubKeyFinder(network); err != nil {
				return err
			}

			cmd.Printf("Ethereum chain configured\n")
			return nil
		},
	}
}
