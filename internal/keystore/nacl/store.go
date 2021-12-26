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
	"encoding/json"
	"fmt"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/mailchain/mailchain"
	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/multikey"
	"github.com/mailchain/mailchain/internal/keystore"
	"github.com/mailchain/mailchain/internal/keystore/kdf"
	"github.com/mailchain/mailchain/internal/keystore/kdf/multi"
	"github.com/mailchain/mailchain/internal/keystore/kdf/scrypt"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

// Store the private key with the storage key and curve type
func (f FileStore) Store(protocol, network string, private crypto.PrivateKey, deriveKeyOptions multi.OptionsBuilders) (crypto.PublicKey, error) {
	storageKey, keyDefFunc, err := multi.DeriveKey(deriveKeyOptions)
	if err != nil {
		return nil, errors.WithMessage(err, "could not derive storage key")
	}

	encrypted, err := easySeal(private.Bytes(), storageKey, f.rand)
	if err != nil {
		return nil, errors.WithMessage(err, "could seal storage key")
	}

	if keyDefFunc != kdf.Scrypt {
		return nil, errors.Errorf("kdf not supported")
	}

	kind, err := multikey.KindFromPrivateKey(private)
	if err != nil {
		return nil, err
	}

	keyJSON := keystore.EncryptedKey{
		PublicKeyBytes: private.PublicKey().Bytes(),
		StorageCipher:  "nacl",
		CipherText:     encrypted,
		CurveType:      kind,
		ID:             uuid.New().String(),
		KDF:            keyDefFunc,
		Timestamp:      time.Now(),
		Version:        mailchain.Version,
		ScryptParams:   scrypt.CreateOptions(deriveKeyOptions.Scrypt),
	}

	content, err := json.Marshal(keyJSON)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	fileLoc := fmt.Sprintf("./%s/%s/%s", protocol, network, f.filename(private.PublicKey().Bytes()))

	tmpName, err := writeTemporaryKeyFile(f.fs, fileLoc, content)
	if err != nil {
		return nil, err
	}

	return private.PublicKey(), f.fs.Rename(tmpName, fileLoc)
}

func writeTemporaryKeyFile(fs afero.Fs, file string, content []byte) (string, error) {
	// Create the keystore directory with appropriate permissions
	// in case it is not present yet.
	const dirPerm = 0700
	if err := fs.MkdirAll(filepath.Dir(file), dirPerm); err != nil {
		return "", errors.WithStack(err)
	}

	// Atomic write: create a temporary hidden file first
	// then move it into place. TempFile assigns mode 0600.
	f, err := afero.TempFile(fs, filepath.Dir(file), "."+filepath.Base(file)+".tmp")
	if err != nil {
		return "", errors.WithStack(err)
	}

	if _, err := f.Write(content); err != nil {
		f.Close()

		return "", fs.Remove(f.Name())
	}

	f.Close()

	return f.Name(), nil
}
