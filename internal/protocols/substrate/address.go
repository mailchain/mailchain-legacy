package substrate

import (
	"github.com/minio/blake2b-simd"
	"github.com/pkg/errors"
)

func SS58AddressFormat(network string, publicKey []byte) ([]byte, error) {
	if len(publicKey) != 32 {
		return nil, errors.Errorf("public key must be 32 bytes")
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

func prefixWithNetwork(network string, publicKey []byte) ([]byte, error) {
	// https://github.com/paritytech/substrate/wiki/External-Address-Format-(SS58)#address-type defines different prefixes by network
	switch network {
	case EdgewareTestnet:
		// 42 = 0x2a
		return append([]byte{0x2a}, publicKey...), nil
	default:
		return nil, errors.Errorf("unknown address prefix for %q", network)
	}
}
