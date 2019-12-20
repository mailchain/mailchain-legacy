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

package ethereum

import (
	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/secp256k1"
	"github.com/pkg/errors"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/sha3"
)

// Address returns the byte representation of the address
func Address(pubKey crypto.PublicKey) ([]byte, error) {
	switch pubKey.(type) {
	case secp256k1.PublicKey, *secp256k1.PublicKey:
		ecdsa, err := ethcrypto.UnmarshalPubkey(append([]byte{byte(4)}, pubKey.Bytes()...))
		if err != nil {
			return nil, errors.WithMessage(err, "could not decompress key")
		}

		return Keccak256(ethcrypto.FromECDSAPub(ecdsa)[1:])[12:], nil
	default:
		return nil, errors.Errorf("invalid public key type: %T", pubKey)
	}
}

// Keccak256 calculates and returns the Keccak256 hash of the input data.
func Keccak256(data ...[]byte) []byte {
	d := sha3.NewLegacyKeccak256()

	for _, b := range data {
		_, _ = d.Write(b)
	}

	return d.Sum(nil)
}
