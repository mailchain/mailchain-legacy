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
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings"
	"github.com/spf13/cobra"
	"github.com/spf13/viper" //nolint: depguard
)

func rootCmd(v *viper.Viper) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "mailchain",
		Short: "Mailchain node.",
		Long: `Decentralized Mailchain client, run it locally.
Complete documentation is available at https://github.com/mailchain/mailchain`,
		PersistentPreRunE: prerunInitConfig(v),
	}
	cmd.PersistentFlags().String("config", "", "config file (default is $HOME/.mailchain/.mailchain.yaml)")
	cmd.PersistentFlags().String("log-level", "warn", "log level [Panic,Fatal,Error,Warn,Info,Debug]")
	cmd.PersistentFlags().Bool("prevent-init-config", false, "stop automatically creating config if no file is found")

	config := settings.FromStore(v)
	account, err := accountCmd(config)
	if err != nil {
		return nil, err
	}

	cmd.AddCommand(account)
	cmd.AddCommand(settingsCmd(config))

	serve, err := serveCmd()
	if err != nil {
		return nil, err
	}

	cmd.AddCommand(serve)

	cmd.AddCommand(versionCmd())
	return cmd, nil
}
