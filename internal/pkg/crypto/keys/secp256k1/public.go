package secp256k1

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
)

// PublicKey based on the secp256k1 curve
type PublicKey struct {
	ecdsa ecdsa.PublicKey
}

// Bytes returns the byte representation of the public key
func (pk PublicKey) Bytes() []byte {
	return crypto.CompressPubkey(&pk.ecdsa)
}

// Address returns the byte representation of the address
func (pk PublicKey) Address() []byte {
	return crypto.PubkeyToAddress(pk.ecdsa).Bytes()
}

// PublicKeyFromBytes create a public key from []byte
func PublicKeyFromBytes(pk []byte) (*PublicKey, error) {
	rpk, err := crypto.UnmarshalPubkey(pk)
	if err != nil {
		return nil, errors.WithMessage(err, "could not convert pk")
	}
	return &PublicKey{ecdsa: *rpk}, nil
}
