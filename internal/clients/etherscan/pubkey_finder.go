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

package etherscan

import (
	"context"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/secp256k1"
	"github.com/mailchain/mailchain/encoding"
	"github.com/mailchain/mailchain/internal/protocols/ethereum"
	"github.com/pkg/errors"
)

// PublicKeyFromAddress get public key from the recipient address, this will only work if the recipient has previously sent a message.
func (c APIClient) PublicKeyFromAddress(ctx context.Context, protocol, network string, address []byte) (crypto.PublicKey, error) {
	if !c.isNetworkSupported(network) {
		return nil, errors.Errorf("network not supported")
	}

	txResult, err := c.getTransactionsByAddress(network, address)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	hash, err := getFromResultHash(encoding.EncodeHexZeroX(address), txResult)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	tx, err := c.getTransactionByHash(network, hash)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	v, r, s := tx.RawSignatureValues()

	keyBytes, err := ethereum.GetPublicKeyFromTransaction(
		r, s, v,
		tx.To().Bytes(),
		tx.Data(),
		tx.Nonce(),
		tx.GasPrice(),
		tx.Gas(),
		tx.Value())
	if err != nil {
		return nil, errors.WithMessage(err, "could not get public key from raw hash")
	}

	return secp256k1.PublicKeyFromBytes(keyBytes)
}

func getFromResultHash(address string, txResult *txList) (common.Hash, error) {
	if len(txResult.Result) == 0 {
		return common.Hash{}, errors.Errorf("No transactions found for address: %v", address)
	}

	for i := range txResult.Result {
		x := txResult.Result[i]
		if strings.EqualFold(x.From, address) {
			return common.HexToHash(x.Hash), nil
		}
	}

	return common.Hash{}, errors.Errorf("No transactions from address found")
}
