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
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/mailchain/mailchain/internal/keystore"
	"github.com/mailchain/mailchain/internal/encoding"
	"github.com/mailchain/mailchain/internal/keystore"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

// NewFileStore create a new filestore with the path specified
func NewFileStore(path string, logger io.Writer) FileStore {
	return FileStore{
		fs:     afero.NewBasePathFs(afero.NewOsFs(), path),
		rand:   rand.Reader,
		logger: logger,
	}
}

// FileStore object
type FileStore struct {
	fs     afero.Fs
	rand   io.Reader
	logger io.Writer
}

func (f FileStore) getEncryptedKeys() ([]keystore.EncryptedKey, error) {
	files, err := afero.ReadDir(f.fs, "./")
	if err != nil {
		return nil, err
	}

	encryptedKeys := []keystore.EncryptedKey{}

	for _, file := range files {
		fileName := file.Name()
		if !strings.HasSuffix(fileName, ".json") {
			fmt.Fprintf(f.logger, "skipping non .json filename %s\n", fileName)
			continue
		}

		fileName = strings.TrimSuffix(fileName, ".json")
		splits := strings.Split(fileName, "/")
		pubKeyPortion := splits[len(splits)-1]

		pubKeyBytes, err := hex.DecodeString(pubKeyPortion)
		if err != nil {
			fmt.Fprintf(f.logger, "skipping invalid filename %s: %v\n", fileName, err)
		}

		encryptedKey, err := f.getEncryptedKey(pubKeyBytes)
		if err != nil {
			fmt.Fprintf(f.logger, "skipping invalid file %s: %v\n", fileName, err)
			continue
		}

		encryptedKeys = append(encryptedKeys, *encryptedKey)
	}

	return encryptedKeys, nil
}

func (f FileStore) getEncryptedKey(pubKeyBytes []byte) (*keystore.EncryptedKey, error) {
	fd, err := f.fs.Open(f.filename(pubKeyBytes))
	if err != nil {
		return nil, errors.WithMessage(err, "could not find key file")
	}

	defer fd.Close()

	encryptedKey := new(keystore.EncryptedKey)
	if err := json.NewDecoder(fd).Decode(encryptedKey); err != nil {
		return nil, err
	}

	if !bytes.Equal(encryptedKey.PublicKeyBytes, pubKeyBytes) {
		return nil, fmt.Errorf("key content mismatch: have pubKey %x, want %x", encryptedKey.PublicKeyBytes, pubKeyBytes)
	}

	return encryptedKey, nil
}

func (f FileStore) filename(pubKeyBytes []byte) string {
	return fmt.Sprintf("%s.json", encoding.EncodeHex(pubKeyBytes))
}
