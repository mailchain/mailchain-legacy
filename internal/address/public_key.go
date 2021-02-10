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
	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/internal/protocols"
	"github.com/mailchain/mailchain/internal/protocols/ethereum"
	"github.com/mailchain/mailchain/internal/protocols/substrate"
	"github.com/pkg/errors"
)

// FromPublicKey creates an address from public key.
func FromPublicKey(pubKey crypto.PublicKey, protocol, network string) (address []byte, err error) {
	switch protocol {
	case protocols.Ethereum:
		return ethereum.Address(pubKey)
	case protocols.Substrate:
		return substrate.SS58AddressFormat(network, pubKey)
	default:
		return nil, errors.Errorf("%q unsupported protocol", protocol)
	}
}
