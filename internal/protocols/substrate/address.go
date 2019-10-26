package substrate

import (
	"github.com/mailchain/mailchain/crypto"
	"github.com/minio/blake2b-simd"
	"github.com/pkg/errors"
)

func SS58AddressFormat(network string, publicKey crypto.PublicKey) ([]byte, error) {
	if publicKey == nil {
		return nil, errors.Errorf("public key must not be nil")
	}

	prefixedKey, err := prefixWithNetwork(network, publicKey)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	hash := blake2b.Sum512(addSS58Prefix(prefixedKey))

	// take first 2 bytes of hash since public key
	return append(prefixedKey, hash[:2]...), nil
}

func addSS58Prefix(pubKey []byte) []byte {
	prefix := []byte("SS58PRE")
	return append(prefix, pubKey...)
}

func prefixWithNetwork(network string, publicKey crypto.PublicKey) ([]byte, error) {
	// https://github.com/paritytech/substrate/wiki/External-Address-Format-(SS58)#address-type defines different prefixes by network
	switch network {
	case EdgewareTestnet:
		// 42 = 0x2a
		return append([]byte{0x2a}, publicKey.Bytes()...), nil
	default:
		return nil, errors.Errorf("unknown address prefix for %q", network)
	}
}
