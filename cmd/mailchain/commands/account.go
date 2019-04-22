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
	"encoding/hex"
	"fmt"

	"github.com/mailchain/mailchain/cmd/mailchain/config"
	"github.com/mailchain/mailchain/internal/mailchain/commands/prerun"
	"github.com/mailchain/mailchain/internal/pkg/crypto/keys/multikey"
	"github.com/mailchain/mailchain/internal/pkg/keystore/kdf/multi"
	"github.com/mailchain/mailchain/internal/pkg/keystore/kdf/scrypt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/ttacon/chalk"
)

// account represents the say command
func accountCmd() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "account",
		Short: "Manage Accounts",
		Long: `
Manage accounts, list all existing accounts, import a private key into a new
account, create a new account or update an existing account.
Make sure you remember the password you gave when creating a new account (with
either new or import). Without it you are not able to unlock your account.
Keys are stored under <DATADIR>/keystore.
It is safe to transfer the entire directory or the individual keys therein
between ethereum nodes by simply copying.

Make sure you backup your keys regularly.`,
		PersistentPreRunE: prerun.InitConfig,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Usage()
		},
	}
	cmd.PersistentFlags().StringP("key-type", "", "", "Select the chain [secp256k1]")

	addCmd, err := accountAddCmd()
	if err != nil {
		return nil, err
	}
	cmd.AddCommand(addCmd)
	cmd.AddCommand(accountListCmd())

	return cmd, nil
}

func accountAddCmd() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add private key",
		RunE: func(cmd *cobra.Command, args []string) error {
			keytype, err := getKeyType(cmd)
			if err != nil {
				return errors.WithMessage(err, "could not determine `key-type`")
			}
			privateKey, err := cmd.Flags().GetString("private-key")
			if err != nil {
				return errors.WithMessage(err, "could not get `private-key`")
			}
			if privateKey == "" {
				return errors.Errorf("private key must be specified")
			}
			pk, err := multikey.PrivateKeyFromHex(privateKey, keytype)
			if err != nil {
				return errors.WithMessage(err, "`private-key` could not be decoded")
			}

			passphrase, err := config.Passphrase(cmd)
			if err != nil {
				return errors.WithMessage(err, "could not get `passphrase`")
			}
			randomSalt, err := scrypt.RandomSalt()
			if err != nil {
				return errors.WithMessage(err, "could not create `random salt`")
			}
			// TODO: currently this only does scrypt need flag + config etc
			ks, err := config.GetKeystore()
			if err != nil {
				return errors.WithMessage(err, "could not create `keystore`")
			}
			address, err := ks.Store(pk, keytype,
				multi.OptionsBuilders{
					Scrypt: []scrypt.DeriveOptionsBuilder{
						scrypt.DefaultDeriveOptions(),
						scrypt.WithPassphrase(passphrase),
						randomSalt,
					},
				})
			if err != nil {
				return errors.WithMessage(err, "key could not be stored")
			}

			cmd.Printf(chalk.Green.Color("Private key added")+" for address=%s\n", hex.EncodeToString(address))
			return nil
		},
	}

	cmd.Flags().StringP("chain", "C", "", "Select the chain [ethereum]")
	cmd.Flags().StringP("private-key", "K", "", "Specify the private key to store")
	cmd.Flags().String("passphrase", "", "Passphrase to encrypt/decrypt key with")

	if err := cmd.MarkFlagRequired("private-key"); err != nil {
		return nil, err
	}

	return cmd, nil
}

func accountListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List accounts",
		RunE: func(cmd *cobra.Command, args []string) error {
			ks, err := config.GetKeystore()
			if err != nil {
				return errors.WithMessage(err, "could not create `keystore`")
			}
			addresses, err := ks.GetAddresses()
			if err != nil {
				return errors.WithMessage(err, "could not get addresses")
			}
			for _, x := range addresses {
				fmt.Println(hex.EncodeToString(x))
			}
			return nil
		},
	}
}

func getKeyType(cmd *cobra.Command) (string, error) {
	keyType, err := cmd.Flags().GetString("key-type")
	if err != nil {
		return "", errors.WithMessage(err, "failed to get `key-type`")
	}
	if keyType != "" {
		return keyType, nil
	}
	chain, err := cmd.Flags().GetString("chain")
	if err != nil {
		return "", errors.WithMessage(err, "failed to get `chain`")
	}
	if chain == "" {
		return "", errors.Errorf("either `key-type` or `chain` must be specified")
	}
	return multikey.GetKeyTypeFromChain(chain)
}
