package ed25519

import (
	"github.com/mailchain/mailchain/crypto"
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

// PublicKeyFromBytes create a public key from []byte
func PublicKeyFromBytes(keyBytes []byte) (crypto.PublicKey, error) {
	if len(keyBytes) != ed25519.PublicKeySize {
		return nil, errors.Errorf("public key must be 32 bytes")
	}
	return &PublicKey{key: keyBytes}, nil
}
