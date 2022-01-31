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

package envelope

import (
	"github.com/mailchain/mailchain/internal/mli"
	"github.com/pkg/errors"
)

// CreateOptionsBuilder creates the options to derive a key from scrypt.
type CreateOptionsBuilder func(*CreateOpts)

// CreateOpts for building an envelope.
type CreateOpts struct {
	// URL of message.
	URL string
	// DecryptedHash use to verify the decrypted contents have not been tampered with.
	DecryptedHash []byte
	// EncryptedHash use to verify the encrypted contents have not been tampered with.
	EncryptedHash []byte
	// Resource id of the message.
	Resource string
	// Kind type of envelope used
	Kind byte
	// Location maps to an addressable location.
	Location uint64
	// EncryptedContents message after its been encrypted.
	EncryptedContents []byte
}

// func (d CreateOpts) Kind() byte { return d.Kind }

// WithKind creates options builder with envelope type identifier.
func WithKind(kind byte) CreateOptionsBuilder {
	return func(o *CreateOpts) { o.Kind = kind }
}

// WithURL creates options builder with an encrypted URL.
func WithURL(address string) CreateOptionsBuilder {
	return func(o *CreateOpts) { o.URL = address }
}

// WithResource creates options builder with a resource location.
func WithResource(resource string) CreateOptionsBuilder {
	return func(o *CreateOpts) { o.Resource = resource }
}

// WithEncryptedContents creates options builder with a the encrypted content of the message.
func WithEncryptedContents(encryptedContents []byte) CreateOptionsBuilder {
	return func(o *CreateOpts) { o.EncryptedContents = encryptedContents }
}

// WithMessageLocationIdentifier creates options builder with a message location identifier.
func WithMessageLocationIdentifier(msgLocInd uint64) (CreateOptionsBuilder, error) {
	_, ok := mli.ToAddress()[msgLocInd]
	if !ok && msgLocInd != 0 {
		return func(o *CreateOpts) {}, errors.Errorf("unknown mli %q", msgLocInd)
	}

	return func(o *CreateOpts) { o.Location = msgLocInd }, nil
}

// WithDecryptedHash creates options builder with the decrypted hash.
func WithDecryptedHash(decryptedHash []byte) CreateOptionsBuilder {
	return func(o *CreateOpts) { o.DecryptedHash = decryptedHash }
}

// WithEncryptedHash creates options builder with the encrypted hash.
func WithEncryptedHash(encryptedHash []byte) CreateOptionsBuilder {
	return func(o *CreateOpts) { o.EncryptedHash = encryptedHash }
}
