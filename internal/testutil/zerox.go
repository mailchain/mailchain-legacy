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

import "github.com/mailchain/mailchain/internal/encoding"

// MustDecodeZeroX decodes a hex string. It panics for invalid input.
func MustDecodeZeroX(in string) []byte {
	dec, err := encoding.DecodeZeroX(in)
	if err != nil {
		panic(err)
	}
	return dec
}
