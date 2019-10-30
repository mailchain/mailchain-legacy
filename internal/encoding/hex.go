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

func EncodeZeroX(in []byte) (encoded, encoding string) {
	out := make([]byte, len(in)*2+2)
	copy(out, "0x")
	hex.Encode(out[2:], in)

	return string(out), TypeHex0XPrefix
}

func DecodeZeroX(in string) ([]byte, error) {
	if in == "" {
		return nil, errors.Errorf("empty hex string")
	}

	if !strings.HasPrefix(in, "0x") {
		return nil, errors.Errorf("missing \"0x\" prefix from hex string")
	}
	
	return hex.DecodeString(in[2:])
}
