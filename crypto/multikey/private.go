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
	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/ed25519"
	"github.com/mailchain/mailchain/crypto/secp256k1"
	"github.com/pkg/errors"
)

// PrivateKeyFromBytes returns a private key from `[]byte`.
//
// The function used to create the private key is based on the key type.
// Supported key types are secp256k1, ed25519.
func PrivateKeyFromBytes(keyType string, data []byte) (crypto.PrivateKey, error) {
	switch keyType {
	case crypto.SECP256K1:
		return secp256k1.PrivateKeyFromBytes(data)
	case crypto.ED25519:
		return ed25519.PrivateKeyFromBytes(data)
	default:
		return nil, errors.Errorf("unsupported key type: %q", keyType)
	}
}
