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
	"bytes"

	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/multikey"
	"github.com/mailchain/mailchain/internal/address"
	"github.com/mailchain/mailchain/internal/keystore"
	"github.com/pkg/errors"
)

// HasAddress check for the presence of the address in the store
func (f FileStore) HasAddress(searchAddress []byte, protocol, network string) bool {
	_, err := f.getEncryptedKeyByAddress(searchAddress, protocol, network)
	return err == nil
}

// GetPublicKeys that are stored on disk.
func (f FileStore) GetPublicKeys() ([]crypto.PublicKey, error) {
	rawKeys, err := f.getEncryptedKeys()
	if err != nil {
		return nil, err
	}

	publicKeys := []crypto.PublicKey{}

	for i := range rawKeys {
		pubKey, err := multikey.PublicKeyFromBytes(rawKeys[i].CurveType, rawKeys[i].PublicKeyBytes)
		if err != nil {
			continue
		}

		publicKeys = append(publicKeys, pubKey)
	}

	return publicKeys, nil
}

// GetAddresses list all the address this key store has
func (f FileStore) GetAddresses(protocol, network string) ([][]byte, error) {
	addresses := [][]byte{}

	keys, err := f.GetPublicKeys()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for _, pubKey := range keys {
		pubkeyAddress, err := address.FromPublicKey(pubKey, protocol, network)
		if err != nil {
			continue
		}

		addresses = append(addresses, pubkeyAddress)
	}

	return addresses, nil
}

func (f FileStore) getEncryptedKeyByAddress(searchAddress []byte, protocol, network string) (*keystore.EncryptedKey, error) {
	rawKeys, err := f.getEncryptedKeys()
	if err != nil {
		return nil, err
	}

	for i := range rawKeys {
		pubKey, err := multikey.PublicKeyFromBytes(rawKeys[i].CurveType, rawKeys[i].PublicKeyBytes)
		if err != nil {
			continue
		}

		pubkeyAddress, err := address.FromPublicKey(pubKey, protocol, network)
		if err != nil {
			continue
		}

		if bytes.Equal(pubkeyAddress, searchAddress) {
			return &rawKeys[i], nil
		}
	}

	return nil, errors.Errorf("not found")
}
