package ecdh

import (
	"bytes"
	"errors"
	"io"

	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/sr25519"
)

type SR25519 struct {
	rand io.Reader
}

func NewSR25519(rand io.Reader) (*SR25519, error) {
	if rand == nil {
		return nil, errors.New("rand must not be nil")
	}

	return &SR25519{rand: rand}, nil
}

func (kx SR25519) EphemeralKey() (crypto.PrivateKey, error) {
	return sr25519.GenerateKey(kx.rand)
}

// SharedSecret computes a secret value from a private / public key pair.
// On sending a message the private key should be an ephemeralKey or generated private key,
// the public key is the recipient public key.
// On reading a message the private key is the recipient private key, the public key is the
// ephemeralKey or generated public key.
func (kx SR25519) SharedSecret(privateKey crypto.PrivateKey, publicKey crypto.PublicKey) ([]byte, error) {
	sr25519PrivateKey, err := kx.privateKey(privateKey)
	if err != nil {
		return nil, ErrSharedSecretGenerate
	}

	sr25519PublicKey, err := kx.publicKey(publicKey)
	if err != nil {
		return nil, ErrSharedSecretGenerate
	}

	ephemeralPublicKey, _ := kx.publicKey(privateKey.PublicKey())

	if bytes.Equal(ephemeralPublicKey.Bytes(), sr25519PublicKey.Bytes()) {
		return nil, ErrSharedSecretGenerate
	}

	sharedSecret, err := sr25519.ExchangeKeys(sr25519PrivateKey, sr25519PublicKey)
	if err != nil {
		return nil, ErrSharedSecretGenerate
	}

	return sharedSecret, nil
}

func (kx SR25519) publicKey(pubKey crypto.PublicKey) (*sr25519.PublicKey, error) {
	switch pk := pubKey.(type) {
	case *sr25519.PublicKey:
		return pk, nil
	default:
		return nil, ErrSharedSecretGenerate
	}
}

func (kx SR25519) privateKey(privKey crypto.PrivateKey) (*sr25519.PrivateKey, error) {
	switch pk := privKey.(type) {
	case *sr25519.PrivateKey:
		return pk, nil
	default:
		return nil, ErrSharedSecretGenerate
	}
}
