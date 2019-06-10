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

type Receiver struct {
	viper        *viper.Viper
	clientGetter ClientsGetter
	clientSetter ClientsSetter
	mapMerge     func(dst interface{}, src interface{}, opts ...func(*mergo.Config)) error
}

func (r Receiver) Set(chain, network, receiver string) error {
	if err := r.clientSetter.SetClient(receiver, network); err != nil {
		return err
	}
	r.viper.Set(fmt.Sprintf("chains.%s.networks.%s.receiver", chain, network), receiver)
	fmt.Printf("%s used for receiving messages\n", receiver)
	return nil
}

// GetReceivers in configured state
func (r Receiver) GetReceivers() (map[string]mailbox.Receiver, error) {
	receivers := make(map[string]mailbox.Receiver)
	for chain := range r.viper.GetStringMap("chains") {
		chRcvrs, err := r.getChainReceivers(chain)
		if err != nil {
			return nil, err
		}
		if err := r.mapMerge(&receivers, chRcvrs); err != nil {
			return nil, err
		}
	}
	return receivers, nil
}

func (r Receiver) getChainReceivers(chain string) (map[string]mailbox.Receiver, error) {
	receivers := make(map[string]mailbox.Receiver)
	for network := range r.viper.GetStringMap(fmt.Sprintf("chains.%s.networks", chain)) {
		receiver, err := r.getReceiver(chain, network)
		if err != nil {
			return nil, err
		}
		receivers[fmt.Sprintf("%s.%s", chain, network)] = receiver
	}

	return receivers, nil
}

func (r Receiver) getReceiver(chain, network string) (mailbox.Receiver, error) {
	switch r.viper.GetString(fmt.Sprintf("chains.%s.networks.%s.receiver", chain, network)) {
	case names.Etherscan:
		return r.clientGetter.GetEtherscanClient()
	case names.EtherscanNoAuth:
		return r.clientGetter.GetEtherscanNoAuthClient()
	default:
		return nil, errors.Errorf("unsupported receiver")
	}
}
