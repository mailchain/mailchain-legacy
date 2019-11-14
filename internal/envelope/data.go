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
	// Kind0x01 envelope identifier for Message Location Identifier Envelope.
	// The first byte of the envelope is used to identify which programmable envelope is used.
	Kind0x01 byte = 0x01
	// Kind0x50 envelope identifier for Alpha Envelope.
	// The first byte of the envelope is used to identify which programmable envelope is used.
	Kind0x50 byte = 0x50
)

// Data definition for programmable envelopes.
type Data interface {
	// URL returns the addressable location of the message, the URL may be encrypted requiring decrypter to be supplied.
	URL(decrypter cipher.Decrypter) (*url.URL, error)
	// IntegrityHash returns a hash of the encrypted content. This can be used to validate the integrity of the contents before decrypting.
	IntegrityHash(decrypter cipher.Decrypter) ([]byte, error)
	// ContentsHash returns a hash of the decrypted content.
	// This can be used to verify the contents of the message have not been tampered with.
	ContentsHash(decrypter cipher.Decrypter) ([]byte, error)
	// Valid will verify the contents of the envelope.
	// Checks the envelopes contents for no integrity issues which would prevent the envelope from being read.
	Valid() error
}
