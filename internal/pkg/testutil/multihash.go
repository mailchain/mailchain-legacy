package testutil

import (
	"github.com/multiformats/go-multihash"
)

// MustHexDecodeMultiHashID takes an hex string and forces it to be parsed with multihash
func MustHexDecodeMultiHashID(input string) multihash.Multihash {
	ret, err := multihash.FromHexString(input)
	if err != nil {
		panic(err)
	}
	return ret
}
