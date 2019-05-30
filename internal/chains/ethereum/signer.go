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
	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/secp256k1"
	"github.com/mailchain/mailchain/internal/mailbox"
	"github.com/pkg/errors"
)

type SignerOptions struct {
	Tx      *types.Transaction
	ChainID *big.Int
}

func NewSigner(privateKey crypto.PrivateKey) Signer {
	return Signer{privateKey: privateKey}
}

type Signer struct {
	privateKey crypto.PrivateKey
}

func (e Signer) Sign(opts mailbox.SignerOpts) (signedTransaction interface{}, err error) {
	if opts == nil {
		return nil, errors.New("opts must not be nil")
	}
	pk, ok := e.privateKey.(*secp256k1.PrivateKey)
	if !ok {
		return nil, errors.Errorf("invalid key type")
	}
	ecdsa, err := pk.ECDSA()
	if err != nil {
		return nil, err
	}
	switch opts := opts.(type) {
	case SignerOptions:
		// Depending on the presence of the chain ID, sign with EIP155 or homestead
		if opts.ChainID != nil {
			return types.SignTx(opts.Tx, types.NewEIP155Signer(opts.ChainID), ecdsa)
		}
		return types.SignTx(opts.Tx, types.HomesteadSigner{}, ecdsa)
	default:
		return nil, errors.New("invalid options for ethereum signing")
	}
}
