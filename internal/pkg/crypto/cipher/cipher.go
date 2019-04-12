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

//go:generate mockgen -source=crypto.go -package=mocks -destination=$PACKAGE_PATH/internal/pkg/testutil/mocks/cipher.go
package cipher

import (
	"io"

	"github.com/mailchain/mailchain/internal/pkg/crypto/keys"
)

const (
	AES256CBC byte = 0x2e // TODO: Not merged yet
)

// EncryptedContent typed version of byte array that holds encrypted data
type EncryptedContent []byte

// PlainContent typed version of byte array that holds plain data
type PlainContent []byte

// Decrypter will decrypt data using specified method
type Decrypter interface {
	Decrypt(EncryptedContent) (PlainContent, error)
}

// Encrypter will encrypt data using public key
type Encrypter interface {
	Encrypt(rand io.Reader, pub keys.PublicKey, plain PlainContent) (EncryptedContent, error)
}
