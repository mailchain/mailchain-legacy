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

package nameservice

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

// Lookup methods for address and domain resolution.
type Lookup interface {
	ForwardLookup
	ReverseLookup
}

// NewLookupService creates a new lookup services.
func NewLookupService(baseURL string) Lookup {
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	return &LookupService{
		baseURL:    strings.TrimSuffix(baseURL, "/"),
		newRequest: http.NewRequest,
		doRequest:  client.Do,
	}
}

// LookupService is the default lookup service.
type LookupService struct {
	baseURL    string
	newRequest func(method string, url string, body io.Reader) (*http.Request, error)
	doRequest  func(req *http.Request) (*http.Response, error)
}

// ResolveName look up a domain name and return address on the related protocol and network pair.
func (s LookupService) ResolveName(ctx context.Context, protocol, network, domainName string) ([]byte, error) {
	req, err := s.newRequest("GET", fmt.Sprintf("%s/%s/%s/name?domain-name=%s", s.baseURL, protocol, network, domainName), nil)
	if err != nil {
		return nil, err
	}

	res, err := s.doRequest(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == http.StatusOK {
		var okRes struct {
			Address string `json:"address"`
		}

		if err := json.NewDecoder(res.Body).Decode(&okRes); err != nil {
			return nil, err
		}

		return common.FromHex(okRes.Address), nil
	}

	var errRes struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
	}

	if err := json.NewDecoder(res.Body).Decode(&errRes); err != nil {
		return nil, err
	}

	return nil, WrapError(errors.Errorf(errRes.Message))
}

// ResolveAddress looks up an address on the related protocol and network pair.
func (s LookupService) ResolveAddress(ctx context.Context, protocol, network string, address []byte) (string, error) {
	// {protocol}/{network}/address?address={address}
	req, err := s.newRequest("GET", fmt.Sprintf("%s/%s/%s/address?address=%s", s.baseURL, protocol, network, common.BytesToAddress(address).Hex()), nil)
	if err != nil {
		return "", err
	}

	res, err := s.doRequest(req)
	if err != nil {
		return "", err
	}

	if res.StatusCode == http.StatusOK {
		type response struct {
			Name string `json:"name"`
		}

		var okRes response
		if err := json.NewDecoder(res.Body).Decode(&okRes); err != nil {
			return "", err
		}

		return okRes.Name, nil
	}

	type response struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
	}

	var errRes response
	if err := json.NewDecoder(res.Body).Decode(&errRes); err != nil {
		return "", err
	}

	return "", WrapError(errors.Errorf(errRes.Message))
}
