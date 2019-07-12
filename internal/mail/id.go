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

package mail

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/pkg/errors"
)

// NewID create a new secure random ID
func NewID() (ID, error) {
	id, err := generateRandomID(44)
	return id, errors.WithMessage(err, "could not generate ID")
}

// FromHexString create ID from multihash hex string
func FromHexString(h string) (ID, error) {
	return hex.DecodeString(h)
}

// HexString create a multihash representation of ID as hex string
func (id ID) HexString() string {
	return hex.EncodeToString(id)
}

// ID create the mail message ID header
type ID []byte

// generateRandomID returns a securely generated random bytes encoded with multihash 0x00 prefix.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func generateRandomID(n int) (ID, error) {
	bytes := make([]byte, n)
	_, err := rand.Read(bytes)
	return bytes, errors.WithStack(err)
}
