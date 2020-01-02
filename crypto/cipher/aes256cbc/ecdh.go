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
	"crypto/elliptic"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"io"

	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/pkg/errors"
)

// deriveSharedSecret create a shared secret between public and private key.
func deriveSharedSecret(pub *ecies.PublicKey, private *ecies.PrivateKey) ([]byte, error) {
	x, _ := pub.ScalarMult(pub.X, pub.Y, private.D.Bytes())
	if x == nil {
		return nil, errors.New("failed to derive shared secret")
	}
	return x.Bytes(), nil
}

func generateMacKeyAndEncryptionKey(sharedSecret []byte) (macKey, encryptionKey []byte) {
	hash := sha512.Sum512(sharedSecret)
	encryptionKey = hash[:32]
	macKey = hash[32:]
	return macKey, encryptionKey
}

func (e Encrypter) generateIV() ([]byte, error) {
	iv := make([]byte, 16)
	_, err := io.ReadFull(e.rand, iv)
	return iv, err
}

func generateMac(macKey, iv []byte, ephemeralPublicKey ecies.PublicKey, ciphertext []byte) ([]byte, error) {
	// TODO: curve is hard code yet the type is stored in the keystore. Can aes256cbc work with other curves?
	pub := elliptic.Marshal(ecies.DefaultCurve, ephemeralPublicKey.X, ephemeralPublicKey.Y)
	dataToMac := append(iv, pub...)
	dataToMac = append(dataToMac, ciphertext...)
	mac := hmac.New(sha256.New, macKey)
	_, err := mac.Write(dataToMac)
	if err != nil {
		return nil, err
	}
	return mac.Sum(nil), nil
}
