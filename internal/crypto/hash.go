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

package crypto

import (
	"github.com/multiformats/go-multihash"
	"github.com/pkg/errors"
)

func CreateLocationHash(encryptedData []byte) multihash.Multihash {
	hash, _ := createLocationHash(encryptedData)
	return hash
}

func createLocationHash(encryptedData []byte) (multihash.Multihash, error) {
	hash, err := multihash.Sum(encryptedData, multihash.MURMUR3, -1)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to create hash")
	}
	casted, err := multihash.Cast(hash)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to cast hash")
	}
	return casted, nil
}

// CreateMessageHash used to verify if the contents of the message match the hash.
func CreateMessageHash(encodedData []byte) (multihash.Multihash, error) {
	hash, err := multihash.Sum(encodedData, multihash.SHA3_256, -1)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to create hash")
	}
	casted, err := multihash.Cast(hash)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to cast hash")
	}
	return casted, nil
}
