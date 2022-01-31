// Copyright 2022 Mailchain Ltd
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

package blockscout

import (
	"bytes"
	"context"
	"strconv"

	"github.com/mailchain/mailchain/encoding"
	"github.com/mailchain/mailchain/stores"
	"github.com/pkg/errors"
)

// Receive check ethereum transactions for mailchain messages
func (c APIClient) Receive(ctx context.Context, protocol, network string, address []byte) ([]stores.Transaction, error) {
	if !c.isNetworkSupported(network) {
		return nil, errors.Errorf("network not supported")
	}

	txResult, err := c.getTransactionsByAddress(network, address)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	res := []stores.Transaction{}
	txHashes := map[string]bool{}

	for i := range txResult.Result {
		x := txResult.Result[i]

		_, ok := txHashes[x.Hash]
		if ok {
			continue
		}

		txHashes[x.Hash] = true

		encryptedTransactionData, err := encoding.DecodeHexZeroX(x.Input)
		if err != nil {
			continue // invalid data should move to next record
		}

		if !bytes.HasPrefix(encryptedTransactionData, encoding.DataPrefix()) {
			continue
		}

		blkNo, err := strconv.ParseInt(x.BlockNumber, 10, 32)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		res = append(res, stores.Transaction{
			EnvelopeData: encryptedTransactionData[len(encoding.DataPrefix()):],
			BlockNumber:  blkNo,
			Hash:         []byte(x.Hash),
		})
	}

	return res, nil
}

func (c APIClient) Kind() string {
	return "blockscout"
}
