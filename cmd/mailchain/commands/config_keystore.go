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

// import (
// 	"github.com/mailchain/mailchain"
// 	"github.com/mailchain/mailchain/cmd/mailchain/internal/config/defaults"
// 	"github.com/mailchain/mailchain/cmd/mailchain/internal/setup"
// 	"github.com/mailchain/mailchain/internal/keystore"
// 	"github.com/pkg/errors"
// 	"github.com/spf13/cobra"
// )

// func configKeystore(keystoreSelector setup.KeystoreSelector) *cobra.Command {
// 	validArgs := keystore.KeystoreNames()
// 	cmd := &cobra.Command{
// 		Use:   "keystore KEYSTORE",
// 		Short: "configure storage for private keys",
// 		Long: `Mailchain stores the private keys in an encrypted format. Private Keys are used when 
//   sending messages: creating a transactions
//   reading messages: decrypting a message the sender encrypted with the corresponding public key`,
// 		Example:   formatExampleText("mailchain config storage keys mailchain", validArgs),
// 		Args:      exactAndOnlyValid(1),
// 		ValidArgs: validArgs,
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			store := args[0]
// 			path, _ := cmd.Flags().GetString("keystore-path")
// 			if path == "" {
// 				path = mailchain.RequiresValue
// 			}

// 			keystoreType, err := keystoreSelector.Select(store, path)
// 			if err != nil {
// 				return errors.WithStack(err)
// 			}
// 			cmd.Printf("Key store %q configured\n", keystoreType)
// 			return nil
// 		},
// 	}
// 	cmd.Flags().String("keystore-path", defaults.KeystorePath, "Path where keys will be stored.")
// 	return cmd
// }
