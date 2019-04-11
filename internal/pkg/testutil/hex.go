package testutil

import (
	"encoding/hex"
)

// MustHexDecodeString decodes a hex string. It panics for invalid input.
func MustHexDecodeString(input string) []byte {
	dec, err := hex.DecodeString(input)
	if err != nil {
		panic(err)
	}
	return dec
}
