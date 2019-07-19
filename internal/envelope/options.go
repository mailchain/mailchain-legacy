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

import "github.com/pkg/errors"

// CreateOptionsBuilder creates the options to derive a key from scrypt.
type CreateOptionsBuilder func(*CreateOpts)

type CreateOpts struct {
	URL           string
	DecryptedHash []byte
	EncryptedHash []byte
	Resource      string
	// Kind type of envelope used
	Kind byte
	// Location maps to an addressable location.
	Location uint64
}

// func (d CreateOpts) Kind() byte { return d.Kind }

// WithKind adds passphrase to the dervive options
func WithKind(kind byte) CreateOptionsBuilder {
	return func(o *CreateOpts) { o.Kind = kind }
}

// WithURL the encrypted location.
func WithURL(address string) CreateOptionsBuilder {
	return func(o *CreateOpts) { o.URL = address }
}

func WithResource(resource string) CreateOptionsBuilder {
	return func(o *CreateOpts) { o.Resource = resource }
}

func WithMessageLocationIdentifier(locCode uint64) (CreateOptionsBuilder, error) {
	_, ok := MLIToAddress()[locCode]
	if !ok && locCode != 0 {
		return func(o *CreateOpts) {}, errors.Errorf("unknown location code %q", locCode)
	}

	return func(o *CreateOpts) { o.Location = locCode }, nil
}

// WithDecryptedHash the encrypted resource name.
func WithDecryptedHash(decryptedHash []byte) CreateOptionsBuilder {
	return func(o *CreateOpts) { o.DecryptedHash = decryptedHash }
}

// WithencryptedHash the encrypted resource name.
func WithEncryptedHash(encryptedHash []byte) CreateOptionsBuilder {
	return func(o *CreateOpts) { o.EncryptedHash = encryptedHash }
}
