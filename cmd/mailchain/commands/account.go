package commands

import (
	"fmt"
	"strings"

	"github.com/mailchain/mailchain/cmd/mailchain/internal/prompts"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings"
	"github.com/mailchain/mailchain/crypto/multikey"
	"github.com/mailchain/mailchain/encoding"
	"github.com/mailchain/mailchain/internal/addressing"
	"github.com/mailchain/mailchain/internal/keystore"
	"github.com/mailchain/mailchain/internal/keystore/kdf/multi"
	"github.com/mailchain/mailchain/internal/keystore/kdf/scrypt"
	"github.com/mailchain/mailchain/internal/protocols"
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

func getPrivateKeyBytes(privateKeyEncoding, privateKeyInput string) ([]byte, error) {
	switch privateKeyEncoding {
	case encoding.KindHex:
		return encoding.DecodeHex(privateKeyInput)
	case encoding.KindMnemonicAlgorand:
		return encoding.DecodeMnemonicAlgorand(privateKeyInput)
	default:
		return nil, errors.New("private key encoding type not supported.")
	}
}

func accountAddCmd(ks keystore.Store, passphrasePrompt, privateKeyPrompt prompts.SecretFunc) *cobra.Command { //nolint: funlen
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add private key",
		RunE: func(cmd *cobra.Command, args []string) error {
			keyType, _ := cmd.Flags().GetString("key-type")
			cmdPK, _ := cmd.Flags().GetString("private-key")
			cmdPassphrase, _ := cmd.Flags().GetString("passphrase")
			privateKeyEncoding, _ := cmd.Flags().GetString("private-key-encoding")

			privateKey, err := privateKeyPrompt(cmdPK, "", "Private Key", false, false)
			if err != nil {
				return errors.WithMessage(err, "could not get private key")
			}

			privKeyBytes, err := getPrivateKeyBytes(privateKeyEncoding, privateKey)
			if err != nil {
				return errors.WithMessage(err, "`private-key` could not be decoded")
			}

			privKey, err := multikey.PrivateKeyFromBytes(keyType, privKeyBytes)
			if err != nil {
				return errors.WithMessage(err, "`private-key` could not be created from bytes")
			}

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
			cmd.Printf("Public key=%s\n", encoding.EncodeHex(pubKey.Bytes()))
			return nil
		},
	}

	cmd.Flags().StringP("protocol", "P", "", fmt.Sprintf("Select the protocol [%s]", strings.Join([]string{protocols.Algorand, protocols.Ethereum, protocols.Substrate}, ", ")))
	cmd.Flags().StringP("key-type", "", "", "Select the key type [secp256k1, ed25519, sr25519]")
	_ = cmd.MarkFlagRequired("key-type")
	cmd.Flags().StringP("private-key", "K", "", "Private key to store encoded")
	cmd.Flags().StringP("private-key-encoding", "E", encoding.KindHex, fmt.Sprintf("Encoding used for supplied private key [%s]", strings.Join([]string{encoding.KindHex, encoding.KindMnemonicAlgorand}, ", ")))
	cmd.Flags().String("passphrase", "", "Passphrase to encrypt/decrypt key with")

	return cmd
}

func accountListCmd(ks keystore.Store) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List accounts",
		RunE: func(cmd *cobra.Command, args []string) error {
			protocol, _ := cmd.Flags().GetString("protocol")
			network, _ := cmd.Flags().GetString("network")

			addresses, err := ks.GetAddresses(protocol, network)
			if err != nil {
				return errors.WithMessage(err, "could not get addresses")
			}
			for _, x := range addresses {
				encoded, encoding, err := addressing.EncodeByProtocol(x, protocol)
				if err != nil {
					return errors.WithMessage(err, "could not encode address")
				}

				cmd.Printf("Encoding: %s, address: %s\n", encoding, encoded)
			}
			return nil
		},
	}

	cmd.Flags().StringP("protocol", "", "", "Protocol to search for")
	_ = cmd.MarkFlagRequired("protocol")
	cmd.Flags().StringP("network", "", "", "Network to search for")
	_ = cmd.MarkFlagRequired("network")

	return cmd
}
