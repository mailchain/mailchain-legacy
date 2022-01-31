// Copyright 2022 Mailchain Ltd
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

package algorand

import (
	"crypto/sha512"

	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/ed25519"
	"github.com/pkg/errors"
)

const (
	checksumLength = 4
)

// Address returns the address from the public key.
func Address(pubKey crypto.PublicKey) ([]byte, error) {
	switch pubKey.(type) {
	case ed25519.PublicKey, *ed25519.PublicKey:
		return append(pubKey.Bytes(), checksum(pubKey.Bytes())...), nil
	default:
		return nil, errors.Errorf("invalid public key type: %T", pubKey)
	}
}

// Algorand 4-byte checksum. Added to a public key it combines make an algorand address.
func checksum(data []byte) []byte {
	fullHash := sha512.Sum512_256(data)
	return fullHash[len(fullHash)-checksumLength:]
}
