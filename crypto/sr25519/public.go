package sr25519

import (
	"errors"

	"github.com/ChainSafe/go-schnorrkel"
	"github.com/mailchain/mailchain/crypto"
)

const (
	publicKeySize = 32
)

// PublicKey is a interface
type PublicKey struct {
	key *schnorrkel.PublicKey
}

// PublicKey to Bytes
func (pk PublicKey) Bytes() []byte {
	b := pk.key.Encode()

	return b[:]
}

// Kind returns the key type
func (pk PublicKey) Kind() string {
	return crypto.SR25519
}

// Verify uses the sr25519 signature algorithm to verify that the message was signed by
// this public key; it returns true if this key created the signature for the message,
// false otherwise
func (pk PublicKey) Verify(message, sig []byte) bool {
	if pk.key == nil {
		return false
	}

	b := [64]byte{}
	copy(b[:], sig)

	s := &schnorrkel.Signature{}
	err := s.Decode(b)

	if err != nil {
		return false
	}

	t := schnorrkel.NewSigningContext(SigningContext, message)

	return pk.key.Verify(s, t)
}

// Encode returns the 32-byte encoding of the public key
func (pk *PublicKey) Encode() []byte {
	if pk.key == nil {
		return nil
	}

	enc := pk.key.Encode()

	return enc[:]
}

// Decode decodes the input bytes into a public key and sets the receiver the decoded key
// Input must be 32 bytes, or else this function will error
func (pk *PublicKey) Decode(in []byte) error {
	if len(in) != publicKeySize {
		return errors.New("input to sr25519 public key decode is not 32 bytes")
	}

	b := [32]byte{}
	copy(b[:], in)

	pk.key = &schnorrkel.PublicKey{}

	return pk.key.Decode(b)
}

func schnorrkelPublicKeyFromBytes(in []byte) (*schnorrkel.PublicKey, error) {
	if len(in) != publicKeySize {
		return nil, errors.New("input to sr25519 public key decode is not 32 bytes")
	}

	b := [32]byte{}
	copy(b[:], in)

	key := &schnorrkel.PublicKey{}
	err := key.Decode(b)

	return key, err
}

// Convert this public key to a byte array.
func PublicKeyFromBytes(keyBytes []byte) (*PublicKey, error) {
	switch len(keyBytes) {
	case publicKeySize:
		pubKey, err := schnorrkelPublicKeyFromBytes(keyBytes)
		if err != nil {
			return nil, err
		}

		return &PublicKey{pubKey}, nil
	default:
		return nil, errors.New("public key must be 32 bytes")
	}
}
