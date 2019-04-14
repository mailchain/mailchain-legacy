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
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/mailchain/mailchain/internal/pkg/crypto/cipher"
	"github.com/mailchain/mailchain/internal/pkg/encoding"
	"github.com/pkg/errors"
)

// Receive check ethereum transactions for mailchain messages
func (c APIClient) Receive(ctx context.Context, network string, address []byte) ([]cipher.EncryptedContent, error) {
	if !c.isNetworkSupported(network) {
		return nil, errors.Errorf("network not supported")
	}
	txResult, err := c.getTranscationsByAddress(network, address)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	res := []cipher.EncryptedContent{}
	for i := range txResult.Result { // TODO: paging
		x := txResult.Result[i]
		if !strings.HasPrefix(x.Input, "0x6d61696c636861696e") {
			continue
		}

		encryptedTransactionData, err := hexutil.Decode(x.Input)
		if err != nil {
			return nil, errors.WithMessage(err, "can not decode `data`")
		}
		if !bytes.HasPrefix(encryptedTransactionData, encoding.DataPrefix()) {
			return nil, errors.New("missing `mailchain` prefix")
		}

		res = append(res, encryptedTransactionData[len(encoding.DataPrefix()):])
	}
	return res, nil
}
