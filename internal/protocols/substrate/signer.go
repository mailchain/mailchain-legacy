package substrate

import (
	"github.com/centrifuge/go-substrate-rpc-client/types"
	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/sr25519"
	"github.com/mailchain/mailchain/internal/mailbox/signer"
	"github.com/pkg/errors"
)

type SignerOptions struct {
	Extrinsic        types.Extrinsic
	SignatureOptions types.SignatureOptions
}

func NewSigner(privateKey crypto.PrivateKey) (*Signer, error) {
	return &Signer{privateKey: privateKey}, nil
}

type Signer struct {
	privateKey crypto.PrivateKey
}

func (e Signer) createSignature(signedData []byte) (*types.MultiSignature, error) {
	sig := types.NewSignature(signedData)
	switch e.privateKey.(type) {
	case *sr25519.PrivateKey:
		return &types.MultiSignature{IsSr25519: true, AsSr25519: sig}, nil
	default:
		return nil, errors.New("unsupported private key type")
	}
}

func (e Signer) prepareData(ext *types.Extrinsic, opts *types.SignatureOptions) ([]byte, error) {
	mb, err := types.EncodeToBytes(ext.Method)
	if err != nil {
		return nil, err
	}
	era := opts.Era
	if !opts.Era.IsMortalEra {
		era = types.ExtrinsicEra{IsImmortalEra: true}
	}
	payload := types.ExtrinsicPayloadV3{
		Method:      mb,
		Era:         era,
		Nonce:       opts.Nonce,
		Tip:         opts.Tip,
		SpecVersion: opts.SpecVersion,
		GenesisHash: opts.GenesisHash,
		BlockHash:   opts.BlockHash,
	}
	return types.EncodeToBytes(payload)
}

func (e Signer) Sign(opts signer.Options) (signedTransaction interface{}, err error) {
	switch opts := opts.(type) {
	case *SignerOptions:
		data, err := e.prepareData(&opts.Extrinsic, &opts.SignatureOptions)
		if err != nil {
			return nil, err
		}
		signedData, err := e.privateKey.Sign(data)
		if err != nil {
			return nil, err
		}
		signature, err := e.createSignature(signedData)
		if err != nil {
			return nil, err
		}
		extSig := types.ExtrinsicSignatureV4{
			Signer:    types.NewAddressFromAccountID(e.privateKey.PublicKey().Bytes()),
			Signature: *signature,
			Era:       opts.SignatureOptions.Era,
			Nonce:     opts.SignatureOptions.Nonce,
			Tip:       opts.SignatureOptions.Tip,
		}
		ext := &opts.Extrinsic
		ext.Signature = extSig
		ext.Version |= types.ExtrinsicBitSigned
		return ext, nil
	default:
		return nil, errors.New("invalid options for substrate signing")
	}
}
