package algorand

import (
	"crypto/ed25519"

	algocrypto "github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/types"
	"github.com/mailchain/mailchain/crypto"
	mced25519 "github.com/mailchain/mailchain/crypto/ed25519"
	"github.com/mailchain/mailchain/internal/mailbox/signer"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type SignerOptions struct {
	Transaction types.Transaction
}

func NewSigner(privateKey crypto.PrivateKey) (*Signer, error) {
	pk, err := validatePrivateKeyType(privateKey)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &Signer{privateKey: pk}, nil
}

type Signer struct {
	privateKey ed25519.PrivateKey
}

func (e Signer) Sign(opts signer.Options) (signedTransaction interface{}, err error) {
	switch opts := opts.(type) {
	case SignerOptions:
		logger := log.With().Str("component", "signer").Logger()

		txID, signedTxn, err := algocrypto.SignTransaction(e.privateKey, opts.Transaction)
		if err != nil {
			return nil, errors.Wrap(err, "failed to sign transaction")
		}

		logger = logger.With().Str("transaction-id", txID).Logger()
		logger.Info().Msg("transaction signed")

		return signedTxn, nil
	default:
		return nil, errors.New("invalid algorand signer options")
	}
}

func validatePrivateKeyType(pk crypto.PrivateKey) (ed25519.PrivateKey, error) {
	switch pk := pk.(type) {
	case *mced25519.PrivateKey:
		return pk.Bytes(), nil
	case mced25519.PrivateKey:
		return pk.Bytes(), nil
	default:
		return nil, errors.New("invalid key type")
	}
}
