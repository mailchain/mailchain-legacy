package hash

import "github.com/multiformats/go-multihash"

// CreateIntegrityHash returns a hash of the encrypted `[]byte` to allow easy checking it has not been tampered with.
func CreateIntegrityHash(encryptedData []byte) multihash.Multihash {
	hash, _ := Create(MurMur3128, encryptedData)
	return hash
}
