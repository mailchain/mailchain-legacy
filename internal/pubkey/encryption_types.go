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

package pubkey

import (
	"fmt"

	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/cipher/encrypter"
)

// EncryptionMethods returns supported encryption methods.
func EncryptionMethods(kind string) ([]string, error) {
	switch kind {
	case crypto.ED25519:
		return []string{encrypter.NACL, encrypter.NoOperation}, nil
	case crypto.SECP256K1:
		return []string{encrypter.AES256CBC, encrypter.NoOperation}, nil
	default:
		return nil, fmt.Errorf("%q unsuported public key type", kind)
	}
}
