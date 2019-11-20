package sr25519

import (
	"errors"

	sr25519 "github.com/ChainSafe/go-schnorrkel"
	r255 "github.com/gtank/ristretto255"
)

const (
	publicKeySize = 32
	seedSize      = 32
)

// PublicKey is a member
type PublicKey struct {
	key *sr25519.PublicKey
}

// Verify uses the sr25519 signature algorithm to verify that the message was signed by
// this public key; it returns true if this key created the signature for the message,
// false otherwise
func (k *PublicKey) Verify(msg, sig []byte) bool {
	if k.key == nil {
		return false
	}

	b := [64]byte{}
	copy(b[:], sig)

	s := &sr25519.Signature{}
	err := s.Decode(b)
	if err != nil {
		return false
	}

	t := sr25519.NewSigningContext(SigningContext, msg)
	return k.key.Verify(s, t)
}

// Encode returns the 32-byte encoding of the public key
func (k *PublicKey) Encode() []byte {
	if k.key == nil {
		return nil
	}

	enc := k.key.Encode()
	return enc[:]
}

// Decode decodes the input bytes into a public key and sets the receiver the decoded key
// Input must be 32 bytes, or else this function will error
func (k *PublicKey) Decode(in []byte) error {
	if len(in) != publicKeySize {
		return errors.New("input to sr25519 public key decode is not 32 bytes")
	}
	b := [32]byte{}
	copy(b[:], in)
	k.key = &sr25519.PublicKey{}
	return k.key.Decode(b)
}

// New Public key sr25519
// PublicKey -> Key -> PublicKey -> key
func NewPublicKey(b [32]byte) *PublicKey {
	e := r255.NewElement()
	e.Decode(b[:])

	srPubKey := &sr25519.PublicKey{e}
	// pubKey := PublicKey{srPubKey}
	// pbk := &pubKey

	return &PublicKey{key: srPubKey}
}

// Convert this public key to a byte array.
func PublicKeyFromBytes(keyBytes []byte) (*PublicKey, error) {
	if len(keyBytes) != publicKeySize {
		return nil, errors.New("public key must be 32 bytes")
	}

	kb := [32]byte{}
	copy(kb[:], keyBytes)
	pub := NewPublicKey(kb)
	return pub, nil
}
