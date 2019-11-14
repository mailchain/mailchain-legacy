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

package substrate

import (
	"context"

	"github.com/pkg/errors"
)

// NewPublicKeyFinder create a default substrate public key finder.
func NewPublicKeyFinder() *PublicKeyFinder {
	return &PublicKeyFinder{}
}

// PublicKeyFinder for substrate.
type PublicKeyFinder struct {
}

// PublicKeyFromAddress returns the public key from the address.
func (pkf *PublicKeyFinder) PublicKeyFromAddress(ctx context.Context, protocol, network string, address []byte) ([]byte, error) {
	if protocol != "substrate" {
		return nil, errors.New("protocol must be 'substrate'")
	}

	if len(address) != 35 {
		return nil, errors.New("address must be 35 bytes in length")
	}

	// Remove the 1st byte (network identifier)
	// Remove last 2 bytes (blake2b hash)
	newAddress := address[1:33]

	return newAddress, nil
}
