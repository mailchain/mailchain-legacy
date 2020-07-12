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

package mailchain

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mailchain/mailchain/encoding"
	"github.com/mailchain/mailchain/internal/address"
	"github.com/mailchain/mailchain/internal/mailbox"
	"github.com/pkg/errors"
	"gopkg.in/resty.v1"
)

// NewReceiver create new API client.
func NewReceiver(address string) (*Receiver, error) {
	return &Receiver{
		address: address,
	}, nil
}

// Receiver for talking to etherscan.
type Receiver struct {
	address string
}

// GetTransactionsByAddress get transactions from address via etherscan.
func (c Receiver) getTransactionsByAddress(protocol, network string, addr []byte) (*envelopeList, error) {
	encodedAddress, _, err := address.EncodeByProtocol(addr, protocol)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	envelopeListResponse, err := resty.R().
		SetQueryParams(map[string]string{
			"protocol": protocol,
			"network":  network,
			"address":  encodedAddress,
		}).
		Get(fmt.Sprintf("%s/to", strings.Trim(c.address, "/")))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	txResult := &envelopeList{}
	if err := json.Unmarshal(envelopeListResponse.Body(), txResult); err != nil {
		return nil, errors.WithMessage(err, string(envelopeListResponse.Body()))
	}

	return txResult, nil
}

// Receive check ethereum transactions for mailchain messages
func (c Receiver) Receive(ctx context.Context, protocol, network string, address []byte) ([]mailbox.Transaction, error) {
	txResult, err := c.getTransactionsByAddress(protocol, network, address)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	res := []mailbox.Transaction{}
	txHashes := map[string]bool{}

	for i := range txResult.Envelopes { //nolint TODO: paging
		x := txResult.Envelopes[i]

		_, ok := txHashes[x.Hash]
		if ok {
			continue
		}

		txHashes[x.Hash] = true

		encryptedTransactionData, err := encoding.DecodeHexZeroX(x.Data)
		if err != nil {
			continue // invalid data should move to next record
		}

		if !bytes.HasPrefix(encryptedTransactionData, encoding.DataPrefix()) {
			continue
		}

		res = append(res, mailbox.Transaction{
			Data: encryptedTransactionData[len(encoding.DataPrefix()):],
			Hash: []byte(x.Hash),
		})
	}

	return res, nil
}

type envelopeList struct {
	Envelopes []Envelope `json:"envelopes"`
}

type Envelope struct {
	From      string `json:"from"`
	To        string `json:"to"`
	Data      string `json:"data"`
	BlockHash string `json:"block-hash"`
	Hash      string `json:"hash"`

	Value    string `json:"value"`
	GasUsed  string `json:"gas-used"`
	GasPrice string `json:"gas-price"`
}
