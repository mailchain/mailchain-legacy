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
	"github.com/mailchain/mailchain/cmd/mailchain/config/names"
	"github.com/mailchain/mailchain/cmd/mailchain/prompts"
	"github.com/mailchain/mailchain/internal/pkg/mailbox"
	"github.com/pkg/errors"
	"github.com/spf13/viper" // nolint: depguard
)

func SetReceiver(network string) error {
	receiver, skipped, err := prompts.SelectItemSkipable(
		"Receiver",
		[]string{names.Etherscan},
		viper.GetString(fmt.Sprintf("chains.ethereum.networks.%s.receiver", network)) != "")
	if err != nil || skipped {
		return err
	}
	viper.Set(fmt.Sprintf("chains.ethereum.networks.%s.receiver", network), receiver)
	if err := setClient(receiver, network); err != nil {
		return err
	}
	fmt.Printf("%s used for receiving messages\n", receiver)
	return nil
}

// GetReceivers in configured state
func GetReceivers() (map[string]mailbox.Receiver, error) {
	receivers := make(map[string]mailbox.Receiver)
	for chain := range viper.GetStringMap("chains") {
		chRcvrs, err := getChainReceivers(chain)
		if err != nil {
			return nil, err
		}
		if err := mergo.Merge(&receivers, chRcvrs); err != nil {
			return nil, err
		}
	}
	return receivers, nil
}

func getChainReceivers(chain string) (map[string]mailbox.Receiver, error) {
	receivers := make(map[string]mailbox.Receiver)
	for network := range viper.GetStringMap(fmt.Sprintf("chains.%s.networks", chain)) {
		receiver, err := getReceiver(chain, network)
		if err != nil {
			return nil, err
		}
		receivers[fmt.Sprintf("%s.%s", chain, network)] = receiver
	}

	return receivers, nil
}

func getReceiver(chain, network string) (mailbox.Receiver, error) {
	switch viper.GetString(fmt.Sprintf("chains.%s.networks.%s.receiver", chain, network)) {
	case names.Etherscan:
		return getEtherscanClient()
	default:
		return nil, errors.Errorf("unsupported receiver")
	}
}
