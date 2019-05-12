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
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/mailchain/mailchain/internal/pkg/crypto/keys"
	"github.com/pkg/errors"
)

// PublicKey based on the secp256k1 curve
type PublicKey struct {
	ecdsa ecdsa.PublicKey
}

// Bytes returns the byte representation of the public key
func (pk PublicKey) Bytes() []byte {
	return crypto.CompressPubkey(&pk.ecdsa)
}

// Address returns the byte representation of the address
func (pk PublicKey) Address() []byte {
	return crypto.PubkeyToAddress(pk.ecdsa).Bytes()
}

// PublicKeyFromBytes create a public key from []byte
func PublicKeyFromBytes(pk []byte) (*PublicKey, error) {
	rpk, err := crypto.UnmarshalPubkey(pk)
	if err != nil {
		return nil, errors.WithMessage(err, "could not convert pk")
	}
	return &PublicKey{ecdsa: *rpk}, nil
}

// PublicKeyFromHex create a public key from hex
func PublicKeyFromHex(input string) (*PublicKey, error) {
	input = strings.TrimPrefix(input, "0x")
	keyBytes, err := hex.DecodeString(input)
	if err != nil {
		return nil, err
	}
	pk, err := crypto.DecompressPubkey(keyBytes)
	if err != nil {
		return nil, errors.WithMessage(err, "could not decompress pk")
	}
	return &PublicKey{ecdsa: *pk}, nil
}

// TODO: hang off object instead
func PublicKeyToECIES(pk keys.PublicKey) (*ecies.PublicKey, error) {
	rpk, err := crypto.DecompressPubkey(pk.Bytes())
	if err != nil {
		return nil, errors.WithMessage(err, "could not convert pk")
	}
	return ecies.ImportECDSAPublic(rpk), nil
}
