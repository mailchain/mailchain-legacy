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

package multikey

import (
	"github.com/mailchain/mailchain/internal/pkg/crypto/keys"
	"github.com/mailchain/mailchain/internal/pkg/crypto/keys/secp256k1"
	"github.com/pkg/errors"
)

// PrivateKeyFromHex get private key from hex.
func PrivateKeyFromHex(hex, keyType string) (keys.PrivateKey, error) {
	table := map[string]privateKeyFromHex{
		SECP256K1: func(hex string) (keys.PrivateKey, error) {
			return secp256k1.PrivateKeyFromHex(hex)
		},
	}

	f, ok := table[keyType]
	if !ok {
		return nil, errors.Errorf("func for key type %v not registered", keyType)
	}
	return f(hex)
}

type privateKeyFromHex func(hex string) (keys.PrivateKey, error)
