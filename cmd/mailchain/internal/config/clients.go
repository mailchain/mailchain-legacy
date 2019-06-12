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

	"github.com/mailchain/mailchain"
	"github.com/mailchain/mailchain/internal/clients/etherscan"
	"github.com/mailchain/mailchain/internal/clients/ethrpc"
	"github.com/pkg/errors"
	"github.com/spf13/viper" // nolint: depguard
)

//go:generate mockgen -source=clients.go -package=configtest -destination=./configtest/clients_mock.go
type ClientsSetter interface {
	SetClient(client, network string) error
}

type ClientsGetter interface {
	GetEtherRPC2Client(network string) (*ethrpc.EthRPC2, error)
	GetEtherscanClient() (*etherscan.APIClient, error)
	GetEtherscanNoAuthClient() (*etherscan.APIClient, error)
}

type Clients struct {
	viper         *viper.Viper
	requiredInput func(label string) (string, error)
}

func (c Clients) GetEtherRPC2Client(network string) (*ethrpc.EthRPC2, error) {
	address := c.viper.GetString(fmt.Sprintf("clients.ethereum-rpc2.%s.address", network))
	if address == "" {
		return nil, errors.Errorf("`clients.ethereum-rpc2.%s.address` must not be empty", network)
	}
	return ethrpc.New(address)
}

func (c Clients) GetEtherscanClient() (*etherscan.APIClient, error) {
	apiKey := c.viper.GetString("clients.etherscan.api-key")
	if apiKey == "" {
		return nil, errors.Errorf("`clients.etherscan.api-key` must not be empty")
	}
	return etherscan.NewAPIClient(apiKey)
}

func (c Clients) GetEtherscanNoAuthClient() (*etherscan.APIClient, error) {
	return etherscan.NewAPIClient("")
}

func (c Clients) SetClient(client, network string) error {
	switch client {
	case mailchain.ClientEthereumRPC2:
		return c.setEthRPC(network)
	case mailchain.ClientEtherscan:
		return c.setEtherscan()
	case mailchain.ClientEtherscanNoAuth:
		return nil
	default:
		return errors.Errorf("unsupported client type")
	}
}

func (c Clients) setEthRPC(network string) error {
	client := mailchain.ClientEthereumRPC2
	if c.viper.GetString(fmt.Sprintf("clients.%s.%s.address", client, network)) != "" {
		fmt.Printf("%s already configured\n", client)
		return nil
	}
	address, err := c.requiredInput("Address")
	if err != nil {
		return err
	}
	c.viper.Set(fmt.Sprintf("clients.%s.%s.address", client, network), address)
	return nil
}

func (c Clients) setEtherscan() error {
	client := mailchain.ClientEtherscan
	if c.viper.GetString(fmt.Sprintf("clients.%s.api-key", client)) != "" {
		fmt.Printf("%s already configured\n", client)
		return nil
	}
	apiKey, err := c.requiredInput("Api Key")
	if err != nil {
		return err
	}
	c.viper.Set(fmt.Sprintf("clients.%s.api-key", client), apiKey)
	fmt.Printf("%s configured\n", client)

	return nil
}
