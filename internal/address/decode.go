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

package address

import (
	"github.com/mailchain/mailchain/internal/encoding"
	"github.com/mailchain/mailchain/internal/protocols"
	"github.com/mr-tron/base58"
	"github.com/pkg/errors"
)

// DecodeByProtocol returns the raw `[]byte` from the supplied address.
func DecodeByProtocol(in, protocol string) ([]byte, error) {
	switch protocol {
	case protocols.Ethereum:
		return encoding.DecodeZeroX(in)
	case protocols.Substrate:
		return base58.Decode(in)
	default:
		return nil, errors.Errorf("%q unsupported protocol", protocol)
	}
}
