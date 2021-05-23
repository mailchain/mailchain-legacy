// Copyright 2021 Finobo
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
	"github.com/mailchain/mailchain/internal/protocols/algorand"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// New algod client is created that stores config for each supported network.
func New(algodToken string) (*Client, error) {
	return &Client{
		networkConfigs: map[string]networkConfig{
			algorand.Mainnet: {url: "https://api.algoexplorer.io/idx2"},
			algorand.Testnet: {url: "https://api.testnet.algoexplorer.io/idx2"},
			algorand.Betanet: {url: "https://api.betanet.algoexplorer.io/idx2"},
		},
		algodToken: algodToken,
		logger:     log.With().Str("component", "receiver").Str("client", "algod").Logger(),
	}, nil
}

// Client ethereum JSON-RPC2 client that is used to send transactions.
type Client struct {
	networkConfigs map[string]networkConfig
	algodToken     string
	logger         zerolog.Logger
}

type networkConfig struct {
	url string
}

// IsNetworkSupported checks if the network is supported by etherscan API.
func (c *Client) isNetworkSupported(network string) bool {
	_, ok := c.networkConfigs[network]

	return ok
}
