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
	"io/ioutil"
	"os"
	"strings"

	"github.com/mailchain/mailchain/internal/encoding"
	"github.com/pkg/errors"
)

// HasAddress check for the presence of the address in the store
func (fs FileStore) HasAddress(address []byte) bool {
	fd, err := os.Open(fs.filename(address))
	if err != nil {
		return false
	}
	defer fd.Close()

	return true
}

// GetAddresses list all the address this key store has
func (fs FileStore) GetAddresses() ([][]byte, error) {
	files, err := ioutil.ReadDir(fs.path)
	if err != nil {
		return nil, err
	}
	addresses := [][]byte{}
	for _, f := range files {
		fileName := f.Name()
		if !strings.HasSuffix(fileName, ".json") {
			continue
		}
		fileName = strings.TrimSuffix(fileName, ".json")
		splits := strings.Split(fileName, "/")
		addressPortion := splits[len(splits)-1]
		address, err := hex.DecodeString(addressPortion)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		encryptedKey, err := fs.getEncryptedKey(address)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		addresses = append(addresses, encoding.HexToAddress(encryptedKey.Address))
	}
	return addresses, nil
}
