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

package ens

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/mailchain/mailchain/internal/nameservice"
	ens "github.com/wealdtech/go-ens"
)

func NewLookupService(clientURL string) (nameservice.Lookup, error) {
	client, err := ethclient.Dial(clientURL)
	if err != nil {
		return nil, err
	}
	return &LookupService{
		client: client,
	}, nil
}

type LookupService struct {
	client *ethclient.Client
}

func (s LookupService) ResolveName(ctx context.Context, protocol, network, domainName string) ([]byte, error) {
	address, err := ens.Resolve(s.client, domainName)
	if err != nil {
		return nil, wrapError(err)
	}
	return address.Bytes(), nil
}

func (s LookupService) ResolveAddress(ctx context.Context, protocol, network string, address []byte) (string, error) {
	ethAddress := common.BytesToAddress(address)
	reverse, err := ens.ReverseResolve(s.client, &ethAddress)
	if err != nil {
		return "", wrapError(err)
	}
	return reverse, nil
}
