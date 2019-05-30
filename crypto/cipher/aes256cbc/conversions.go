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

package aes256cbc

import (
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/secp256k1"
	"github.com/pkg/errors"
)

func asPublicECIES(pk crypto.PublicKey) (*ecies.PublicKey, error) {
	switch rpk := pk.(type) {
	case *secp256k1.PublicKey:
		return rpk.ECIES()
	case secp256k1.PublicKey:
		return rpk.ECIES()
	default:
		return nil, errors.Errorf("could not convert public key")
	}
}
func asPrivateECIES(pk crypto.PrivateKey) (*ecies.PrivateKey, error) {
	switch rpk := pk.(type) {
	case *secp256k1.PrivateKey:
		return rpk.ECIES(), nil
	case secp256k1.PrivateKey:
		return rpk.ECIES(), nil
	default:
		return nil, errors.Errorf("could not convert private key")
	}
}
