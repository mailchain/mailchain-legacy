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

	"github.com/mailchain/mailchain/cmd/mailchain/internal/http"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings"
	"github.com/spf13/cobra"
	"github.com/spf13/viper" //nolint: depguard
	"github.com/ttacon/chalk"
)

func serveCmd() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Serve the mailchain application",
		// PersistentPreRunE: preRun,
		RunE: func(cmd *cobra.Command, args []string) error {
			config := settings.FromStore(viper.GetViper())
			router, err := http.CreateRouter(config, cmd)
			if err != nil {
				return err
			}
			fmt.Println(chalk.Bold.TextStyle(fmt.Sprintf(
				"Find out more by visiting the docs http://127.0.0.1:%d/api/docs",
				config.Server.Port.Get())))

			http.CreateNegroni(config.Server, router).Run(fmt.Sprintf(":%d", config.Server.Port.Get()))
			return nil
		},
	}

	if err := http.SetupFlags(cmd); err != nil {
		return nil, err
	}
	return cmd, nil
}
