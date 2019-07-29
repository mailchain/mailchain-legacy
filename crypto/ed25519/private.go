package ed25519

import (
	"github.com/mailchain/mailchain/crypto"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ed25519"
)

// PrivateKey based on the secp256k1 curve
type PrivateKey struct {
	key ed25519.PrivateKey
}

// Bytes returns the byte representation of the private key
func (pk PrivateKey) Bytes() []byte {
	return pk.key
}

// PublicKey return the public key that is derived from the private key
func (pk PrivateKey) PublicKey() crypto.PublicKey {
	publicKey := make([]byte, ed25519.PublicKeySize)
	copy(publicKey, pk.key[32:])
	return PublicKey{key: publicKey}
}

// PrivateKeyFromBytes get a private key from []byte
func PrivateKeyFromBytes(pk []byte) (*PrivateKey, error) {
	if len(pk) != ed25519.PrivateKeySize {
		return nil, errors.Errorf("len must be %v", ed25519.PrivateKeySize)
	}
	return &PrivateKey{key: pk}, nil
}
