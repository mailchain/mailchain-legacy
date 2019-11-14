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

import (
	"encoding/hex"
	"strings"

	"github.com/pkg/errors"
)

// EncodeZeroX encodes src into "0x"+hex.Encode. As a convenience, it returns the encoding type used,
// but this value is always TypeHex0XPrefix.
// EncodeZeroX uses hexadecimal encoding prefixed with "0x".
func EncodeZeroX(src []byte) (encoded, encoding string) {
	out := make([]byte, len(src)*2+2)
	copy(out, "0x")
	hex.Encode(out[2:], src)

	return string(out), TypeHex0XPrefix
}

// DecodeZeroX returns the bytes represented by the hexadecimal string src.
//
// DecodeZeroX expects that src contains only hex characters and must contain a `0x` prefix.
// If the input is malformed, DecodeZeroX returns an error.
func DecodeZeroX(in string) ([]byte, error) {
	if in == "" {
		return nil, errors.Errorf("empty hex string")
	}

	if !strings.HasPrefix(in, "0x") {
		return nil, errors.Errorf("missing \"0x\" prefix from hex string")
	}

	return hex.DecodeString(in[2:])
}
