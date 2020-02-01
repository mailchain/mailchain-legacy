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

// SharedSecret computes a secret value from a private / public key pair.
// On sending a message the private key should be an ephemeralKey or generated private key,
// the public key is the recipient public key.
// On reading a message the private key is the recipient private key, the public key is the
// ephemeralKey or generated public key.
func (kx SECP256K1) SharedSecret(privateKey crypto.PrivateKey, publicKey crypto.PublicKey) ([]byte, error) {
	secp256k1PrivateKey, err := kx.privateKey(privateKey)
	if err != nil {
		return nil, ErrSharedSecretGenerate
	}

	secp256k1PublicKey, err := kx.publicKey(publicKey)
	if err != nil {
		return nil, ErrSharedSecretGenerate
	}

	ephemeralPublicKey, _ := kx.publicKey(privateKey.PublicKey())
	if ephemeralPublicKey.X == secp256k1PublicKey.X && ephemeralPublicKey.Y == secp256k1PublicKey.Y {
		return nil, ErrSharedSecretGenerate
	}

	sX, _ := kx.curve.ScalarMult(secp256k1PublicKey.X, secp256k1PublicKey.Y, secp256k1PrivateKey.D.Bytes())

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
