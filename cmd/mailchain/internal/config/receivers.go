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

// nolint: dupl
package config

import (
	"fmt"

	"github.com/imdario/mergo"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/config/names"
	"github.com/mailchain/mailchain/internal/mailbox"
	"github.com/pkg/errors"
	"github.com/spf13/viper" // nolint: depguard
)

func SetReceiver(v *viper.Viper, chain, network, receiver string) error {
	if err := setClient(v, receiver, network); err != nil {
		return err
	}
	v.Set(fmt.Sprintf("chains.%s.networks.%s.receiver", chain, network), receiver)
	fmt.Printf("%s used for receiving messages\n", receiver)
	return nil
}

// GetReceivers in configured state
func GetReceivers(v *viper.Viper) (map[string]mailbox.Receiver, error) {
	receivers := make(map[string]mailbox.Receiver)
	for chain := range v.GetStringMap("chains") {
		chRcvrs, err := getChainReceivers(v, chain)
		if err != nil {
			return nil, err
		}
		if err := mergo.Merge(&receivers, chRcvrs); err != nil {
			return nil, err
		}
	}
	return receivers, nil
}

func getChainReceivers(v *viper.Viper, chain string) (map[string]mailbox.Receiver, error) {
	receivers := make(map[string]mailbox.Receiver)
	for network := range v.GetStringMap(fmt.Sprintf("chains.%s.networks", chain)) {
		receiver, err := getReceiver(v, chain, network)
		if err != nil {
			return nil, err
		}
		receivers[fmt.Sprintf("%s.%s", chain, network)] = receiver
	}

	return receivers, nil
}

func getReceiver(v *viper.Viper, chain, network string) (mailbox.Receiver, error) {
	switch v.GetString(fmt.Sprintf("chains.%s.networks.%s.receiver", chain, network)) {
	case names.Etherscan:
		return getEtherscanClient(v)
	case names.EtherscanNoAuth:
		return getEtherscanNoAuthClient()
	default:
		return nil, errors.Errorf("unsupported receiver")
	}
}
