package crypto

import (
	"github.com/multiformats/go-multihash"
	"github.com/pkg/errors"
)

func CreateLocationHash(encryptedData []byte) (multihash.Multihash, error) {
	hash, err := multihash.Sum(encryptedData, multihash.MURMUR3, -1)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to create hash")
	}
	casted, err := multihash.Cast(hash)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to cast hash")
	}
	return casted, nil
}

// CreateMessageHash used to verify if the contents of the message match the hash.
func CreateMessageHash(encodedData []byte) (multihash.Multihash, error) {
	hash, err := multihash.Sum(encodedData, multihash.SHA3_256, -1)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to create hash")
	}
	casted, err := multihash.Cast(hash)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to cast hash")
	}
	return casted, nil
}
