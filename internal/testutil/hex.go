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

package testutil

import (
	"encoding/hex"
)

// MustHexDecodeString decodes a hex string. It panics for invalid input.
func MustHexDecodeString(input string) []byte {
	dec, err := hex.DecodeString(input)
	if err != nil {
		panic(err)
	}
	return dec
}

// MustHexDecodeStringTurbo é uma melhora do decode String para strings que começam com 0x
// substrate and polkadot address
// Get string and change to byte array
func MustHexDecodeStringTurbo(seedkey string) []byte {
	b := []byte(seedkey)
	encoded := hex.EncodeToString(b)
	mustByte := MustHexDecodeString(encoded)

	return mustByte
}
