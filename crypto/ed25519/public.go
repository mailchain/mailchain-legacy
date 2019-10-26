package ed25519

import (
	"github.com/pkg/errors"
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

// PublicKeyFromBytes create a public key from []byte
func PublicKeyFromBytes(keyBytes []byte) (*PublicKey, error) {
	if len(keyBytes) != 32 {
		return nil, errors.Errorf("public key must be 32 bytes")
	}
	return &PublicKey{key: keyBytes}, nil
}
