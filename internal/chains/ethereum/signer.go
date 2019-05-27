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

package ethereum

import (
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/mailchain/mailchain/internal/crypto/keys"
	"github.com/mailchain/mailchain/internal/crypto/keys/secp256k1"
	"github.com/mailchain/mailchain/internal/mailbox"
	"github.com/pkg/errors"
)

type SignerOptions struct {
	Tx      *types.Transaction
	ChainID *big.Int
}

func NewSigner(privateKey keys.PrivateKey) Signer {
	return Signer{privateKey: privateKey}
}

type Signer struct {
	privateKey keys.PrivateKey
}

func (e Signer) Sign(opts mailbox.SignerOpts) (signedTransaction interface{}, err error) {
	if opts == nil {
		return nil, errors.New("opts must not be nil")
	}

	pk, err := secp256k1.PrivateKeyToECDSA(e.privateKey)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	switch opts := opts.(type) {
	case SignerOptions:
		// Depending on the presence of the chain ID, sign with EIP155 or homestead
		if opts.ChainID != nil {
			return types.SignTx(opts.Tx, types.NewEIP155Signer(opts.ChainID), pk)
		}
		return types.SignTx(opts.Tx, types.HomesteadSigner{}, pk)
	default:
		return nil, errors.New("invalid options for ethereum signing")
	}
}
