package hash

import "github.com/multiformats/go-multihash"

// CreateMessageHash used to verify if the contents of the message match the hash.
func CreateMessageHash(encodedData []byte) multihash.Multihash {
	// No err: SHA3_256 does not error
	hash, _ := Create(SHA3256, encodedData)
	return hash
}
