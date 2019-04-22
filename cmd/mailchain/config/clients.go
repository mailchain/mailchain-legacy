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

	"github.com/mailchain/mailchain/cmd/mailchain/config/names"
	"github.com/mailchain/mailchain/internal/pkg/clients/etherscan"
	"github.com/mailchain/mailchain/internal/pkg/clients/ethrpc"
	"github.com/mailchain/mailchain/internal/pkg/cmd/prompts"
	"github.com/pkg/errors"
	"github.com/spf13/viper" // nolint: depguard
)

func getEtherRPC2Client(network string) (*ethrpc.EthRPC2, error) {
	address := viper.GetString(fmt.Sprintf("clients.ethereum-rpc2.%s.address", network))
	if address == "" {
		return nil, errors.Errorf("`clients.ethereum-rpc2.%s.address` must not be empty", network)
	}
	return ethrpc.New(address)
}

func getEtherscanClient() (*etherscan.APIClient, error) {
	apiKey := viper.GetString("clients.etherscan.api-key")
	if apiKey == "" {
		return nil, errors.Errorf("`clients.etherscan.api-key` must not be empty")
	}
	return etherscan.NewAPIClient(apiKey)
}

func setClient(client, network string) error {
	switch client {
	case names.EthereumRPC2:
		return setEthRPC(network)
	case names.Etherscan:
		return setEtherscan()
	default:
		return errors.Errorf("unsupported client type")
	}
}

func setEthRPC(network string) error {
	client := names.EthereumRPC2
	if viper.GetString(fmt.Sprintf("clients.%s.%s.address", client, network)) != "" {
		fmt.Printf("%s already configured\n", client)
		return nil
	}
	address, err := prompts.RequiredInput("Address")
	if err != nil {
		return err
	}
	viper.Set(fmt.Sprintf("clients.%s.%s.address", client, network), address)
	return nil
}

func setEtherscan() error {
	client := names.Etherscan
	if viper.GetString(fmt.Sprintf("clients.%s.api-key", client)) != "" {
		fmt.Printf("%s already configured\n", client)
		return nil
	}
	apiKey, err := prompts.RequiredInput("Api Key")
	if err != nil {
		return err
	}
	viper.Set(fmt.Sprintf("clients.%s.api-key", client), apiKey)
	fmt.Printf("%s configured\n", client)

	return nil
}
