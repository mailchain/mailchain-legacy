package envelope

import (
	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/cipher"
	"github.com/pkg/errors"
)

func NewEnvelope(encrypter cipher.Encrypter, pubkey crypto.PublicKey, o []CreateOptionsBuilder) (Data, error) {
	opts := &CreateOpts{}
	apply(opts, o)
	switch opts.Kind {
	case Kind0x01:
		return NewZeroX01(encrypter, pubkey, opts)
	default:
		return nil, errors.Errorf("unknown kind")
	}
}

func apply(o *CreateOpts, opts []CreateOptionsBuilder) {
	for _, f := range opts {
		f(o)
	}
}
