package multikey

import (
	"errors"

	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/ed25519"
	"github.com/mailchain/mailchain/crypto/secp256k1"
	"github.com/mailchain/mailchain/crypto/sr25519"
)

func KindFromPublicKey(key crypto.PublicKey) (string, error) {
	switch key.(type) {
	case *ed25519.PublicKey, ed25519.PublicKey:
		return crypto.KindED25519, nil
	case *secp256k1.PublicKey, secp256k1.PublicKey:
		return crypto.KindSECP256K1, nil
	case *sr25519.PublicKey, sr25519.PublicKey:
		return crypto.KindSR25519, nil
	default:
		return "", errors.New("unknown public key type")
	}
}

func KindFromPrivateKey(key crypto.PrivateKey) (string, error) {
	switch key.(type) {
	case *ed25519.PrivateKey, ed25519.PrivateKey:
		return crypto.KindED25519, nil
	case *secp256k1.PrivateKey, secp256k1.PrivateKey:
		return crypto.KindSECP256K1, nil
	case *sr25519.PrivateKey, sr25519.PrivateKey:
		return crypto.KindSR25519, nil
	default:
		return "", errors.New("unknown private key type")
	}
}
