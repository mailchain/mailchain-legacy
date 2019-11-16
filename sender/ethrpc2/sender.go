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

package ethrpc2

import (
	"context"
	"math/big"

	geth "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/mailchain/mailchain/internal/mailbox/signer"
	"github.com/mailchain/mailchain/internal/protocols/ethereum"
	"github.com/mailchain/mailchain/sender"
	"github.com/pkg/errors"
)

// Send transaction using the RPC2 client.
func (e EthRPC2) Send(ctx context.Context, network string, to, from, data []byte, txSigner signer.Signer, opts sender.SendOpts) error {
	chainID, err := e.client.NetworkID(ctx)
	if err != nil {
		return errors.WithMessage(err, "could not determine chain id")
	}
	gasPrice, err := e.client.SuggestGasPrice(ctx)
	if err != nil {
		return errors.WithMessage(err, "could not determine gas price")
	}

	value := big.NewInt(0)
	addrTo := common.BytesToAddress(to)
	addrFrom := common.BytesToAddress(from)
	gas, err := e.client.EstimateGas(ctx, geth.CallMsg{
		Data:     data,
		From:     addrFrom,
		GasPrice: gasPrice,
		To:       &addrTo,
		Value:    value,
	})
	if err != nil {
		return errors.WithMessage(err, "could not estimate gas")
	}
	nonce, err := e.client.NonceAt(ctx, addrFrom, nil) // from address
	if err != nil {
		return errors.WithStack(err)
	}

	rawSignedTx, err := txSigner.Sign(ethereum.SignerOptions{
		Tx:      types.NewTransaction(nonce, addrTo, value, gas, gasPrice, data),
		ChainID: chainID,
	})
	if err != nil {
		return errors.WithMessage(err, "could not sign transaction")
	}
	signedTx, ok := rawSignedTx.(*types.Transaction)
	if !ok {
		return errors.Errorf("sign did not return an ethereum transaction")
	}

	return e.client.SendTransaction(ctx, signedTx)
}
