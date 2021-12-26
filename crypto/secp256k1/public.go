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
	"crypto/sha256"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/mailchain/mailchain/crypto"
	"github.com/pkg/errors"
)

// PublicKey based on the secp256k1 curve
type PublicKey struct {
	ecdsa ecdsa.PublicKey
}

// Verify verifies whether sig is a valid signature of message.
func (pk PublicKey) Verify(message, sig []byte) bool {
	// VerifySignature requires the signature to be in
	// [ R || S ] format, so we remove the recid if present.
	if len(sig) == 65 {
		sig = sig[:64]
	}

	hash := sha256.Sum256(message)

	return ethcrypto.VerifySignature(ethcrypto.CompressPubkey(&pk.ecdsa), hash[:], sig)
}

// Bytes returns the byte representation of the public key
func (pk PublicKey) Bytes() []byte {
	return ethcrypto.CompressPubkey(&pk.ecdsa)
}

func (pk PublicKey) UncompressedBytes() []byte {
	return append(pk.ecdsa.X.Bytes(), pk.ecdsa.Y.Bytes()...)
}

// PublicKeyFromBytes create a public key from []byte
func PublicKeyFromBytes(keyBytes []byte) (crypto.PublicKey, error) {
	switch len(keyBytes) {
	case 65:
		rpk, err := ethcrypto.UnmarshalPubkey(keyBytes)
		if err != nil {
			return nil, errors.WithMessage(err, "could not convert pk")
		}

		return &PublicKey{ecdsa: *rpk}, nil
	case 64:
		rpk, err := ethcrypto.UnmarshalPubkey(append([]byte{byte(4)}, keyBytes...))
		if err != nil {
			return nil, errors.WithMessage(err, "could not convert pk")
		}

		return &PublicKey{ecdsa: *rpk}, nil
	case 33:
		pk, err := ethcrypto.DecompressPubkey(keyBytes)
		if err != nil {
			return nil, errors.WithMessage(err, "could not decompress pk")
		}

		return &PublicKey{ecdsa: *pk}, nil
	default:
		return nil, errors.Errorf("invalid key length %v", len(keyBytes))
	}
}

// ECIES returns an ECIES representation of the public key.
func (pk PublicKey) ECIES() (*ecies.PublicKey, error) {
	return ecies.ImportECDSAPublic(&pk.ecdsa), nil
}

func (pk PublicKey) ECDSA() *ecdsa.PublicKey {
	return &pk.ecdsa
}
