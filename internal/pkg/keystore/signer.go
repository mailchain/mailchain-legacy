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

package keystore

import (
	"github.com/mailchain/mailchain/internal/pkg/chains/ethereum"
	"github.com/mailchain/mailchain/internal/pkg/crypto/keys"
	"github.com/mailchain/mailchain/internal/pkg/encoding"
	"github.com/mailchain/mailchain/internal/pkg/mailbox"
	"github.com/pkg/errors"
)

type signerFunc func(pk keys.PrivateKey) (mailbox.Signer, error)

// Signer use the correct function to get the decrypter from private key
func Signer(chain string, pk keys.PrivateKey) (mailbox.Signer, error) {
	table := map[string]signerFunc{
		encoding.Ethereum: func(pk keys.PrivateKey) (mailbox.Signer, error) {
			return ethereum.NewSigner(pk), nil
		},
	}

	f, ok := table[chain]
	if !ok {
		return nil, errors.Errorf("unsupported signer type")
	}
	return f(pk)
}
