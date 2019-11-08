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

package substrate

import (
	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/ed25519"
	"github.com/minio/blake2b-simd"
	"github.com/pkg/errors"
)

// SS58AddressFormat creates the address `[]byte` from substrate network and a public key.
func SS58AddressFormat(network string, pubKey crypto.PublicKey) ([]byte, error) {
	if err := validPublicKeyType(pubKey); err != nil {
		return nil, errors.WithStack(err)
	}

	prefixedKey, err := prefixWithNetwork(network, pubKey)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	hash := blake2b.Sum512(addSS58Prefix(prefixedKey))

	// take first 2 bytes of hash since public key
	return append(prefixedKey, hash[:2]...), nil
}

func validPublicKeyType(pubKey crypto.PublicKey) error {
	switch pubKey.(type) {
	case ed25519.PublicKey, *ed25519.PublicKey:
		return nil
	default:
		return errors.Errorf("invalid public key type: %T", pubKey)
	}
}

func addSS58Prefix(pubKey []byte) []byte {
	prefix := []byte("SS58PRE")
	return append(prefix, pubKey...)
}

func prefixWithNetwork(network string, publicKey crypto.PublicKey) ([]byte, error) {
	// https://github.com/paritytech/substrate/wiki/External-Address-Format-(SS58)#address-type defines different prefixes by network
	switch network {
	case EdgewareTestnet:
		// 42 = 0x2a
		return append([]byte{0x2a}, publicKey.Bytes()...), nil
	default:
		return nil, errors.Errorf("unknown address prefix for %q", network)
	}
}
