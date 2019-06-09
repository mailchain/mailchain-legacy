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

func SetPubKeyFinder(v *viper.Viper, chain, network, pubkeyFinder string) error {
	if err := setClient(v, pubkeyFinder, network); err != nil {
		return err
	}
	v.Set(fmt.Sprintf("chains.%s.networks.%s.pubkey-finder", chain, network), pubkeyFinder)
	fmt.Printf("%s used for looking up public key\n", pubkeyFinder)
	return nil
}

// GetPublicKeyFinders in configured state
func GetPublicKeyFinders(v *viper.Viper) (map[string]mailbox.PubKeyFinder, error) {
	finders := make(map[string]mailbox.PubKeyFinder)
	for chain := range v.GetStringMap("chains") {
		chFinders, err := getChainFinders(v, chain)
		if err != nil {
			return nil, err
		}
		if err := mergo.Merge(&finders, chFinders); err != nil {
			return nil, err
		}
	}
	return finders, nil
}

func getChainFinders(v *viper.Viper, chain string) (map[string]mailbox.PubKeyFinder, error) {
	finders := make(map[string]mailbox.PubKeyFinder)
	for network := range v.GetStringMap(fmt.Sprintf("chains.%s.networks", chain)) {
		finder, err := getFinder(v, chain, network)
		if err != nil {
			return nil, err
		}
		finders[fmt.Sprintf("%s.%s", chain, network)] = finder
	}

	return finders, nil
}

func getFinder(v *viper.Viper, chain, network string) (mailbox.PubKeyFinder, error) {
	switch v.GetString(fmt.Sprintf("chains.%s.networks.%s.pubkey-finder", chain, network)) {
	case names.Etherscan:
		return getEtherscanClient(v)
	case names.EtherscanNoAuth:
		return getEtherscanNoAuthClient()
	default:
		return nil, errors.Errorf("unsupported pubkey finder")
	}
}
