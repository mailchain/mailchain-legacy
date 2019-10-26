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
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

// HasAddress check for the presence of the address in the store
func (f FileStore) HasAddress(address []byte) bool {
	fd, err := f.fs.Open(f.filename(address))
	if err != nil {
		return false
	}

	defer fd.Close()

	return true
}

// GetAddresses list all the address this key store has
func (f FileStore) GetAddresses() ([][]byte, error) {
	files, err := afero.ReadDir(f.fs, "./")
	if err != nil {
		return nil, err
	}
	addresses := [][]byte{}

	for _, file := range files {
		fileName := file.Name()
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

		encryptedKey, err := f.getEncryptedKey(address)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		addresses = append(addresses, HexToAddress(encryptedKey.Address))
	}
	return addresses, nil
}

// TODO: this needs to be removed only works for ethereum, using mailchain address package instead, or rethink how keys are stored. Should public key be the index then generate address
func HexToAddress(s string) []byte {
	return FromHex(s)
}

// FromHex returns the bytes represented by the hexadecimal string s.
// s may be prefixed with "0x".
func FromHex(s string) []byte {
	if len(s) > 1 {
		if s[0:2] == "0x" || s[0:2] == "0X" {
			s = s[2:]
		}
	}
	if len(s)%2 == 1 {
		s = "0" + s
	}
	return Hex2Bytes(s)
}

// Hex2Bytes returns the bytes represented by the hexadecimal string str.
func Hex2Bytes(str string) []byte {
	h, _ := hex.DecodeString(str)
	return h
}
