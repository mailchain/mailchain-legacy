package ecdh

import (
	"bytes"
	"errors"
	"io"

	"github.com/agl/ed25519/extra25519"
	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/ed25519"
	"golang.org/x/crypto/curve25519"
)

type ED25519 struct {
	rand io.Reader
}

func NewED25519(rand io.Reader) (*ED25519, error) {
	if rand == nil {
		return nil, errors.New("rand must not be nil")
	}

	return &ED25519{rand: rand}, nil
}

func (kx ED25519) EphemeralKey() (crypto.PrivateKey, error) {
	return ed25519.GenerateKey(kx.rand)
}

func (kx ED25519) SharedSecret(ephemeralKey crypto.PrivateKey, recipientKey crypto.PublicKey) ([]byte, error) {
	ephemeralPrivateKey, err := kx.privateKey(ephemeralKey)
	if err != nil {
		return nil, ErrSharedSecretGenerate
	}

	recipientPublicKey, err := kx.publicKey(recipientKey)
	if err != nil {
		return nil, ErrSharedSecretGenerate
	}

	ephemeralPublicKey, _ := kx.publicKey(ephemeralKey.PublicKey())

	if bytes.Equal(ephemeralPublicKey[:], recipientPublicKey[:]) {
		return nil, ErrSharedSecretGenerate
	}

	var secret [32]byte

	curve25519.ScalarMult(&secret, &ephemeralPrivateKey, &recipientPublicKey)

	return secret[:], nil
}

func (kx ED25519) publicKey(pubKey crypto.PublicKey) (key [32]byte, err error) {
	switch pk := pubKey.(type) {
	case *ed25519.PublicKey:
		var ed25519Key, key [32]byte

		copy(ed25519Key[:], pk.Bytes())
		extra25519.PublicKeyToCurve25519(&key, &ed25519Key)

		return key, nil
	default:
		return [32]byte{}, ErrSharedSecretGenerate
	}
}

func (kx ED25519) privateKey(privKey crypto.PrivateKey) (key [32]byte, err error) {
	switch pk := privKey.(type) {
	case *ed25519.PrivateKey:
		var ed25519Key [64]byte

		copy(ed25519Key[:], pk.Bytes())
		extra25519.PrivateKeyToCurve25519(&key, &ed25519Key)

		return key, nil
	default:
		return [32]byte{}, ErrSharedSecretGenerate
	}
}
