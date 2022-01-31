// Copyright 2022 Mailchain Ltd.
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
	"github.com/mailchain/mailchain/crypto/cipher"
	"github.com/mailchain/mailchain/internal/keystore"
	"github.com/mailchain/mailchain/internal/keystore/kdf/multi"
	"github.com/pkg/errors"
)

// GetDecrypter the decrypter for the specified address.
func (f FileStore) GetDecrypter(address []byte, protocol, network string, decrypterType byte, deriveKeyOptions multi.OptionsBuilders) (cipher.Decrypter, error) {
	encryptedKey, err := f.getEncryptedKeyByAddress(address, protocol, network)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	pk, err := f.getPrivateKey(encryptedKey, deriveKeyOptions)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return keystore.Decrypter(decrypterType, pk)
}
