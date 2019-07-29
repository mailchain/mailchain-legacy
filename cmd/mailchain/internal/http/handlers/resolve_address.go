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
	"encoding/json"
	"encoding/hex"
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
	// Get swagger:route GET /nameservice/address/{address}/resolve ResolveName NameService GetResolveAddress
	//
	// Get name from address.
	//
	// Get name.
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
			errs.JSONWriter(w, http.StatusUnprocessableEntity, errors.Errorf("no name servier resolver for chain.network configured"))
			return
		}

		name, err := resolver.ResolveAddress(ctx, protocol, network, address)
		if mailbox.IsNetworkNotSupportedError(err) {
			errs.JSONWriter(w, http.StatusNotAcceptable, errors.Errorf("%q not supported", protocol+"/"+network))
			return
		}
		if nameservice.IsInvalidAddressError(err) {
			errs.JSONWriter(w, http.StatusPreconditionFailed, err)
			return
		}
		if nameservice.IsNoResolverError(err) {
			errs.JSONWriter(w, http.StatusNotFound, err)
			return
		}
		if nameservice.IsNotFoundError(err) {
			errs.JSONWriter(w, http.StatusNotFound, err)
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
	// enum: mainnet,ropsten,rinkeby,local
	// in: path
	// required: true
	// example: ropsten
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
	Body GetPublicKeyResponseBody
}

// GetBody body response
//
// swagger:model GetResolveAddressResponseBody
type GetResolveAddressResponseBody struct {
	// The public key
	//
	// Required: true
	// example: mailchain.eth
	Name string `json:"name"`
}
