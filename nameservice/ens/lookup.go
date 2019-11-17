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
	"github.com/mailchain/mailchain/nameservice"
	ens "github.com/wealdtech/go-ens"
)

// NewLookupService creates a new ethereum name service (ENS) lookup service.
func NewLookupService(clientURL string) (nameservice.Lookup, error) {
	client, err := ethclient.Dial(clientURL)
	if err != nil {
		return nil, err
	}
	return &LookupService{
		client: client,
	}, nil
}

// LookupService for ethereum name service (ENS).
type LookupService struct {
	client *ethclient.Client
}

// ResolveName against the ethereum name service (ENS).
func (s LookupService) ResolveName(ctx context.Context, protocol, network, domainName string) ([]byte, error) {
	address, err := ens.Resolve(s.client, domainName)
	if err != nil {
		return nil, nameservice.WrapError(err)
	}
	return address.Bytes(), nil
}

// ResolveAddress against the ethereum name service (ENS).
func (s LookupService) ResolveAddress(ctx context.Context, protocol, network string, address []byte) (string, error) {
	ethAddress := common.BytesToAddress(address)
	reverse, err := ens.ReverseResolve(s.client, &ethAddress)
	if err != nil {
		return "", nameservice.WrapError(err)
	}
	return reverse, nil
}
