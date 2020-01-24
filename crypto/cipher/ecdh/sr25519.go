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

func (kx SR25519) SharedSecret(ephemeralKey crypto.PrivateKey, recipientKey crypto.PublicKey) ([]byte, error) {
	ephemeralPrivateKey, err := kx.privateKey(ephemeralKey)
	if err != nil {
		return nil, ErrSharedSecretGenerate
	}

	recipientPublicKey, err := kx.publicKey(recipientKey)
	if err != nil {
		return nil, ErrSharedSecretGenerate
	}

	ephemeralPublicKey, _ := kx.publicKey(ephemeralKey.PublicKey())

	if bytes.Equal(ephemeralPublicKey.Bytes(), recipientPublicKey.Bytes()) {
		return nil, ErrSharedSecretGenerate
	}

	sharedSecret, err := sr25519.ExchangeKeys(ephemeralPrivateKey, recipientPublicKey)
	if err != nil {
		return nil, ErrSharedSecretGenerate
	}
	return sharedSecret[:], nil
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
