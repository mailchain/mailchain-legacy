package sr25519

import (
	"errors"

	"github.com/developerfred/go-schnorrkel"
	"github.com/mailchain/mailchain/crypto"
)

const (
	publicKeySize   = 32
	signatureLength = 64
)

// PublicKey is a interface
type PublicKey struct {
	key []byte
}

func (pk PublicKey) generate() (*schnorrkel.PublicKey, error) {
	if pk.key == nil {
		return nil, errors.New("cannot create public key: input is not 32 bytes")
	}

	buf := [32]byte{}
	copy(buf[:], pk.key)

	public := schnorrkel.NewPublicKey(buf)

	return &public, nil
}

// Bytes return Publickey Bytes
func (pk PublicKey) Bytes() []byte {
	return pk.key
}

// Kind returns the key type
func (pk PublicKey) Kind() string {
	return crypto.SR25519
}

func newPublicKey(b []byte) PublicKey { //nolint
	return PublicKey{key: b}
}

// Verify uses the sr25519 signature algorithm to verify that the message was signed by
// this public key; it returns true if this key created the signature for the message,
// false otherwise
func (pk PublicKey) Verify(message, sig []byte) bool {
	b := [signatureLength]byte{}
	copy(b[:], sig)

	s := &schnorrkel.Signature{}

	err := s.Decode(b)
	if err != nil {
		return false
	}

	signingContext := schnorrkel.NewSigningContext(SigningContext, message)

	pub, err := pk.generate()
	if err != nil {
		return false
	}

	return pub.Verify(s, signingContext)
}

// PublicKeyFromBytes - Convert byte array to PublicKey
func PublicKeyFromBytes(keyBytes []byte) (crypto.PublicKey, error) {
	switch len(keyBytes) {
	case publicKeySize:
		pubKey := newPublicKey(keyBytes)

		return pubKey, nil
	default:
		return nil, errors.New("public key must be 32 bytes")
	}
}
