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
	"github.com/mailchain/mailchain/cmd/mailchain/config/defaults"
	"github.com/mailchain/mailchain/cmd/mailchain/config/names"
	"github.com/mailchain/mailchain/internal/pkg/cmd/prerun"
	"github.com/mailchain/mailchain/internal/pkg/cmd/setup"
	"github.com/spf13/cobra"
)

func cfgKeystore() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "keys",
		Short: "setup keystore",
		// Long:  ``,
		PreRunE:  prerun.InitConfig,
		PostRunE: config.WriteConfig,
		RunE: func(cmd *cobra.Command, args []string) error {
			keystoreType, err := setup.Keystore(cmd, names.Empty)
			if err != nil {
				return err
			}
			cmd.Printf("Key store %q configured\n", keystoreType)
			return cmd.Usage()
		},
	}
	cmd.AddCommand(cfgKeystoreNaclFilestore())
	return cmd
}

func cfgKeystoreNaclFilestore() *cobra.Command {
	cmd := &cobra.Command{
		Use:      "nacl-filestore",
		Short:    "setup nacl filestore",
		PreRunE:  prerun.InitConfig,
		PostRunE: config.WriteConfig,
		RunE: func(cmd *cobra.Command, args []string) error {
			keystoreType, err := setup.Keystore(cmd, names.KeystoreNACLFilestore)
			if err != nil {
				return err
			}
			cmd.Printf("Key store %q configured\n", keystoreType)
			return nil
		},
	}
	cmd.Flags().String("keystore-path", defaults.KeystorePath, "Path where keys will be stored.")
	return cmd
}
