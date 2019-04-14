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

package config

import (
	"fmt"

	"github.com/mailchain/mailchain/internal/pkg/clients/etherscan"
	"github.com/mailchain/mailchain/internal/pkg/clients/ethrpc"
	"github.com/pkg/errors"
	"github.com/spf13/viper" // nolint: depguard
)

func etherRPC2Client(network string) (*ethrpc.EthRPC2, error) {
	address := viper.GetString(fmt.Sprintf("clients.ethereum-rpc2.%s.address", network))
	if address == "" {
		return nil, errors.Errorf("`clients.ethereum-rpc2.%s.address` must not be empty", network)
	}
	return ethrpc.New(address)
}

func etherscanClient() (*etherscan.APIClient, error) {
	apiKey := viper.GetString("clients.etherscan.api-key")
	if apiKey == "" {
		return nil, errors.Errorf("`clients.etherscan.api-key` must not be empty")
	}
	return etherscan.NewAPIClient(apiKey)
}
