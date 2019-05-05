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

	"github.com/imdario/mergo"
	"github.com/mailchain/mailchain/cmd/mailchain/config/names"
	"github.com/mailchain/mailchain/internal/pkg/mailbox"
	"github.com/pkg/errors"
	"github.com/spf13/viper" // nolint: depguard
)

func SetSender(vpr *viper.Viper, chain, network, sender string) error {
	viper.Set(fmt.Sprintf("chains.%s.networks.%s.sender", chain, network), sender)
	if err := setClient(vpr, sender, network); err != nil {
		return err
	}
	fmt.Printf("%s used for sending messages\n", sender)
	return nil
}

// GetSenders in configured state
func GetSenders(vpr *viper.Viper) (map[string]mailbox.Sender, error) {
	senders := make(map[string]mailbox.Sender)
	for chain := range viper.GetStringMap("chains") {
		chSenders, err := getChainSenders(vpr, chain)
		if err != nil {
			return nil, err
		}
		if err := mergo.Merge(&senders, chSenders); err != nil {
			return nil, err
		}
	}
	return senders, nil
}

func getChainSenders(vpr *viper.Viper, chain string) (map[string]mailbox.Sender, error) {
	senders := make(map[string]mailbox.Sender)
	for network := range viper.GetStringMap(fmt.Sprintf("chains.%s.networks", chain)) {
		sender, err := getSender(vpr, chain, network)
		if err != nil {
			return nil, err
		}
		senders[fmt.Sprintf("%s.%s", chain, network)] = sender
	}

	return senders, nil
}

func getSender(vpr *viper.Viper, chain, network string) (mailbox.Sender, error) {
	switch viper.GetString(fmt.Sprintf("chains.%s.networks.%s.sender", chain, network)) {
	case names.EthereumRPC2:
		return getEtherRPC2Client(vpr, network)
	default:
		return nil, errors.Errorf("unsupported receiver")
	}
}
