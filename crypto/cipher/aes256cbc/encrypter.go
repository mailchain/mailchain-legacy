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
	"crypto/elliptic"
	"crypto/rand"
	"io"

	"github.com/andreburgaud/crypt2go/padding"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/mailchain/mailchain/crypto"
	mc "github.com/mailchain/mailchain/crypto/cipher"
	"github.com/pkg/errors"
)

// NewEncrypter create a new encrypter with crypto rand for reader
// and attaching the public key to the encrypter.
func NewEncrypter(pubKey crypto.PublicKey) (*Encrypter, error) {
	return &Encrypter{rand: rand.Reader, publicKey: pubKey}, nil
}

// Encrypter will encrypt data using AES256CBC method
type Encrypter struct {
	rand      io.Reader
	publicKey crypto.PublicKey
}

// Encrypt data using recipient public key with AES in CBC mode.  Generate an ephemeral private key and IV.
func (e Encrypter) Encrypt(message mc.PlainContent) (mc.EncryptedContent, error) {
	epk, err := asPublicECIES(e.publicKey)
	if err != nil {
		return nil, errors.WithMessage(err, "could not convert")
	}

	ephemeral, err := ecies.GenerateKey(e.rand, ecies.DefaultCurve, nil)
	if err != nil {
		return nil, errors.WithMessage(err, "could not generate ephemeral key")
	}

	iv, err := e.generateIV()
	if err != nil {
		return nil, errors.WithMessage(err, "could not generate iv")
	}

	encryptedData, err := encrypt(ephemeral, epk, message, iv)
	if err != nil {
		return nil, errors.WithMessage(err, "could not encrypt data")
	}

	return bytesEncode(encryptedData)
}

func encrypt(ephemeralPrivateKey *ecies.PrivateKey, pub *ecies.PublicKey, input, iv []byte) (*encryptedData, error) {
	ephemeralPublicKey := ephemeralPrivateKey.PublicKey
	sharedSecret, err := deriveSharedSecret(pub, ephemeralPrivateKey)
	if err != nil {
		return nil, err
	}
	macKey, encryptionKey := generateMacKeyAndEncryptionKey(sharedSecret)
	ciphertext, err := encryptCBC(input, iv, encryptionKey)
	if err != nil {
		return nil, errors.WithMessage(err, "encryptCBC failed")
	}

	mac, err := generateMac(macKey, iv, ephemeralPublicKey, ciphertext)
	if err != nil {
		return nil, errors.WithMessage(err, "generateMac failed")
	}

	return &encryptedData{
		MessageAuthenticationCode: mac,
		InitializationVector:      iv,
		EphemeralPublicKey:        elliptic.Marshal(ecies.DefaultCurve, ephemeralPublicKey.X, ephemeralPublicKey.Y),
		Ciphertext:                ciphertext,
	}, nil
}

func encryptCBC(data, iv, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	data, err = padding.NewPkcs7Padding(block.BlockSize()).Pad(data)
	if err != nil {
		return nil, errors.WithMessage(err, "could not pad")
	}

	if len(iv) != block.BlockSize() {
		return nil, errors.Errorf("cipher.NewCBCEncrypter: IV length must equal block size")
	}

	ciphertext := make([]byte, len(data))
	cbc := cipher.NewCBCEncrypter(block, iv)
	cbc.CryptBlocks(ciphertext, data)

	return ciphertext, nil
}
