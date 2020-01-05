package ecdh

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"errors"
	"io"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/secp256k1"
)

type SECP256K1 struct {
	rand  io.Reader
	curve elliptic.Curve
}

func NewSECP256K1(rand io.Reader) (*SECP256K1, error) {
	if rand == nil {
		return nil, errors.New("rand must not be nil")
	}
	return &SECP256K1{rand: rand, curve: ethcrypto.S256()}, nil
}

func (kx SECP256K1) EphemeralKey() (crypto.PrivateKey, error) {
	return secp256k1.GenerateKey(kx.rand)
}

func (kx SECP256K1) SharedSecret(ephemeralKey crypto.PrivateKey, recipientKey crypto.PublicKey) ([]byte, error) {
	ephemeralPrivateKey, err := kx.privateKey(ephemeralKey)
	if err != nil {
		return nil, ErrSharedSecretGenerate
	}
	recipientPublicKey, err := kx.publicKey(recipientKey)
	if err != nil {
		return nil, ErrSharedSecretGenerate
	}

	ephemeralPublicKey, _ := kx.publicKey(ephemeralKey.PublicKey())
	if ephemeralPublicKey.X == recipientPublicKey.X && ephemeralPublicKey.Y == recipientPublicKey.Y {
		return nil, ErrSharedSecretGenerate
	}

	sX, _ := kx.curve.ScalarMult(recipientPublicKey.X, recipientPublicKey.Y, ephemeralPrivateKey.D.Bytes())

	return sX.Bytes(), nil
}

func (kx SECP256K1) publicKey(pubKey crypto.PublicKey) (*ecdsa.PublicKey, error) {
	switch pk := pubKey.(type) {
	case *secp256k1.PublicKey:
		return pk.ECDSA(), nil
	default:
		return nil, errors.New("unknown public key")
	}
}

func (kx SECP256K1) privateKey(privKey crypto.PrivateKey) (*ecdsa.PrivateKey, error) {
	switch pk := privKey.(type) {
	case *secp256k1.PrivateKey:
		return pk.ECDSA()
	default:
		return nil, errors.New("unknown private key")
	}
}
