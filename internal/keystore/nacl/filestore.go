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
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/mailchain/mailchain/encoding"
	"github.com/mailchain/mailchain/internal/keystore"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/afero"
)

// NewFileStore create a new filestore with the path specified
func NewFileStore(path string) *FileStore {
	return &FileStore{
		fs:   afero.NewBasePathFs(afero.NewOsFs(), path),
		rand: rand.Reader,
	}
}

// FileStore object
type FileStore struct {
	fs   afero.Fs
	rand io.Reader
}

func (f FileStore) getEncryptedKeys(protocol, network string) ([]keystore.EncryptedKey, error) {
	logger := log.With().Str("component", "nacl-filestore").Str("action", "getEncryptedKeys").Logger()

	// directory needs to be created if looking for a key before one has been added and sub directories created
	const dirPerm = 0700
	if err := f.fs.MkdirAll(fmt.Sprintf("./%s/%s", protocol, network), dirPerm); err != nil {
		return nil, errors.WithStack(err)
	}

	files, err := afero.ReadDir(f.fs, fmt.Sprintf("./%s/%s", protocol, network))
	if err != nil {
		return nil, err
	}

	encryptedKeys := []keystore.EncryptedKey{}

	for _, file := range files {
		fileName := file.Name()
		if !strings.HasSuffix(fileName, ".json") {
			logger.Warn().Msgf("skipping non .json filename %s", fileName)

			continue
		}

		fileName = strings.TrimSuffix(fileName, ".json")
		splits := strings.Split(fileName, "/")
		pubKeyPortion := splits[len(splits)-1]

		pubKeyBytes, err := encoding.DecodeHex(pubKeyPortion)
		if err != nil {
			logger.Warn().Err(err).Msgf("invalid filename %s", fileName)

			continue
		}

		encryptedKey, err := f.getEncryptedKey(protocol, network, pubKeyBytes)
		if err != nil {
			logger.Warn().Err(err).Msgf("invalid file %s", file.Name())

			continue
		}

		encryptedKeys = append(encryptedKeys, *encryptedKey)
	}

	return encryptedKeys, nil
}

func (f FileStore) getEncryptedKey(protocol, network string, pubKeyBytes []byte) (*keystore.EncryptedKey, error) {
	fd, err := f.fs.Open(fmt.Sprintf("./%s/%s/%s", protocol, network, f.filename(pubKeyBytes)))
	if err != nil {
		return nil, errors.WithMessage(err, "could not find key file")
	}

	defer fd.Close()

	encryptedKey := new(keystore.EncryptedKey)
	if err := json.NewDecoder(fd).Decode(encryptedKey); err != nil {
		return nil, errors.WithStack(err)
	}

	if !bytes.Equal(encryptedKey.PublicKeyBytes, pubKeyBytes) {
		return nil, errors.Errorf("key content mismatch: have pubKey %x, want %x", encryptedKey.PublicKeyBytes, pubKeyBytes)
	}

	return encryptedKey, nil
}

func (f FileStore) filename(pubKeyBytes []byte) string {
	return fmt.Sprintf("%s.json", encoding.EncodeHex(pubKeyBytes))
}
