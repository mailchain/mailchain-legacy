package substrate

import (
	"context"

	"github.com/pkg/errors"
)

func NewPublicKeyFinder() *PublicKeyFinder {
	return &PublicKeyFinder{}
}

type PublicKeyFinder struct {
	supportedNetworks []string
}

func (pkf *PublicKeyFinder) PublicKeyFromAddress(ctx context.Context, protocol, network string, address []byte) ([]byte, error) {
	if protocol != "substrate" {
		return nil, errors.New("protocol must be 'substrate'")
	}
	if len(address) != 35 {
		return nil, errors.New("address must be 35 bytes in length")
	}

	// Remove the 1st byte (network identifier)
	// Remove last 2 bytes (blake2b hash)
	newAddress := address[1:33]

	return newAddress, nil
}
