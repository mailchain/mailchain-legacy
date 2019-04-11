// Copyright (c) 2019 Finobo
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package mail

import (
	"crypto/rand"

	"github.com/multiformats/go-multihash"
	"github.com/pkg/errors"
)

// NewID create a new secure random ID
func NewID() (ID, error) {
	id, err := generateRandomID(44)
	if err != nil {
		return nil, errors.WithMessage(err, "could not generate ID")
	}
	return id, nil
}

// FromHexString create ID from multihash hex string
func FromHexString(hex string) (ID, error) {
	id, err := multihash.FromHexString(hex)
	if err != nil {
		return nil, errors.WithMessage(err, "could not generate ID")
	}
	return ID(id), nil
}

// String create a multihash representation of ID
func (id ID) String() string {
	mh := multihash.Multihash(id)
	return mh.String()
}

// ID create the mail message ID header
type ID multihash.Multihash

// generateRandomID returns a securely generated random bytes encoded with multihash 0x00 prefix.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func generateRandomID(n int) (ID, error) {
	bytes := make([]byte, n)
	_, err := rand.Read(bytes)
	if err != nil {
		return nil, err
	}

	id, err := multihash.Encode(bytes, multihash.ID)
	if err != nil {
		return nil, err
	}
	return id, nil
}
