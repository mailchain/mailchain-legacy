package multikey

import (
	"github.com/mailchain/mailchain/internal/pkg/crypto/keys"
	"github.com/mailchain/mailchain/internal/pkg/crypto/keys/secp256k1"
	"github.com/pkg/errors"
)

// PrivateKeyFromHex get private key from hex.
func PrivateKeyFromHex(hex string, keyType string) (keys.PrivateKey, error) {
	f, ok := privateKeyFromHexTable[keyType]
	if !ok {
		return nil, errors.Errorf("func for key type %v not registered", keyType)
	}
	return f(hex)
}

type privateKeyFromHex func(hex string) (keys.PrivateKey, error)

// privateKeyFromHexTable maps key types values to hash functions.
var privateKeyFromHexTable = make(map[string]privateKeyFromHex)

func init() {
	registerPrivateKeyFromHex(SECP256K1, func(hex string) (keys.PrivateKey, error) {
		return secp256k1.PrivateKeyFromHex(hex)
	})
}

func registerPrivateKeyFromHex(keyType string, privateKeyFromHex privateKeyFromHex) error {
	_, ok := privateKeyFromHexTable[keyType]
	if ok {
		return errors.Errorf("func for key type %v already registered", keyType)
	}

	privateKeyFromHexTable[keyType] = privateKeyFromHex
	return nil
}
