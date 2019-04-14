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

package ethrpc

import (
	"context"
	"math/big"

	eth "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/mailchain/mailchain/internal/pkg/chains/ethereum"
	"github.com/mailchain/mailchain/internal/pkg/mailbox"
	"github.com/pkg/errors"
)

type Options struct {
	Tx      *types.Transaction
	ChainID *big.Int
}

func New(address string) (*EthRPC2, error) {
	client, err := ethclient.Dial(address)
	if err != nil {
		return nil, err
	}
	return &EthRPC2{Client: *client}, nil
}

type EthRPC2 struct {
	ethclient.Client
}

func (e EthRPC2) Send(ctx context.Context, to, from, data []byte, signer mailbox.Signer, opts mailbox.SenderOpts) error {
	chainID, err := e.NetworkID(ctx)
	if err != nil {
		return errors.WithMessage(err, "could not determine chain id")
	}
	gasPrice, err := e.SuggestGasPrice(ctx)
	if err != nil {
		return errors.WithMessage(err, "could not determine gas price")
	}

	value := big.NewInt(0)
	addrTo := common.BytesToAddress(to)
	addrFrom := common.BytesToAddress(from)
	gas, err := e.EstimateGas(ctx, eth.CallMsg{
		Data:     data,
		From:     addrFrom,
		GasPrice: gasPrice,
		To:       &addrTo,
		Value:    value,
	})
	if err != nil {
		return errors.WithMessage(err, "could not estimate gas")
	}
	nonce, err := e.NonceAt(ctx, addrFrom, nil) // from address
	if err != nil {
		return errors.WithStack(err)
	}

	rawSignedTx, err := signer.Sign(ethereum.SignerOptions{
		Tx:      types.NewTransaction(nonce, addrTo, value, gas, gasPrice, data),
		ChainID: chainID,
	})
	if err != nil {
		return errors.WithMessage(err, "could not sign transaction")
	}
	signedTx, ok := rawSignedTx.(*types.Transaction)
	if !ok {
		return errors.WithMessage(err, "could not cast transaction")
	}

	return e.SendTransaction(ctx, signedTx)
}
