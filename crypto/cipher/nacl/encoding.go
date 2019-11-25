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
	"github.com/mailchain/mailchain/crypto/cipher"
	"github.com/pkg/errors"
)

// bytesEncode encode the encrypted data to the hex format
func bytesEncode(data cipher.EncryptedContent) cipher.EncryptedContent {
	encodedData := make(cipher.EncryptedContent, 1+len(data))
	encodedData[0] = cipher.NACL
	copy(encodedData[1:], data)

	return encodedData
}

// bytesDecode convert the hex format in to the encrypted data format
func bytesDecode(raw cipher.EncryptedContent) (cipher.EncryptedContent, error) {
	if raw[0] != cipher.NACL {
		return nil, errors.Errorf("invalid prefix")
	}

	return raw[1:], nil
}
