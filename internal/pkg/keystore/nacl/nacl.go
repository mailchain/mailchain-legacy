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

package nacl

import (
	"github.com/kevinburke/nacl/secretbox"
	"github.com/mailchain/mailchain/internal/pkg/crypto/keys"
	"github.com/mailchain/mailchain/internal/pkg/crypto/keys/multikey"
	"github.com/mailchain/mailchain/internal/pkg/keystore"
	"github.com/mailchain/mailchain/internal/pkg/keystore/kdf/multi"
	"github.com/mailchain/mailchain/internal/pkg/keystore/kdf/scrypt"
	"github.com/pkg/errors"
)

func easyOpen(box, key []byte) ([]byte, error) {
	if len(box) < 24 {
		return nil, errors.New("secretbox: message too short")
	}
	decryptNonce := new([24]byte)
	copy(decryptNonce[:], box[:24])

	var secretKey [32]byte
	copy(secretKey[:], key)

	decrypted, ok := secretbox.Open([]byte{}, box[24:], decryptNonce, &secretKey)
	if !ok {
		return nil, errors.New("secretbox: Could not decrypt invalid input")
	}
	return decrypted, nil
}

func deriveKey(ek *keystore.EncryptedKey, deriveKeyOptions multi.OptionsBuilders) ([]byte, error) {
	switch ek.KDF {
	case "scrypt":
		if ek.ScryptParams == nil {
			return nil, errors.New("scryptParams are required")
		}
		storageOpts := scrypt.FromEncryptedKey(ek.ScryptParams.Len, ek.ScryptParams.N, ek.ScryptParams.P, ek.ScryptParams.R, ek.ScryptParams.Salt)

		return scrypt.DeriveKey(append(deriveKeyOptions.Scrypt, storageOpts))
	default:
		return nil, errors.New("KDF is not supported")
	}
}

func (fs FileStore) getPrivateKey(address []byte, deriveKeyOptions multi.OptionsBuilders) (keys.PrivateKey, error) {
	encryptedKey, err := fs.getEncryptedKey(address)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	storageKey, err := deriveKey(encryptedKey, deriveKeyOptions)
	if err != nil {
		return nil, errors.WithMessage(err, "storage key could not be derived")
	}
	pkBytes, err := easyOpen(encryptedKey.CipherText, storageKey)
	if err != nil {
		return nil, errors.WithMessage(err, "could not decrypt key file")
	}
	pk, err := multikey.PrivateKeyFromBytes(encryptedKey.CurveType, pkBytes)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return pk, nil
}
