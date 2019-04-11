// Copyright (c) 2019 Finobo
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

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
func ParseAddress(input string, chain string, network string) (*Address, error) {
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

func tryAddChainNetwork(input string, chain string, network string) (string, error) {
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
