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
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/mailchain/mailchain/encoding"
	"github.com/mailchain/mailchain/internal/protocols/ethereum"
	"github.com/pkg/errors"
	"gopkg.in/resty.v1"
)

// NewAPIClient create new API client.
func NewAPIClient(apiKey string) (*APIClient, error) {
	return &APIClient{
		key: apiKey,
		networkConfigs: map[string]networkConfig{
			ethereum.Mainnet: {url: "https://api.etherscan.io/api"},
			ethereum.Ropsten: {url: "https://api-ropsten.etherscan.io/api"},
			ethereum.Kovan:   {url: "https://api-kovan.etherscan.io/api"},
			ethereum.Rinkeby: {url: "https://api-rinkeby.etherscan.io/api"},
			ethereum.Goerli:  {url: "https://api-goerli.etherscan.io/api"},
		},
	}, nil
}

// APIClient for talking to etherscan.
type APIClient struct {
	key            string
	networkConfigs map[string]networkConfig
}

type networkConfig struct {
	url string
}

// IsNetworkSupported checks if the network is supported by etherscan API.
func (c APIClient) isNetworkSupported(network string) bool {
	_, ok := c.networkConfigs[network]
	return ok
}

// GetTransactionByHash get transaction details from transaction hash via etherscan.
func (c APIClient) getTransactionByHash(network string, hash common.Hash) (*types.Transaction, error) {
	config, ok := c.networkConfigs[network]
	if !ok {
		return nil, errors.Errorf("network not supported")
	}

	txListResponse, err := resty.R().
		SetQueryParams(map[string]string{
			"module": "proxy",
			"action": "eth_getTransactionByHash",
			"txhash": hash.Hex(),
			"apikey": c.key,
		}).Get(config.url)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	res := &jsonrpcMessage{}
	if err := json.Unmarshal(txListResponse.Body(), res); err != nil {
		return nil, errors.WithStack(err)
	}

	if res.Error != nil {
		return nil, errors.Errorf(res.Error.Message)
	}

	if res.Result == nil {
		return nil, errors.Errorf("not found")
	}

	var ts *types.Transaction

	if err := json.Unmarshal(res.Result, &ts); err != nil {
		return nil, errors.WithMessage(err, string(txListResponse.Body()))
	}

	return ts, nil
}

// GetTransactionsByAddress get transactions from address via etherscan.
func (c APIClient) getTransactionsByAddress(network string, address []byte) (*txList, error) {
	config, ok := c.networkConfigs[network]
	if !ok {
		return nil, errors.Errorf("network not supported")
	}

	txListResponse, err := resty.R().
		SetQueryParams(map[string]string{
			"module":     "account",
			"action":     "txlist",
			"startblock": "0",
			"endblock":   "99999999",
			"sort":       "desc",
			"apikey":     c.key,
			"address":    encoding.EncodeHexZeroX(address),
		}).Get(config.url)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	txResult := &txList{}
	if err := json.Unmarshal(txListResponse.Body(), txResult); err != nil {
		return nil, errors.WithMessage(err, string(txListResponse.Body()))
	}

	return txResult, nil
}

// getBalanceByAddress get transactions from address via etherscan.
func (c APIClient) getBalanceByAddress(network string, address []byte) (*balanceResult, error) {
	config, ok := c.networkConfigs[network]
	if !ok {
		return nil, errors.Errorf("network not supported")
	}

	balanceResponse, err := resty.R().
		SetQueryParams(map[string]string{
			"module":  "account",
			"action":  "balance",
			"apikey":  c.key,
			"address": encoding.EncodeHexZeroX(address),
		}).Get(config.url)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	fmt.Println("%s", balanceResponse)

	balanceresult := &balanceResult{}
	if err := json.Unmarshal(balanceResponse.Body(), balanceresult); err != nil {
		return nil, errors.WithMessage(err, string(balanceResponse.Body()))
	}

	return balanceresult, nil
}
