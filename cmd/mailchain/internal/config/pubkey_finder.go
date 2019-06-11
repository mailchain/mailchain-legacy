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

//go:generate mockgen -source=pubkey_finder.go -package=configtest -destination=./configtest/pubkey_finder_mock.go

type PubKeyFinderSetter interface {
	Set(chain, network, pubkeyFinder string) error 
}

type PubKeyFinder struct {
	viper        *viper.Viper
	clientGetter ClientsGetter
	clientSetter ClientsSetter
	mapMerge     func(dst interface{}, src interface{}, opts ...func(*mergo.Config)) error
}

func (p PubKeyFinder) Set(chain, network, pubkeyFinder string) error {
	if err := p.clientSetter.SetClient(pubkeyFinder, network); err != nil {
		return err
	}
	p.viper.Set(fmt.Sprintf("chains.%s.networks.%s.pubkey-finder", chain, network), pubkeyFinder)
	fmt.Printf("%s used for looking up public key\n", pubkeyFinder)
	return nil
}

// GetFinders in configured state
func (p PubKeyFinder) GetFinders() (map[string]mailbox.PubKeyFinder, error) {
	finders := make(map[string]mailbox.PubKeyFinder)
	for chain := range p.viper.GetStringMap("chains") {
		chFinders, err := p.getChainFinders(chain)
		if err != nil {
			return nil, err
		}
		if err := p.mapMerge(&finders, chFinders); err != nil {
			return nil, err
		}
	}
	return finders, nil
}

func (p PubKeyFinder) getChainFinders(chain string) (map[string]mailbox.PubKeyFinder, error) {
	finders := make(map[string]mailbox.PubKeyFinder)
	for network := range p.viper.GetStringMap(fmt.Sprintf("chains.%s.networks", chain)) {
		finder, err := p.getFinder(chain, network)
		if err != nil {
			return nil, err
		}
		finders[fmt.Sprintf("%s.%s", chain, network)] = finder
	}

	return finders, nil
}

func (p PubKeyFinder) getFinder(chain, network string) (mailbox.PubKeyFinder, error) {
	switch p.viper.GetString(fmt.Sprintf("chains.%s.networks.%s.pubkey-finder", chain, network)) {
	case names.Etherscan:
		return p.clientGetter.GetEtherscanClient()
	case names.EtherscanNoAuth:
		return p.clientGetter.GetEtherscanNoAuthClient()
	default:
		return nil, errors.Errorf("unsupported pubkey finder")
	}
}
