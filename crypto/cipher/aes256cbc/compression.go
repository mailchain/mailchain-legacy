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

package aes256cbc

import (
	"crypto/elliptic"

	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/pkg/errors"
)

const (
	pubKeyBytesLenCompressed   = 33
	pubKeyBytesLenUncompressed = 65
	compressedKeyPrefix        = 4
)

// compress a 65 byte uncompressed public key
func compress(publicKey []byte) ([]byte, error) {
	if len(publicKey) == pubKeyBytesLenUncompressed-1 && publicKey[0] != compressedKeyPrefix {
		publicKey = append([]byte{compressedKeyPrefix}, publicKey...)
	}
	if len(publicKey) != pubKeyBytesLenUncompressed {
		return nil, errors.Errorf("length of uncompressed public key is invalid")
	}
	x, y := elliptic.Unmarshal(ecies.DefaultCurve, publicKey)

	return secp256k1.CompressPubkey(x, y), nil
}

// decompress a 33 byte compressed public key
func decompress(publicKey []byte) []byte {
	x, y := secp256k1.DecompressPubkey(publicKey)
	return elliptic.Marshal(ecies.DefaultCurve, x, y)
}
