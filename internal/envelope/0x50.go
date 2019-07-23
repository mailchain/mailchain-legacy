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

package envelope

import (
	"encoding/hex"
	"net/url"
	"strings"

	"github.com/mailchain/mailchain/crypto/cipher"
	"github.com/pkg/errors"
)

func (d *ZeroX50) URL(decrypter cipher.Decrypter) (*url.URL, error) {
	loc, err := decrypter.Decrypt(d.EncryptedURL)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return url.Parse(string(loc))
}

func (d *ZeroX50) ContentsHash(decrypter cipher.Decrypter) ([]byte, error) {
	return d.DecryptedHash, nil
}

func (d *ZeroX50) IntegrityHash(decrypter cipher.Decrypter) ([]byte, error) {
	loc, err := decrypter.Decrypt(d.EncryptedURL)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	parts := strings.Split(string(loc), "-")
	if len(parts) < 2 {
		return nil, errors.Errorf("could not safely extract hash from location")
	}
	return hex.DecodeString(parts[len(parts)-1])
}

func (d *ZeroX50) Valid() error {
	if len(d.EncryptedURL) == 0 {
		return errors.Errorf("EncryptedURL must not be empty")
	}

	if len(d.DecryptedHash) == 0 {
		return errors.Errorf("DecryptedHash must not be empty")
	}

	return nil
}
