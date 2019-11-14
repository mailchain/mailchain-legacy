// Copyright 2019 Finobo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package scrypt

import (
	"crypto/rand"
	"io"

	"github.com/pkg/errors"
)

// DeriveOptionsBuilder creates the options to derive a key from scrypt.
type DeriveOptionsBuilder func(*DeriveOpts)

// DeriveOpts available for scrypt key derivation.
type DeriveOpts struct {
	Len        int    `json:"len"`
	N          int    `json:"n"`
	P          int    `json:"p"`
	R          int    `json:"r"`
	Salt       []byte `json:"salt"`
	Passphrase string `json:"-"`
}

// KDF name.
func (d DeriveOpts) KDF() string { return "scrypt" }

// WithPassphrase adds passphrase to the dervive options
func WithPassphrase(passphrase string) DeriveOptionsBuilder {
	return func(o *DeriveOpts) { o.Passphrase = passphrase }
}

// RandomSalt use a secure random number to use as the salt.
func RandomSalt() (DeriveOptionsBuilder, error) {
	salt := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, errors.WithMessage(err, "could not generate salt")
	}
	return func(o *DeriveOpts) { o.Salt = salt }, nil
}

// DefaultDeriveOptions for deriving an encryption.
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

// FromEncryptedKey generate the options builder from an encrypted key.
func FromEncryptedKey(length, n, p, r int, salt []byte) DeriveOptionsBuilder {
	return func(o *DeriveOpts) {
		o.Len = length
		o.N = n
		o.P = p
		o.R = r
		o.Salt = salt
	}
}
