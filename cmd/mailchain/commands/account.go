package commands

import (
	"bytes"
	"encoding/json"
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
	"github.com/mailchain/mailchain/internal/pubkey"
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
	produceKeystore := func() (keystore.Store, error) {
		return config.Keystore.Produce()
	}

	cmd.AddCommand(accountAddCmd(produceKeystore, prompts.Secret, prompts.Secret))
	cmd.AddCommand(accountListCmd(produceKeystore))

	return cmd, nil
}

func accountAddCmd(produceKeyStore func() (keystore.Store, error), passphrasePrompt, privateKeyPrompt prompts.SecretFunc) *cobra.Command { //nolint: funlen
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add private key",
		RunE: func(cmd *cobra.Command, args []string) error {
			protocol, _ := cmd.Flags().GetString("protocol")
			if len(protocols.NetworkNames(protocol)) == 0 {
				return errors.New("protocol has no available networks")
			}

			network, _ := cmd.Flags().GetString("network")
			if network == "" {
				return errors.New("network must be specified")
			}

			if !contains(protocols.NetworkNames(protocol), network) {
				return errors.New("network not found for protocol")
			}

			ks, err := produceKeyStore()
			if err != nil {
				return errors.WithMessage(err, "could not create `keystore`")
			}

			keyType, _ := cmd.Flags().GetString("key-type")
			cmdPK, _ := cmd.Flags().GetString("private-key")
			cmdPassphrase, _ := cmd.Flags().GetString("passphrase")
			privateKeyEncoding, _ := cmd.Flags().GetString("private-key-encoding")

			privateKey, err := privateKeyPrompt(cmdPK, "", "Private Key", false, false)
			if err != nil {
				return errors.WithMessage(err, "could not get private key")
			}

			privKeyBytes, err := encoding.Decode(privateKeyEncoding, privateKey)
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

			pk, err := ks.Store(protocol, network, privKey, multi.OptionsBuilders{Scrypt: []scrypt.DeriveOptionsBuilder{scrypt.DefaultDeriveOptions(), scrypt.WithPassphrase(passphrase), randomSalt}})
			if err != nil {
				return errors.WithMessage(err, "key could not be stored")
			}

			encodedPubKey, publicKeyEncoding, err := pubkey.EncodeByProtocol(pk.Bytes(), protocol)
			if err != nil {
				return errors.WithMessage(err, "public key could not be encoded")
			}

			type response struct {
				Message           string `json:"message"`
				PublicKey         string `json:"public-key"`
				PublicKeyEncoding string `json:"public-key-encoding"`
				Address           string `json:"address"`
				AddressEncoding   string `json:"address-encoding"`
				Protocol          string `json:"protocol"`
				Network           string `json:"network"`
			}

			addressBytes, err := addressing.FromPublicKey(pk, protocol, network)
			if err != nil {
				return errors.WithMessage(err, "could not get address fom public key")
			}

			encodedAddress, addressEncoding, err := addressing.EncodeByProtocol(addressBytes, protocol)
			if err != nil {
				return errors.WithMessage(err, "could not encode address")
			}

			jsonResponse, _ := json.Marshal(response{
				Message:           "private key added",
				Protocol:          protocol,
				Network:           network,
				PublicKey:         encodedPubKey,
				PublicKeyEncoding: publicKeyEncoding,
				Address:           encodedAddress,
				AddressEncoding:   addressEncoding,
			})
			var prettyJSON bytes.Buffer
			if err := json.Indent(&prettyJSON, jsonResponse, "", "  "); err != nil {
				return errors.WithMessage(err, "public key could not be encoded")
			}

			cmd.Print(prettyJSON.String())

			return nil
		},
	}

	cmd.Flags().StringP("protocol", "P", "", fmt.Sprintf("Protocol to add private key to. Available: [%s].", strings.Join(protocols.All(), ", ")))
	_ = cmd.MarkFlagRequired("protool")
	cmd.Flags().StringP("network", "N", "", "Network to add the private key to.")
	_ = cmd.MarkFlagRequired("network")
	cmd.Flags().StringP("key-type", "", "", "Select the key type [secp256k1, ed25519, sr25519]")
	_ = cmd.MarkFlagRequired("key-type")
	cmd.Flags().StringP("private-key", "K", "", "Private key to store encoded")
	cmd.Flags().StringP("private-key-encoding", "E", encoding.KindHex, fmt.Sprintf("Encoding used for supplied private key [%s]", strings.Join([]string{encoding.KindHex, encoding.KindMnemonicAlgorand}, ", ")))
	cmd.Flags().String("passphrase", "", "Passphrase to encrypt/decrypt key with")

	return cmd
}

func accountListCmd(produceKeystore func() (keystore.Store, error)) *cobra.Command {
	type address struct {
		Address         string `json:"address"`
		AddressEncoding string `json:"address-encoding"`
		Protocol        string `json:"protocol"`
		Network         string `json:"network"`
	}

	type response struct {
		Addresses []address `json:"addresses"`
		Protocol  string    `json:"protocol"`
		Network   string    `json:"network"`
	}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List accounts",
		RunE: func(cmd *cobra.Command, args []string) error {
			ks, err := produceKeystore()
			if err != nil {
				return errors.WithMessage(err, "could not create `keystore`")
			}

			protocol, _ := cmd.Flags().GetString("protocol")
			network, _ := cmd.Flags().GetString("network")

			allAddresses, err := ks.GetAddresses(protocol, network)
			if err != nil {
				return errors.WithMessage(err, "could not get addresses")
			}

			addresses := []address{}
			for _, a := range keystore.FlattenAddressesMap(allAddresses) {
				encodedAddress, encodingUsed, err := addressing.EncodeByProtocol(a.Address, a.Protocol)
				if err != nil {
					return errors.WithStack(err)
				}

				addresses = append(addresses, address{
					Protocol:        a.Protocol,
					Network:         a.Network,
					Address:         encodedAddress,
					AddressEncoding: encodingUsed,
				})
			}

			jsonResponse, _ := json.Marshal(response{
				Addresses: addresses,
				Protocol:  protocol,
				Network:   network,
			})
			var prettyJSON bytes.Buffer
			if err := json.Indent(&prettyJSON, jsonResponse, "", "  "); err != nil {
				return errors.WithMessage(err, "list addresses could not be encoded")
			}

			cmd.Print(prettyJSON.String())

			return nil
		},
	}

	cmd.Flags().StringP("protocol", "", "", "Protocol to search for")
	cmd.Flags().StringP("network", "", "", "Network to search for")

	return cmd
}

func contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}

	return false
}
