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

package algod

import (
	"context"

	"github.com/algorand/go-algorand-sdk/client/v2/indexer"
	"github.com/algorand/go-algorand-sdk/types"
	"github.com/mailchain/mailchain/encoding"
	"github.com/mailchain/mailchain/internal/addressing"
	"github.com/mailchain/mailchain/stores"
	"github.com/pkg/errors"
)

// Receive check algorand transactions for mailchain messages.
func (c *Client) Receive(ctx context.Context, protocol, network string, address []byte) ([]stores.Transaction, error) {
	if !c.isNetworkSupported(network) {
		return nil, errors.New("network not supported")
	}

	algodClient, err := indexer.MakeClient(c.networkConfigs[network].url, c.algodToken)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to create algod client")
	}

	txResult, err := algodClient.SearchForTransactions().NotePrefix(encoding.DataPrefix()).AddressRole("receiver").Address(toAlgodAddress(address)).Do(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	res := []stores.Transaction{}

	for i := range txResult.Transactions {
		tx := txResult.Transactions[i]

		var reKeyAddress []byte
		if tx.RekeyTo != "" {
			reKeyAddress, _ = addressing.DecodeByProtocol(tx.RekeyTo, protocol)
		}

		res = append(res, stores.Transaction{
			EnvelopeData: tx.Note[len(encoding.DataPrefix()):],
			BlockNumber:  int64(tx.RoundTime),
			Hash:         []byte(tx.Id),
			RekeyAddress: reKeyAddress,
		})
	}

	return res, nil
}

func (c *Client) Kind() string {
	return "algod"
}

func toAlgodAddress(address []byte) types.Address {
	var out types.Address

	copy(out[:], address)

	return out
}
