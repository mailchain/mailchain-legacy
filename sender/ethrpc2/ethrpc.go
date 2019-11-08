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
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Options to use when sending an ethereum transaction.
type Options struct {
	Tx      *types.Transaction
	ChainID *big.Int
}

// New ethereum RPC2 sender is created that dials the address.
func New(address string) (*EthRPC2, error) {
	client, err := ethclient.Dial(address)
	if err != nil {
		return nil, err
	}
	return &EthRPC2{client: client}, nil
}

// EthRPC2 ethereum JSON-RPC2 client that is used to send transactions.
type EthRPC2 struct {
	client Client
}
