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
	"fmt"

	"github.com/mailchain/mailchain/internal/keystore"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

// NewFileStore create a new filestore with the path specified
func NewFileStore(path string) FileStore {
	return FileStore{fs: afero.NewBasePathFs(afero.NewOsFs(), path)}
}

// FileStore object
type FileStore struct {
	fs afero.Fs
}

func (f FileStore) getEncryptedKey(address []byte) (*keystore.EncryptedKey, error) {
	fd, err := f.fs.Open(f.filename(address))
	if err != nil {
		return nil, errors.WithMessage(err, "could not find key file")
	}

	defer fd.Close()

	encryptedKey := new(keystore.EncryptedKey)
	if err := json.NewDecoder(fd).Decode(encryptedKey); err != nil {
		return nil, err
	}
	if encryptedKey.Address != hex.EncodeToString(address) {
		return nil, fmt.Errorf("key content mismatch: have address %x, want %x", encryptedKey.Address, hex.EncodeToString(address))
	}

	return encryptedKey, nil
}

func (f FileStore) filename(address []byte) string {
	return fmt.Sprintf("%s.json", hex.EncodeToString(address))
}
