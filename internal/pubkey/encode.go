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

package pubkey

import (
	"github.com/mailchain/mailchain/encoding"
	"github.com/mailchain/mailchain/internal/protocols"
	"github.com/pkg/errors"
)

// EncodeByProtocol selects the correct encoding method, then encodes the public key with it.
func EncodeByProtocol(in []byte, protocol string) (encoded, encodingType string, err error) {
	switch protocol {
	case protocols.Algorand:
		encodingType = encoding.KindBase32
		encoded = encoding.EncodeBase32(in)
	case protocols.Ethereum:
		encodingType = encoding.KindHex0XPrefix
		encoded = encoding.EncodeHexZeroX(in)
	case protocols.Substrate:
		encodingType = encoding.KindHex0XPrefix
		encoded = encoding.EncodeHexZeroX(in)
	default:
		err = errors.Errorf("%q unsupported protocol", protocol)
	}

	return encoded, encodingType, err
}
