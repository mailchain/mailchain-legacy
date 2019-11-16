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

// Kind is the type of private key.
func (pk PrivateKey) Kind() string {
	return crypto.ED25519
}

// Sign signs the message with the private key and returns the signature.
func (pk PrivateKey) Sign(message []byte) (signature []byte, err error) {
	if len(pk.key) != ed25519.PrivateKeySize {
		return nil, errors.New("invalid key length")
	}

	return ed25519.Sign(pk.key, message), nil
}

// PublicKey return the public key that is derived from the private key
func (pk PrivateKey) PublicKey() crypto.PublicKey {
	publicKey := make([]byte, ed25519.PublicKeySize)
	copy(publicKey, pk.key[32:])
	return PublicKey{key: publicKey}
}

// PrivateKeyFromBytes get a private key from seed []byte
func PrivateKeyFromBytes(privKey []byte) (*PrivateKey, error) {
	switch len(privKey) {
	case ed25519.SeedSize:
		return &PrivateKey{key: ed25519.NewKeyFromSeed(privKey)}, nil
	case ed25519.PrivateKeySize:
		return &PrivateKey{key: privKey}, nil
	default:
		return nil, errors.Errorf("ed25519: bad key length")
	}
}
