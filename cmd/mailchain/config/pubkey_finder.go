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
	"github.com/mailchain/mailchain/internal/pkg/mailbox"
	"github.com/pkg/errors"
	"github.com/spf13/viper" // nolint: depguard
)

func SetPubKeyFinder(vpr *viper.Viper, chain, network, pubkey string) error {
	viper.Set(fmt.Sprintf("chains.%s.networks.%s.pubkey-finder", chain, network), pubkey)
	if err := setClient(vpr, pubkey, network); err != nil {
		return err
	}
	fmt.Printf("%s used for looking up public key\n", pubkey)
	return nil
}

// GetPublicKeyFinders in configured state
func GetPublicKeyFinders() (map[string]mailbox.PubKeyFinder, error) {
	finders := make(map[string]mailbox.PubKeyFinder)
	for chain := range viper.GetStringMap("chains") {
		chFinders, err := getChainFinders(chain)
		if err != nil {
			return nil, err
		}
		if err := mergo.Merge(&finders, chFinders); err != nil {
			return nil, err
		}
	}
	return finders, nil
}

func getChainFinders(chain string) (map[string]mailbox.PubKeyFinder, error) {
	finders := make(map[string]mailbox.PubKeyFinder)
	for network := range viper.GetStringMap(fmt.Sprintf("chains.%s.networks", chain)) {
		finder, err := getFinder(chain, network)
		if err != nil {
			return nil, err
		}
		finders[fmt.Sprintf("%s.%s", chain, network)] = finder
	}

	return finders, nil
}

func getFinder(chain, network string) (mailbox.PubKeyFinder, error) {
	switch viper.GetString(fmt.Sprintf("chains.%s.networks.%s.pubkey-finder", chain, network)) {
	case names.Etherscan:
		return getEtherscanClient()
	default:
		return nil, errors.Errorf("unsupported pubkey finder")
	}
}
