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

package keystore

import (
	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/cipher"
	"github.com/mailchain/mailchain/crypto/cipher/aes256cbc"
	"github.com/mailchain/mailchain/crypto/cipher/nacl"
	"github.com/pkg/errors"
)

// Decrypter use the correct function to get the decrypter from private key
func Decrypter(cipherType byte, pk crypto.PrivateKey) (cipher.Decrypter, error) {
	switch cipherType {
	case cipher.AES256CBC:
		return aes256cbc.NewDecrypter(pk)
	case cipher.NACL:
		return nacl.NewDecrypter(pk)
	default:
		return nil, errors.Errorf("unsupported decrypter type")
	}
}
