package ed25519

import (
	"golang.org/x/crypto/ed25519"
)

// PublicKey based on the secp256k1 curve
type PublicKey struct {
	key ed25519.PublicKey
}

// Bytes returns the byte representation of the public key
func (pk PublicKey) Bytes() []byte {
	return pk.key
}

// Address returns the byte representation of the address
func (pk PublicKey) Address() []byte {
	return nil
}
