package scrypt

import (
	"crypto/rand"
	"io"

	"github.com/pkg/errors"
)

// DeriveOptionsBuilder creates the options to derive a key from scrypt.
type DeriveOptionsBuilder func(*DeriveOpts)

type DeriveOpts struct {
	Len        int    `json:"len"`
	N          int    `json:"n"`
	P          int    `json:"p"`
	R          int    `json:"r"`
	Salt       []byte `json:"salt"`
	Passphrase string `json:"-"`
}

func (d DeriveOpts) KDF() string { return "scrypt" }

// WithPassphrase adds passphrase to the dervive options
func WithPassphrase(passphrase string) DeriveOptionsBuilder {
	return func(o *DeriveOpts) { o.Passphrase = passphrase }
}

func RandomSalt() (DeriveOptionsBuilder, error) {
	salt := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, errors.WithMessage(err, "could not generate salt")
	}
	return func(o *DeriveOpts) { o.Salt = salt }, nil
}

func DefaultDeriveOptions() DeriveOptionsBuilder {
	return func(o *DeriveOpts) {
		// N is the N parameter of Scrypt encryption algorithm, using 256MB
		// memory and taking approximately 1s CPU time on a modern processor.
		o.N = 1 << 18
		// P is the P parameter of Scrypt encryption algorithm, using 256MB
		// memory and taking approximately 1s CPU time on a modern processor.
		o.P = 1

		o.R = 8
		o.Len = 32
	}
}

func FromEncryptedKey(len int, n int, p int, r int, salt []byte) DeriveOptionsBuilder {
	return func(o *DeriveOpts) {
		o.Len = len
		o.N = n
		o.P = p
		o.R = r
		o.Salt = salt
	}
}
