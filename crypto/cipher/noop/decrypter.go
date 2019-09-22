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

package noop

import (
	"github.com/mailchain/mailchain/crypto/cipher"
)

// NewDecrypter create a new decrypter attaching the private key to it
func NewDecrypter() Decrypter {
	return Decrypter{}
}

// Decrypter will decrypt data using AES256CBC method
type Decrypter struct {
}

// Decrypt data using recipient private key with AES in CBC mode.
func (d Decrypter) Decrypt(data cipher.EncryptedContent) (cipher.PlainContent, error) {
	return cipher.PlainContent(data), nil
}
