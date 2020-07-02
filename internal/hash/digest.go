package hash

import (
	"github.com/ipfs/go-cid"
	"github.com/multiformats/go-multihash"
	"github.com/pkg/errors"
)

func GetDigest(kind int, hash []byte) ([]byte, error) {
	switch kind {
	case SHA3256, MurMur3128:
		o, err := multihash.Decode(hash)
		if err != nil {
			return nil, err
		}

		return o.Digest, err
	case CIVv1SHA2256Raw:
		c, err := cid.Cast(hash)
		if err != nil {
			return nil, err
		}

		o, _ := multihash.Decode(c.Hash()) // cast statement tests known error conditions

		return o.Digest, err
	default:
		return nil, errors.Errorf("unknown hash kind")
	}
}
