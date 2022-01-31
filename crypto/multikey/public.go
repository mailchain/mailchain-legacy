// Copyright 2022 Mailchain Ltd.
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
	"github.com/mailchain/mailchain/crypto/sr25519"
	"github.com/pkg/errors"
)

// PublicKeyFromBytes use the correct function to get the private key from bytes
func PublicKeyFromBytes(keyType string, data []byte) (crypto.PublicKey, error) {
	switch keyType {
	case crypto.KindSECP256K1:
		return secp256k1.PublicKeyFromBytes(data)
	case crypto.KindED25519:
		return ed25519.PublicKeyFromBytes(data)
	case crypto.KindSR25519:
		return sr25519.PublicKeyFromBytes(data)
	default:
		return nil, errors.Errorf("unsupported curve type")
	}
}
