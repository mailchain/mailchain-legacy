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
	"fmt"

	"github.com/mailchain/mailchain/cmd/mailchain/config"
	"github.com/spf13/cobra"
	"github.com/ttacon/chalk"
)

func rootCmd() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "mailchain",
		Short: "MailChain node.",
		Long: `Decentralized Mailchain client, run it locally.
Complete documentation is available at github.com/mailchain/mailchain`,
	}
	var cfgFile string
	var logLevel string
	cmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.mailchain/.mailchain.yaml)")
	cmd.PersistentFlags().StringVar(&logLevel, "log-level", "warn", "log level [Panic,Fatal,Error,Warn,Info,Debug]")

	// TODO: this should not be persistent flags
	cmd.PersistentFlags().Bool("empty-passphrase", false, "no passphrase and no prompt")

	err := config.Init(cfgFile, logLevel)
	if err != nil {
		fmt.Println(err)

		fmt.Printf("Run %s to configure create or specify with %s\n",
			chalk.Bold.TextStyle("`mailchain init`"),
			chalk.Bold.TextStyle("`--config`"))
	}

	account, err := accountCmd()
	if err != nil {
		return nil, err
	}
	serve, err := serveCmd()
	if err != nil {
		return nil, err
	}
	cmd.AddCommand(account)
	cmd.AddCommand(cfgCmd())
	cmd.AddCommand(serve)
	return cmd, nil
}
