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

package addressing

import (
	"github.com/mailchain/mailchain/encoding"
	"github.com/mailchain/mailchain/internal/protocols"
	"github.com/pkg/errors"
)

// EncodeByProtocol takes an address as `[]byte` then selects the relevant encoding method to encode it as string.
func EncodeByProtocol(in []byte, protocol string) (encoded, encodingKind string, err error) {
	switch protocol {
	case protocols.Algorand:
		encodingKind = encoding.KindBase32
		encoded = encoding.EncodeBase32(in)
	case protocols.Ethereum:
		encodingKind = encoding.KindHex0XPrefix
		encoded = encoding.EncodeHexZeroX(in)
	case protocols.Substrate:
		encodingKind = encoding.KindBase58
		encoded = encoding.EncodeBase58(in)
	default:
		err = errors.Errorf("%q unsupported protocol", protocol)
	}

	return encoded, encodingKind, err
}
