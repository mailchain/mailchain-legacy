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
	"github.com/mailchain/mailchain/cmd/mailchain/internal/config/names"
	"github.com/mailchain/mailchain/internal/mailbox"
	"github.com/pkg/errors"
	"github.com/spf13/viper" // nolint: depguard
)

//go:generate mockgen -source=sender.go -package=configtest -destination=./configtest/sender_mock.go

type SenderSetter interface {
	Set(chain, network, sender string) error
}

type Sender struct {
	viper        *viper.Viper
	clientGetter ClientsGetter
	clientSetter ClientsSetter
	mapMerge     func(dst interface{}, src interface{}, opts ...func(*mergo.Config)) error
}

func (s Sender) Set(chain, network, sender string) error {
	if err := s.clientSetter.SetClient(sender, network); err != nil {
		return err
	}
	s.viper.Set(fmt.Sprintf("chains.%s.networks.%s.sender", chain, network), sender)
	fmt.Printf("%q used for sending messages\n", sender)
	return nil
}

// GetSenders in configured state
func (s Sender) GetSenders() (map[string]mailbox.Sender, error) {
	senders := make(map[string]mailbox.Sender)
	for chain := range s.viper.GetStringMap("chains") {
		chSenders, err := s.getChainSenders(chain)
		if err != nil {
			return nil, err
		}
		if err := s.mapMerge(&senders, chSenders); err != nil {
			return nil, err
		}
	}
	return senders, nil
}

func (s Sender) getChainSenders(chain string) (map[string]mailbox.Sender, error) {
	senders := make(map[string]mailbox.Sender)
	for network := range s.viper.GetStringMap(fmt.Sprintf("chains.%s.networks", chain)) {
		sender, err := s.getSender(chain, network)
		if err != nil {
			return nil, err
		}
		senders[fmt.Sprintf("%s.%s", chain, network)] = sender
	}

	return senders, nil
}

func (s Sender) getSender(chain, network string) (mailbox.Sender, error) {
	switch s.viper.GetString(fmt.Sprintf("chains.%s.networks.%s.sender", chain, network)) {
	case names.EthereumRPC2:
		return s.clientGetter.GetEtherRPC2Client(network)
	default:
		return nil, errors.Errorf("unsupported sender")
	}
}
