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
	"bytes"
	"context"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/mailchain/mailchain/crypto/cipher"
	"github.com/mailchain/mailchain/internal/encoding"
	"github.com/pkg/errors"
)

// Receive check ethereum transactions for mailchain messages
func (c APIClient) Receive(ctx context.Context, network string, address []byte) ([]cipher.EncryptedContent, error) {
	if !c.isNetworkSupported(network) {
		return nil, errors.Errorf("network not supported")
	}
	txResult, err := c.getTransactionsByAddress(network, address)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	res := []cipher.EncryptedContent{}
	txHashes := map[string]bool{}

	for i := range txResult.Result { //nolint TODO: paging
		x := txResult.Result[i]

		_, ok := txHashes[x.Hash]
		if ok {
			continue
		}
		txHashes[x.Hash] = true

		encryptedTransactionData, err := hexutil.Decode(x.Input)
		if err != nil {
			continue // invalid data should move to next record
		}

		if !bytes.HasPrefix(encryptedTransactionData, encoding.DataPrefix()) {
			continue
		}

		res = append(res, encryptedTransactionData[len(encoding.DataPrefix()):])
	}
	return res, nil
}
