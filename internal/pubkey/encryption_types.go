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

package pubkey

import (
	"errors"

	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/cipher/encrypter"
	"github.com/mailchain/mailchain/crypto/ed25519"
	"github.com/mailchain/mailchain/crypto/secp256k1"
	"github.com/mailchain/mailchain/crypto/sr25519"
)

// EncryptionMethods returns supported encryption methods.
func EncryptionMethods(key crypto.PublicKey) ([]string, error) {
	switch key.(type) {
	case *ed25519.PublicKey:
		return []string{encrypter.NACLECDH, encrypter.NoOperation}, nil
	case *secp256k1.PublicKey:
		return []string{encrypter.AES256CBC, encrypter.NoOperation}, nil
	case *sr25519.PublicKey:
		return []string{encrypter.NACLECDH, encrypter.NoOperation}, nil
	default:
		return nil, errors.New("unsupported public key")
	}
}
