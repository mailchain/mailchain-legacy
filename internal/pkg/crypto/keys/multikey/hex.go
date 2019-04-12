package multikey

import (
	"github.com/mailchain/mailchain/internal/pkg/crypto/keys"
	"github.com/mailchain/mailchain/internal/pkg/crypto/keys/secp256k1"
	"github.com/pkg/errors"
)

// PrivateKeyFromHex get private key from hex.
func PrivateKeyFromHex(hex, keyType string) (keys.PrivateKey, error) {
	table := map[string]privateKeyFromHex{
		SECP256K1: func(hex string) (keys.PrivateKey, error) {
			return secp256k1.PrivateKeyFromHex(hex)
		},
	}

	f, ok := table[keyType]
	if !ok {
		return nil, errors.Errorf("func for key type %v not registered", keyType)
	}
	return f(hex)
}

type privateKeyFromHex func(hex string) (keys.PrivateKey, error)
