package substrate

import (
	"github.com/centrifuge/go-substrate-rpc-client/signature"
	"github.com/centrifuge/go-substrate-rpc-client/types"
	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/encoding"
	"github.com/mailchain/mailchain/internal/mailbox/signer"
	"github.com/pkg/errors"
)

type SignerOptions struct {
	Tx               types.Extrinsic
	SignatureOptions types.SignatureOptions
}

func NewSigner(privateKey crypto.PrivateKey) (*Signer, error) {
	return &Signer{privateKey: privateKey}, nil
}

type Signer struct {
	privateKey crypto.PrivateKey
}

func (e Signer) Sign(opts signer.Options) (signedTransaction interface{}, err error) {
	switch opts := opts.(type) {
	case SignerOptions:
		pair, err := signature.KeyringPairFromSecret(encoding.EncodeHex(e.privateKey.Bytes()))
		if err != nil {
			return nil, err
		}
		err = opts.Tx.Sign(pair, opts.SignatureOptions)
		if err != nil {
			return nil, err
		}
		return &opts.Tx, nil
	default:
		return nil, errors.New("invalid options for substrate signing")
	}
}
