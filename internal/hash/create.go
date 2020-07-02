package hash

import (
	"github.com/ipfs/go-cid"
	"github.com/multiformats/go-multihash"
	"github.com/pkg/errors"
)

func Create(kind int, data []byte) ([]byte, error) {
	switch kind {
	case SHA3256:
		return multihash.Sum(data, multihash.SHA3_256, -1)
	case MurMur3128:
		return multihash.Sum(data, multihash.MURMUR3_128, -1)
	case CIVv1SHA2256Raw:
		// equivalent // pref := cid.Prefix{Version: 1, Codec: cid.Raw, MhType: multihash.SHA2_256, MhLength: -1}; pref.Sum(data)
		h, err := multihash.Sum(data, multihash.SHA2_256, -1)

		return cid.NewCidV1(cid.Raw, h).Bytes(), err
	}

	return nil, errors.Errorf("unknown hash kind")
}
