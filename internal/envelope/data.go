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
	"net/url"

	"github.com/mailchain/mailchain/crypto/cipher"
)

const (
	Kind0x01 byte = 0x01 // Message locator
	Kind0x50 byte = 0x50 // Alpha
)

type Data interface {
	URL(decrypter cipher.Decrypter) (*url.URL, error)
	IntegrityHash(decrypter cipher.Decrypter) ([]byte, error)
	ContentsHash(decrypter cipher.Decrypter) ([]byte, error)
	Valid() error
}
