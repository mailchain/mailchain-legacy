package commands

import (
	"encoding/hex"
	"fmt"

	"github.com/mailchain/mailchain/cmd/mailchain/internal/prompts"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings"
	"github.com/mailchain/mailchain/crypto/multikey"
	"github.com/mailchain/mailchain/internal/keystore"
	"github.com/mailchain/mailchain/internal/keystore/kdf/multi"
	"github.com/mailchain/mailchain/internal/keystore/kdf/scrypt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/ttacon/chalk"
)

// account represents the say command
func accountCmd(config *settings.Root) (*cobra.Command, error) {
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
	}
	ks, err := config.Keystore.Produce()
	if err != nil {
		return nil, errors.WithMessage(err, "could not create `keystore`")
	}

	cmd.AddCommand(accountAddCmd(ks, prompts.Secret, prompts.Secret))
	cmd.AddCommand(accountListCmd(ks))

	return cmd, nil
}

func accountAddCmd(ks keystore.Store, passphrasePrompt, privateKeyPrompt prompts.SecretFunc) *cobra.Command { //nolint: funlen
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add private key",
		RunE: func(cmd *cobra.Command, args []string) error {
			keyType, _ := cmd.Flags().GetString("key-type")
			if keyType == "" {
				return errors.New("`key-type` must be specified")
			}

			cmdPK, _ := cmd.Flags().GetString("private-key")
			privateKey, err := privateKeyPrompt(cmdPK, "", "Private Key", false, false)
			if err != nil {
				return errors.WithMessage(err, "could not get private key")
			}

			privKeyBytes, err := hex.DecodeString(privateKey)
			if err != nil {
				return errors.WithMessage(err, "`private-key` could not be decoded")
			}
			privKey, err := multikey.PrivateKeyFromBytes(keyType, privKeyBytes)
			if err != nil {
				return errors.WithMessage(err, "`private-key` could not be created from bytes")
			}
			cmdPassphrase, _ := cmd.Flags().GetString("passphrase")
			passphrase, err := passphrasePrompt(cmdPassphrase,
				fmt.Sprint(chalk.Yellow, "Note: To derive a storage key passphrase is required. The passphrase must be secure and not guessable."),
				"Passphrase",
				false,
				true,
			)
			if err != nil {
				return errors.WithMessage(err, "could not get `passphrase`")
			}
			randomSalt, err := scrypt.RandomSalt()
			if err != nil {
				return errors.WithMessage(err, "could not create `random salt`")
			}
			pubKey, err := ks.Store(privKey,
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

			cmd.Printf(chalk.Green.Color("Private key added\n"))
			cmd.Printf("Public key=%s\n", hex.EncodeToString(pubKey.Bytes()))
			return nil
		},
	}

	cmd.Flags().StringP("protocol", "P", "", "Select the protocol [ethereum]")
	cmd.Flags().StringP("key-type", "", "", "Select the key type [secp256k1, ed25519]")
	cmd.Flags().StringP("private-key", "K", "", "Private key to store")
	cmd.Flags().String("passphrase", "", "Passphrase to encrypt/decrypt key with")

	return cmd
}

func accountListCmd(ks keystore.Store) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List accounts",
		RunE: func(cmd *cobra.Command, args []string) error {
			protocol, _ := cmd.Flags().GetString("protocol")
			if protocol == "" {
				return errors.New("`--protocol` must be specified to return address list")
			}
			network, _ := cmd.Flags().GetString("network")
			if network == "" {
				return errors.New("`--network` must be specified to return address list")
			}
			addresses, err := ks.GetAddresses(protocol, network)
			if err != nil {
				return errors.WithMessage(err, "could not get addresses")
			}
			for _, x := range addresses {
				cmd.Println(hex.EncodeToString(x))
			}
			return nil
		},
	}
	cmd.Flags().StringP("protocol", "", "", "Protocol to search for")
	cmd.Flags().StringP("network", "", "", "Network to search for")

	return cmd
}
