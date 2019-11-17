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

package mail

import (
	"fmt"
	nm "net/mail"
	"strings"

	"github.com/pkg/errors"
)

// Address represents a single mail address.
// An address such as "Charlotte <0x92d8f10248c6a3953cc3692a894655ad05d61efb@ropsten.ethereum.mailchain>" is represented
// as Address{Name: "Charlotte", Address: "0x92d8f10248c6a3953cc3692a894655ad05d61efb@ropsten.ethereum.mailchain"}.
type Address struct {
	DisplayName  string // Proper name; may be empty.
	FullAddress  string // 0x92d8f10248c6a3953cc3692a894655ad05d61efb@chain.network.mailchain
	ChainAddress string // 0x92d8f10248c6a3953cc3692a894655ad05d61efb
}

func (a *Address) String() string {
	addr := nm.Address{Address: a.FullAddress, Name: a.DisplayName}
	return addr.String()
}

// ParseAddress parses a single RFC 5322 address, e.g. "Charlotte <0x92d8f10248c6a3953cc3692a894655ad05d61efb@ropsten.ethereum.mailchain>"
// then apply it to the chain address.
func ParseAddress(input, chain, network string) (*Address, error) {
	input = strings.Trim(input, " ")
	if !strings.Contains(input, "@") {
		i, err := tryAddChainNetwork(input, chain, network)
		if err != nil {
			return nil, errors.WithMessage(err, "can not add missing @ in addr-spec")
		}
		input = i
	}

	addr, err := nm.ParseAddress(input)
	if err != nil {
		return nil, errors.WithMessage(err, "could not parse address")
	}

	return fromAddress(addr)
}

func tryAddChainNetwork(input, chain, network string) (string, error) {
	if strings.Contains(input, "@") {
		return input, nil
	}
	if strings.Contains(input, "<") || strings.Contains(input, ">") {
		return "", errors.New("can not add network and chain if display name is used")
	}
	if chain == "" || network == "" {
		return "", errors.New("both network and chain must be set")
	}

	return fmt.Sprintf("%s@%s.%s", input, network, chain), nil
}

func fromAddress(address *nm.Address) (*Address, error) {
	if address == nil {
		return nil, errors.Errorf("can not convert nil address")
	}

	parts := strings.Split(address.Address, "@")
	if len(parts) != 2 {
		return nil, errors.Errorf("missing @ in address")
	}

	return &Address{
		ChainAddress: parts[0],
		DisplayName:  address.Name,
		FullAddress:  address.Address,
	}, nil
}
