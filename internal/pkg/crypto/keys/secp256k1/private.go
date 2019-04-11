package secp256k1

import (
	"crypto/ecdsa"
	"encoding/hex"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/mailchain/mailchain/internal/pkg/crypto/keys"
	"github.com/pkg/errors"
)

// PrivateKey based on the secp256k1 curve
type PrivateKey struct {
	ecdsa ecdsa.PrivateKey
}

// Bytes returns the byte representation of the private key
func (pk PrivateKey) Bytes() []byte {
	return crypto.FromECDSA(&pk.ecdsa)
}

// PublicKey return the public key that is derived from the private key
func (pk PrivateKey) PublicKey() keys.PublicKey {
	return PublicKey{ecdsa: pk.ecdsa.PublicKey}
}

// PrivateKeyFromECDSA get a private key from an ecdsa.PrivateKey
func PrivateKeyFromECDSA(pk ecdsa.PrivateKey) PrivateKey {
	return PrivateKey{ecdsa: pk}
}

// PrivateKeyFromBytes get a private key from []byte
func PrivateKeyFromBytes(pk []byte) (*PrivateKey, error) {
	rpk, err := crypto.ToECDSA(pk)
	if err != nil {
		return nil, errors.Errorf("could not convert private key")
	}
	return &PrivateKey{ecdsa: *rpk}, nil
}

// PrivateKeyFromHex get a private key from hex string
func PrivateKeyFromHex(hexkey string) (*PrivateKey, error) {
	b, err := hex.DecodeString(hexkey)
	if err != nil {
		return nil, errors.New("invalid hex string")
	}
	return PrivateKeyFromBytes(b)
}

// TODO: hang off key object instead
func PrivateKeyToECIES(pk keys.PrivateKey) (*ecies.PrivateKey, error) {
	rpk, err := crypto.ToECDSA(pk.Bytes())
	if err != nil {
		return nil, errors.Errorf("could not convert private key")
	}
	return ecies.ImportECDSA(rpk), nil
}

// TODO: hang off key object instead
func PrivateKeyToECDSA(pk keys.PrivateKey) (*ecdsa.PrivateKey, error) {
	rpk, err := crypto.ToECDSA(pk.Bytes())
	if err != nil {
		return nil, errors.Errorf("could not convert private key")
	}
	return rpk, nil
}
