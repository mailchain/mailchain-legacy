package substrate

import (
	"github.com/centrifuge/go-substrate-rpc-client/types"
	"github.com/mailchain/mailchain/crypto"
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
		ext := &opts.Tx
		o := opts.SignatureOptions
		mb, err := types.EncodeToBytes(ext.Method)
		if err != nil {
			return nil, err
		}
		era := o.Era
		if !o.Era.IsMortalEra {
			era = types.ExtrinsicEra{IsImmortalEra: true}
		}
		payload := types.ExtrinsicPayloadV3{
			Method:      mb,
			Era:         era,
			Nonce:       o.Nonce,
			Tip:         o.Tip,
			SpecVersion: o.SpecVersion,
			GenesisHash: o.GenesisHash,
			BlockHash:   o.BlockHash,
		}
		data, err := types.EncodeToBytes(payload)
		if err != nil {
			return nil, err
		}
		signedData, err := e.privateKey.Sign(data)
		if err != nil {
			return nil, err
		}
		sig := types.NewSignature(signedData)
		signerPubKey := types.NewAddressFromAccountID(e.privateKey.PublicKey().Bytes())
		extSig := types.ExtrinsicSignatureV4{
			Signer:    signerPubKey,
			Signature: types.MultiSignature{IsSr25519: true, AsSr25519: sig},
			Era:       era,
			Nonce:     o.Nonce,
			Tip:       o.Tip,
		}
		ext.Signature = extSig
		ext.Version |= types.ExtrinsicBitSigned
		return ext, nil
	default:
		return nil, errors.New("invalid options for substrate signing")
	}
}
