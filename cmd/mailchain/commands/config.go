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
	"github.com/mailchain/mailchain/cmd/mailchain/internal/setup"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func cfgCmd(viper *viper.Viper) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "config",
		Short:   "Config mailchain",
		Aliases: []string{"cfg"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Usage()
		},
	}
	cmd.AddCommand(cfgChainCmd(viper))
	cmd.AddCommand(cfgStorage(viper))

	return cmd
}

func cfgChainCmd(viper *viper.Viper) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "chain",
		Short: "Select a chain to configure",
		// Long:  ``,
		Example:           "",
		PersistentPreRunE: prerunInitConfig(viper),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Usage()
		},
	}
	cmd.AddCommand(cfgChainEthereum())
	return cmd
}
func cfgStorage(viper *viper.Viper) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "storage",
		Short: "Select a storage backend to configure",
		Long:  "Mailchain has multiple storage backends, this command you can configure each of them.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Usage()
		},
	}
	cmd.AddCommand(cfgKeystore(viper))
	cmd.AddCommand(cfgStorageSent(viper, setup.DefaultSentStorage()))

	return cmd
}
