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
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/secp256k1"
	"github.com/mailchain/mailchain/internal/mailbox/signer"
	"github.com/pkg/errors"
)

// SignerOptions options that can be set when signing an ethereum transaction.
type SignerOptions struct {
	Tx      *types.Transaction
	ChainID *big.Int
}

// NewSigner returns a new ethereum signer that can be used to sign transactions.
func NewSigner(privateKey crypto.PrivateKey) (*Signer, error) {
	if _, err := validatePrivateKeyType(privateKey); err != nil {
		return nil, errors.WithStack(err)
	}

	return &Signer{privateKey: privateKey}, nil
}

// Signer for ethereum.
type Signer struct {
	privateKey crypto.PrivateKey
}

// Sign an ethereum transaction with the private key.
func (e Signer) Sign(opts signer.Options) (signedTransaction interface{}, err error) {
	if opts == nil {
		return nil, errors.New("opts must not be nil")
	}

	ecdsaPrivKey, err := validatePrivateKeyType(e.privateKey)
	if err != nil {
		return nil, err
	}

	switch opts := opts.(type) {
	case SignerOptions:
		// Depending on the presence of the chain ID, sign with EIP155 or homestead
		if opts.ChainID != nil {
			return types.SignTx(opts.Tx, types.NewEIP155Signer(opts.ChainID), ecdsaPrivKey)
		}

		return types.SignTx(opts.Tx, types.HomesteadSigner{}, ecdsaPrivKey)
	default:
		return nil, errors.New("invalid options for ethereum signing")
	}
}

func validatePrivateKeyType(pk crypto.PrivateKey) (*ecdsa.PrivateKey, error) {
	switch pk := pk.(type) {
	case secp256k1.PrivateKey:
		return pk.ECDSA()
	case *secp256k1.PrivateKey:
		return pk.ECDSA()
	default:
		return nil, errors.New("invalid key type")
	}
}
