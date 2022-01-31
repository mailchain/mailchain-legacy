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
	"bytes"
	"strings"

	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/multikey"
	"github.com/mailchain/mailchain/encoding"
	"github.com/mailchain/mailchain/internal/addressing"
	"github.com/mailchain/mailchain/internal/keystore"
	"github.com/mailchain/mailchain/internal/protocols"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// HasAddress check for the presence of the address in the store.
func (f FileStore) HasAddress(searchAddress []byte, protocol, network string) bool {
	_, err := f.getEncryptedKeyByAddress(searchAddress, protocol, network)

	return err == nil
}

// GetPublicKeys that are stored on disk.
func (f FileStore) getPublicKeys(protocol, network string) ([]crypto.PublicKey, error) {
	rawKeys, err := f.getEncryptedKeys(protocol, network)
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

// GetAddresses list all addresses.
func (f FileStore) GetAddresses(protocol, network string) (map[string]map[string][][]byte, error) {
	protocol = strings.TrimSpace(protocol)
	network = strings.TrimSpace(network)
	addresses := map[string]map[string][][]byte{}

	if protocol == "" && network == "" {
		for _, protocol := range protocols.All() {
			addresses[protocol] = map[string][][]byte{}

			for _, network := range protocols.NetworkNames(protocol) {
				a, err := f.getProtocolNetworkAddresses(protocol, network)
				if err != nil {
					return nil, err
				}

				addresses[protocol][network] = a
			}
		}
	} else if protocol != "" && network == "" {
		addresses[protocol] = map[string][][]byte{}

		for _, network := range protocols.NetworkNames(protocol) {
			a, err := f.getProtocolNetworkAddresses(protocol, network)
			if err != nil {
				return nil, err
			}

			addresses[protocol][network] = a
		}
	} else if protocol != "" && network != "" {
		addresses[protocol] = map[string][][]byte{}
		a, err := f.getProtocolNetworkAddresses(protocol, network)
		if err != nil {
			return nil, err
		}

		addresses[protocol][network] = a
	} else if protocol == "" && network != "" {
		return nil, errors.New("protocol must be specified if network is supplied")
	}

	return addresses, nil
}

func (f FileStore) getProtocolNetworkAddresses(protocol, network string) ([][]byte, error) {
	addresses := [][]byte{}

	keys, err := f.getPublicKeys(protocol, network)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for _, pubKey := range keys {
		pubkeyAddress, err := addressing.FromPublicKey(pubKey, protocol, network)
		if err != nil {
			log.Logger.Warn().Str("component", "nacl-filestore").Str("func", "getProtocolNetworkAddresses").Str("protocol", protocol).Str("network", network).Str("public-key", encoding.EncodeHex(pubKey.Bytes())).Err(err).Msg("could not get address from public key")

			continue
		}

		addresses = append(addresses, pubkeyAddress)
	}

	return addresses, nil
}

func (f FileStore) getEncryptedKeyByAddress(searchAddress []byte, protocol, network string) (*keystore.EncryptedKey, error) {
	var out *keystore.EncryptedKey

	rawKeys, err := f.getEncryptedKeys(protocol, network)
	if err != nil {
		return nil, err
	}

	for i := range rawKeys {
		pubKey, err := multikey.PublicKeyFromBytes(rawKeys[i].CurveType, rawKeys[i].PublicKeyBytes)
		if err != nil {
			continue
		}

		pubkeyAddress, err := addressing.FromPublicKey(pubKey, protocol, network)
		if err != nil {
			continue
		}

		if bytes.Equal(pubkeyAddress, searchAddress) {
			out = &rawKeys[i]

			break
		}
	}

	if out == nil {
		return nil, errors.Errorf("not found")
	}

	return out, nil
}
