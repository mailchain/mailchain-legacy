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

package secp256k1

import (
	"crypto/ecdsa"
	"encoding/hex"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/mailchain/mailchain/internal/pkg/crypto/keys"
	"github.com/pkg/errors"
)

// PrivateKey based on the secp256k1 curve
type PrivateKey struct {
	ecdsa ecdsa.PrivateKey
}

// Bytes returns the byte representation of the private key
func (pk PrivateKey) Bytes() []byte {
	return crypto.FromECDSA(&pk.ecdsa)
}

// PublicKey return the public key that is derived from the private key
func (pk PrivateKey) PublicKey() keys.PublicKey {
	return PublicKey{ecdsa: pk.ecdsa.PublicKey}
}

// PrivateKeyFromECDSA get a private key from an ecdsa.PrivateKey
func PrivateKeyFromECDSA(pk ecdsa.PrivateKey) PrivateKey {
	return PrivateKey{ecdsa: pk}
}

// PrivateKeyFromBytes get a private key from []byte
func PrivateKeyFromBytes(pk []byte) (*PrivateKey, error) {
	rpk, err := crypto.ToECDSA(pk)
	if err != nil {
		return nil, errors.Errorf("could not convert private key")
	}
	return &PrivateKey{ecdsa: *rpk}, nil
}

// PrivateKeyFromHex get a private key from hex string
func PrivateKeyFromHex(hexkey string) (*PrivateKey, error) {
	b, err := hex.DecodeString(hexkey)
	if err != nil {
		return nil, errors.New("invalid hex string")
	}
	return PrivateKeyFromBytes(b)
}

// TODO: hang off key object instead
func PrivateKeyToECIES(pk keys.PrivateKey) (*ecies.PrivateKey, error) {
	rpk, err := crypto.ToECDSA(pk.Bytes())
	if err != nil {
		return nil, errors.Errorf("could not convert private key")
	}
	return ecies.ImportECDSA(rpk), nil
}

// TODO: hang off key object instead
func PrivateKeyToECDSA(pk keys.PrivateKey) (*ecdsa.PrivateKey, error) {
	rpk, err := crypto.ToECDSA(pk.Bytes())
	if err != nil {
		return nil, errors.Errorf("could not convert private key")
	}
	return rpk, nil
}
