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
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/mailchain/mailchain/internal/crypto/keys"
	"github.com/mailchain/mailchain/internal/keystore"
	"github.com/mailchain/mailchain/internal/keystore/kdf"
	"github.com/mailchain/mailchain/internal/keystore/kdf/multi"
	"github.com/mailchain/mailchain/internal/keystore/kdf/scrypt"
	"github.com/pkg/errors"
)

// Store the private key with the storage key and curve type
func (fs FileStore) Store(private keys.PrivateKey, curveType string, deriveKeyOptions multi.OptionsBuilders) ([]byte, error) {
	storageKey, keyDefFunc, err := multi.DeriveKey(deriveKeyOptions)
	if err != nil {
		return nil, errors.WithMessage(err, "could not derive storage key")
	}
	encrypted, err := easySeal(private.Bytes(), storageKey)
	if err != nil {
		return nil, errors.WithMessage(err, "could seal storage key")
	}
	address := private.PublicKey().Address()
	keyJSON := keystore.EncryptedKey{
		Address:       hex.EncodeToString(address),
		StorageCipher: "nacl",
		CipherText:    encrypted,
		CurveType:     curveType,
		ID:            uuid.New().String(),
		KDF:           keyDefFunc,
		Timestamp:     time.Now(),
		Version:       "0.0.0", // TODO: get the actual version of this
	}
	if keyDefFunc != kdf.Scrypt {
		return nil, errors.Errorf("kdf not supported")
	}

	keyJSON.ScryptParams = scrypt.CreateOptions(deriveKeyOptions.Scrypt)
	content, err := json.Marshal(keyJSON)
	if err != nil {
		return nil, err
	}
	// Write into temporary file
	tmpName, err := writeTemporaryKeyFile(fs.filename(address), content)
	if err != nil {
		return nil, err
	}
	return address, os.Rename(tmpName, fs.filename(address))
}

func writeTemporaryKeyFile(file string, content []byte) (string, error) {
	// Create the keystore directory with appropriate permissions
	// in case it is not present yet.
	const dirPerm = 0700
	if err := os.MkdirAll(filepath.Dir(file), dirPerm); err != nil {
		return "", err
	}
	// Atomic write: create a temporary hidden file first
	// then move it into place. TempFile assigns mode 0600.
	f, err := ioutil.TempFile(filepath.Dir(file), "."+filepath.Base(file)+".tmp")
	if err != nil {
		return "", err
	}
	if _, err := f.Write(content); err != nil {
		f.Close()
		os.Remove(f.Name())
		return "", err
	}
	f.Close()
	return f.Name(), nil
}
