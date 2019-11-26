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

package multikey

import (
	"errors"

	"github.com/mailchain/mailchain/crypto"
)

var (
	// ErrInconclusive is returned when multiple public keys matches for the same input.
	ErrInconclusive = errors.New("multiple matches found")

	// ErrNoMatch is returned when no public key matches for the input.
	ErrNoMatch = errors.New("no match found")
)

// KeyKindsFromSignature tries to determine the key type from the pubKey, message, sig bytes combination.
// The key kinds against which the function should match are specified in the keyKinds slice.
func KeyKindFromSignature(pubKey, message, sig []byte, keyKinds []string) (crypto.PublicKey, error) {
	matches := make([]crypto.PublicKey, 0, 1)

	keyKinds = removeDuplicates(keyKinds)
	for _, kind := range keyKinds {
		key, err := PublicKeyFromBytes(kind, pubKey)
		if err != nil {
			// skip invalid key type.
			continue
		}

		if key.Verify(message, sig) {
			matches = append(matches, key)
		}
	}

	switch len(matches) {
	case 0:
		return nil, ErrNoMatch
	case 1:
		return matches[0], nil
	default:
		return nil, ErrInconclusive
	}
}

func removeDuplicates(x []string) []string {
	if x == nil {
		return nil
	}

	set := make(map[string]struct{})
	unique := []string{}

	for _, str := range x {
		if _, ok := set[str]; !ok {
			set[str] = struct{}{}

			unique = append(unique, str)
		}
	}

	return unique
}
