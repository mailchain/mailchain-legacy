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
	"net/http"

	mchttp "github.com/mailchain/mailchain/cmd/mailchain/http"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/fetching"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper" //nolint: depguard
	"github.com/ttacon/chalk"
)

func serveCmd() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Serve the mailchain application",
		RunE: func(cmd *cobra.Command, args []string) error {
			config := settings.FromStore(viper.GetViper())
			config.Logger.Init()

			mailboxStore, err := config.MailboxState.Produce()
			if err != nil {
				return errors.WithMessage(err, "Could not config mailbox store")
			}

			if err := fetching.Do(config, mailboxStore); err != nil {
				return err
			}

			router, err := mchttp.CreateRouter(config, mailboxStore, cmd)
			if err != nil {
				return err
			}

			cmd.Println(chalk.Green.Color("Mailchain started."))
			cmd.Println(("Check messages at https://inbox.mailchain.xyz."))
			cmd.Printf("View developer documention at http://127.0.0.1:%d/api/docs.\n", config.Server.Port.Get())

			listenAddress := fmt.Sprintf(":%d", config.Server.Port.Get())
			log.Info().Str("address", listenAddress).Msg("serving http")

			return http.ListenAndServe(
				listenAddress,
				mchttp.CreateNegroni(config.Server, router),
			)
		},
	}

	if err := mchttp.SetupFlags(cmd); err != nil {
		return nil, err
	}

	return cmd, nil
}
