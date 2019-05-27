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

package aes256cbc

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/subtle"

	"github.com/andreburgaud/crypt2go/padding"
	mc "github.com/mailchain/mailchain/internal/crypto/cipher"
	"github.com/mailchain/mailchain/internal/crypto/keys"
	"github.com/mailchain/mailchain/internal/crypto/keys/secp256k1"
	"github.com/pkg/errors"
)

// NewDecrypter create a new decrypter attaching the private key to it
func NewDecrypter(privateKey keys.PrivateKey) Decrypter {
	return Decrypter{privateKey: &privateKey}
}

// Decrypter will decrypt data using AES256CBC method
type Decrypter struct {
	privateKey *keys.PrivateKey
}

// Decrypt data using recipient private key with AES in CBC mode.
func (d Decrypter) Decrypt(data mc.EncryptedContent) (mc.PlainContent, error) {
	encryptedData, err := bytesDecode(data)
	if err != nil {
		return nil, errors.WithMessage(err, "could not convert encryptedData")
	}

	return decryptEncryptedData(*d.privateKey, encryptedData)
}

func decryptEncryptedData(privateKey keys.PrivateKey, data *encryptedData) ([]byte, error) {
	tmpEphemeralPublicKey, err := secp256k1.PublicKeyFromBytes(data.EphemeralPublicKey)
	if err != nil {
		return nil, errors.WithMessage(err, "could not convert ephemeralPublicKey")
	}
	ephemeralPublicKey, err := secp256k1.PublicKeyToECIES(tmpEphemeralPublicKey)
	if err != nil {
		return nil, errors.WithMessage(err, "could not convert to ecies")
	}

	rpk, err := secp256k1.PrivateKeyToECIES(privateKey)
	if err != nil {
		return nil, errors.WithMessage(err, "could not convert private key")
	}

	sharedSecret, err := deriveSharedSecret(ephemeralPublicKey, rpk)
	if err != nil {
		return nil, errors.WithMessage(err, "could not derive shared secret")
	}
	macKey, encryptionKey := generateMacKeyAndEncryptionKey(sharedSecret)
	mac, err := generateMac(macKey, data.InitializationVector, *ephemeralPublicKey, data.Ciphertext)
	if err != nil {
		return nil, errors.WithMessage(err, "generateMac failed")
	}
	if subtle.ConstantTimeCompare(data.MessageAuthenticationCode, mac) != 1 {
		return nil, errors.Errorf("invalid mac")
	}
	return decryptCBC(encryptionKey, data.InitializationVector, data.Ciphertext)
}

func decryptCBC(key, iv, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	plaintext := make([]byte, len(ciphertext))
	cbc := cipher.NewCBCDecrypter(block, iv)
	cbc.CryptBlocks(plaintext, ciphertext)

	plaintext, err = padding.NewPkcs7Padding(block.BlockSize()).Unpad(plaintext)
	if err != nil {
		return nil, errors.WithMessage(err, "could not pad")
	}

	ret := make([]byte, len(plaintext))
	copy(ret, plaintext)
	return ret, nil
}
