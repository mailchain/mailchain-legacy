package factory

import (
	"crypto/rand"
	"io"
)

// NewSeed returns a cryptographically secure random seed. Random seeds are needed
// when creating a new secure private key.

// The length of the seed depends on it's usage. In most cases seed length is between
// 16 and 64 making it (128 to 512 bits).
func NewSeed(length uint8) ([]byte, error) {
	return generateSeed(rand.Reader, length)
}

func generateSeed(r io.Reader, length uint8) ([]byte, error) {
	buf := make([]byte, length)

	if _, err := r.Read(buf); err != nil {
		return nil, err
	}

	return buf, nil
}
