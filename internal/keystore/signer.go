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

package keystore

import (
	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/internal/mailbox/signer"
	"github.com/mailchain/mailchain/internal/protocols"
	"github.com/mailchain/mailchain/internal/protocols/algorand"
	"github.com/mailchain/mailchain/internal/protocols/ethereum"
	"github.com/mailchain/mailchain/internal/protocols/substrate"
	"github.com/pkg/errors"
)

// Signer use the correct function to get the decrypter from private key
func Signer(protocol string, pk crypto.PrivateKey) (signer.Signer, error) {
	switch protocol {
	case protocols.Algorand:
		return algorand.NewSigner(pk)
	case protocols.Ethereum:
		return ethereum.NewSigner(pk)
	case protocols.Substrate:
		return substrate.NewSigner(pk), nil
	default:
		return nil, errors.Errorf("unsupported signer type")
	}
}
