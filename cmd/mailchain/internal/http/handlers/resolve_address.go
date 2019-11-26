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

package handlers

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/http/params"
	"github.com/mailchain/mailchain/errs"
	"github.com/mailchain/mailchain/internal/mailbox"
	"github.com/mailchain/mailchain/nameservice"
	"github.com/pkg/errors"
)

// GetResolveAddress returns a handler get spec
func GetResolveAddress(resolvers map[string]nameservice.ReverseLookup) func(w http.ResponseWriter, r *http.Request) {
	// Get swagger:route GET /nameservice/address/{address}/resolve?network={network}&protocol={protocol} NameService GetResolveAddress
	//
	// Resolve Address Against Name Service
	//
	// Get name for supplied address. The name is typically a human-readable value that can be used in place of the address.
	// Resolve will query the protocol's name service to find the human-readable name for the supplied address.
	//
	// Responses:
	//   200: GetResolveAddressResponse
	//   404: NotFoundError
	//   422: ValidationError
	return func(w http.ResponseWriter, hr *http.Request) {
		ctx := hr.Context()
		protocol, network, address, err := parseGetResolveAddressRequest(hr)
		if err != nil {
			errs.JSONWriter(w, http.StatusUnprocessableEntity, errors.WithStack(err))
			return
		}
		resolver, ok := resolvers[fmt.Sprintf("%s/%s", protocol, network)]
		if !ok {
			errs.JSONWriter(w, http.StatusUnprocessableEntity, errors.Errorf("nameserver not supported on \"%s/%s\"", protocol, network))
			return
		}
		if resolver == nil {
			errs.JSONWriter(w, http.StatusUnprocessableEntity, errors.Errorf("no nameserver configured for \"%s/%s\"", protocol, network))
			return
		}

		name, err := resolver.ResolveAddress(ctx, protocol, network, address)
		if mailbox.IsNetworkNotSupportedError(err) {
			errs.JSONWriter(w, http.StatusNotAcceptable, errors.Errorf("%q not supported", protocol+"/"+network))
			return
		}
		if nameservice.IsRFC1035Error(err) {
			_ = json.NewEncoder(w).Encode(GetResolveAddressResponseBody{
				Status: nameservice.RFC1035StatusMap[err],
			})
			return
		}
		if err != nil {
			errs.JSONWriter(w, http.StatusInternalServerError, errors.WithStack(err))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(GetResolveAddressResponseBody{
			Name: name,
		})
	}
}

// GetResolveAddressRequest pubic key from address request
// swagger:parameters GetResolveAddress
type GetResolveAddressRequest struct {
	// name to query to get address for
	//
	// in: path
	// required: true
	// example: 0x4ad2b251246aafc2f3bdf3b690de3bf906622c51
	Address string `json:"address"`

	// Network for the name to resolve
	//
	// enum: mainnet,goerli,ropsten,rinkeby,local
	// in: path
	// required: true
	// example: goerli
	Network string `json:"network"`

	// Protocol for the name to resolve
	//
	// enum: ethereum
	// in: path
	// required: true
	// example: ethereum
	Protocol string `json:"protocol"`
}

// parseGetResolveAddressRequest get all the details for the get request
func parseGetResolveAddressRequest(r *http.Request) (protocol, network string, address []byte, err error) {
	protocol, err = params.QueryRequireProtocol(r)
	if err != nil {
		return "", "", nil, err
	}
	network, err = params.QueryRequireNetwork(r)
	if err != nil {
		return "", "", nil, err
	}

	address, err = hex.DecodeString(strings.TrimPrefix(strings.ToLower(mux.Vars(r)["address"]), "0x"))
	if err != nil {
		return "", "", nil, err
	}

	return protocol, network, address, nil
}

// GetResolveAddressResponse address of resolved name
//
// swagger:response GetResolveAddressResponse
type GetResolveAddressResponse struct {
	// in: body
	Body GetResolveAddressResponseBody
}

// GetResolveAddressResponseBody body response
//
// swagger:model GetResolveAddressResponseBody
type GetResolveAddressResponseBody struct {
	// The resolved name
	//
	// Required: true
	// example: mailchain.eth
	Name string `json:"name"`

	// The RFC1035 status code describing the outcome of the lookup
	//
	// + 0 - No Error
	// + 1 - Format Error
	// + 2 - Server Failure
	// + 3 - Non-Existent Domain
	// + 4 - Not Implemented
	// + 5 - Query Refused
	//
	// Required: false
	// example: 3
	Status int `json:"status"`
}
