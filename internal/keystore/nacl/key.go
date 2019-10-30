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
	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/multikey"
	"github.com/mailchain/mailchain/internal/keystore"
	"github.com/mailchain/mailchain/internal/keystore/kdf"
	"github.com/mailchain/mailchain/internal/keystore/kdf/multi"
	"github.com/mailchain/mailchain/internal/keystore/kdf/scrypt"
	"github.com/pkg/errors"
)

func deriveKey(ek *keystore.EncryptedKey, deriveKeyOptions multi.OptionsBuilders) ([]byte, error) {
	if ek == nil {
		return nil, errors.New("encryptedKey must not be nil")
	}

	switch ek.KDF {
	case kdf.Scrypt:
		if ek.ScryptParams == nil {
			return nil, errors.New("scryptParams are required")
		}
		storageOpts := scrypt.FromEncryptedKey(ek.ScryptParams.Len, ek.ScryptParams.N, ek.ScryptParams.P, ek.ScryptParams.R, ek.ScryptParams.Salt)

		return scrypt.DeriveKey(append(deriveKeyOptions.Scrypt, storageOpts))
	default:
		return nil, errors.New("KDF is not supported")
	}
}

func (f FileStore) getPrivateKey(encryptedKey *keystore.EncryptedKey, deriveKeyOptions multi.OptionsBuilders) (crypto.PrivateKey, error) {
	if encryptedKey == nil {
		return nil, errors.Errorf("encryptedKey must not be nil")
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
