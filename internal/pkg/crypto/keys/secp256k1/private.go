package secp256k1

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/crypto"
)

// PrivateKey based on the secp256k1 curve
type PrivateKey struct {
	ecdsa ecdsa.PrivateKey
}

// Bytes returns the byte representation of the private key
func (pk PrivateKey) Bytes() []byte {
	return crypto.FromECDSA(&pk.ecdsa)
}

// PrivateKeyFromECDSA get a private key from an ecdsa.PrivateKey
func PrivateKeyFromECDSA(pk ecdsa.PrivateKey) PrivateKey {
	return PrivateKey{ecdsa: pk}
}
