package secp256k1

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/crypto"
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
