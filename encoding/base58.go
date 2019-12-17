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

package encoding

import "github.com/mr-tron/base58"

// DecodeBase58 returns the bytes represented by the base58 string src.
//
// DecodeBase58 expects that src contains only base58 characters.
// If the input is malformed, DecodeBase58 returns an error.
func DecodeBase58(src string) ([]byte, error) {
	return base58.Decode(src)
}

// EncodeBase58 returns the string represented by the base58 byte src.
//
// EncodeBase58 expects that src contains only base58 byte.
// If the input is malformed, EncodeBase58 returns an error.
func EncodeBase58(src []byte) (string, error) {
	return base58.Encode(src), nil
}
