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
	"io"

	"github.com/pkg/errors"
	"golang.org/x/crypto/nacl/secretbox"
)

const nonceSize = 24
const secretKeySize = 32

func easyOpen(box, key []byte) ([]byte, error) {
	var secretKey [secretKeySize]byte

	if len(key) != secretKeySize {
		return nil, errors.New("secretbox: key length must be 32")
	}

	if len(box) < nonceSize {
		return nil, errors.New("secretbox: message too short")
	}

	decryptNonce := new([nonceSize]byte)
	copy(decryptNonce[:], box[:nonceSize])
	copy(secretKey[:], key)

	decrypted, ok := secretbox.Open([]byte{}, box[nonceSize:], decryptNonce, &secretKey)
	if !ok {
		return nil, errors.New("secretbox: could not decrypt data with private key")
	}

	return decrypted, nil
}

func easySeal(message, key []byte, rand io.Reader) ([]byte, error) {
	nonce := new([nonceSize]byte)
	if _, err := rand.Read(nonce[:]); err != nil {
		return nil, err
	}

	var secretKey [secretKeySize]byte

	copy(secretKey[:], key)

	return secretbox.Seal(nonce[:], message, nonce, &secretKey), nil
}
