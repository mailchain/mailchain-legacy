package algorand

import (
	"bytes"
	"context"

	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/ed25519"
	"github.com/pkg/errors"
)

// NewPublicKeyFinder create a default substrate public key finder.
func NewPublicKeyFinder() *PublicKeyFinder {
	return &PublicKeyFinder{}
}

var (
	errInvalidProtocol      = errors.New("protocol must be 'algorand'")
	errAddressLength        = errors.Errorf("address must be %d bytes in length", checksumLength+ed25519.PublicKeySize)
	errInconsistentChecksum = errors.New("invalid address, checksum verification failed")
)

// PublicKeyFinder for substrate.
type PublicKeyFinder struct {
}

// PublicKeyFromAddress returns the public key from the address.
func (pkf *PublicKeyFinder) PublicKeyFromAddress(ctx context.Context, protocol, network string, address []byte) (crypto.PublicKey, error) {
	if protocol != "algorand" {
		return nil, errInvalidProtocol
	}

	if len(address) != checksumLength+ed25519.PublicKeySize {
		return nil, errAddressLength
	}

	providedAddressChecksum := address[len(address)-checksumLength:]
	publicKeyFromAddress := address[:ed25519.PublicKeySize]
	calculatedChecksum := checksum(publicKeyFromAddress)
	isValid := bytes.Equal(providedAddressChecksum, calculatedChecksum)

	if !isValid {
		return nil, errInconsistentChecksum
	}

	return ed25519.PublicKeyFromBytes(publicKeyFromAddress)
}
