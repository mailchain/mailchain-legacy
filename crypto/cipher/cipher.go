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

// Package cipher collects common cryptographic constants and interfaces.
package cipher //go:generate mockgen -source=cipher.go -package=ciphertest -destination=./ciphertest/cipher_mock.go

import (
	"github.com/mailchain/mailchain/crypto"
)

const (
	// NoOperation identified for Encrypt and Decrypter in noop package.
	NoOperation byte = 0x20
	// NACL identified for Encrypt and Decrypter in nacl package.
	NACL byte = 0x2a
	// AES256CBC identified for Encrypt and Decrypter in aes256cbc package.
	AES256CBC byte = 0x2e
)

// EncryptedContent typed version of byte array that holds encrypted data.
//
// Encrypt method returns the encrypted contents as EncryptedContent.
// Decrypt method accepts EncryptedContent as the encrypted contents to decrypt.
type EncryptedContent []byte

// PlainContent typed version of byte array that holds plain data.
//
// Encrypt method returns the encrypted contents as EncryptedContent.
// Decrypt method accepts EncryptedContent as the encrypted contents to decrypt.
type PlainContent []byte

// A Decrypter uses the PrivateKey to decrypt the supplied data.
//
// The decryption method used is dependant on the implementation and
// must check that the data can be decrypted before continuing.
// Returned data should be the plain bytes that were supplied
// originally to the Encrypter.
type Decrypter interface {
	Decrypt(EncryptedContent) (PlainContent, error)
}

// An Encrypter uses the PublicKey to encrypt the supplied data.
//
// The encryption method used is dependant on the implementation and must be included in the response.
// Returned encrypted data must include what encryption method was used as the first byte.
// The data can be decrypted using the corresponding PrivateKey and Decrypter method.
type Encrypter interface {
	Encrypt(pubKey crypto.PublicKey, plain PlainContent) (EncryptedContent, error)
}
