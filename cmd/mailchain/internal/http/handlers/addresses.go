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
	"net/http"

	"github.com/mailchain/mailchain/cmd/mailchain/internal/http/params"
	"github.com/mailchain/mailchain/errs"
	"github.com/mailchain/mailchain/internal/keystore"
	"github.com/pkg/errors"
)

// GetAddresses returns a handler get spec
func GetAddresses(ks keystore.Store) func(w http.ResponseWriter, r *http.Request) {
	// Get swagger:route GET /addresses Addresses GetAddresses
	//
	// Get addresses.
	//
	// Get all address that this user has access to. The addresses can be used to send or receive messages.
	// Responses:
	//   200: GetAddressesResponse
	//   404: NotFoundError
	//   422: ValidationError
	return func(w http.ResponseWriter, r *http.Request) {
		protocol, err := params.QueryRequireProtocol(r)
		if err != nil {
			errs.JSONWriter(w, http.StatusUnprocessableEntity, errors.WithStack(err))
			return
		}

		network, err := params.QueryRequireNetwork(r)
		if err != nil {
			errs.JSONWriter(w, http.StatusUnprocessableEntity, errors.WithStack(err))
			return
		}
		addresses := []string{}
		rawAddresses, err := ks.GetAddresses(protocol, network)
		if err != nil {
			errs.JSONWriter(w, http.StatusInternalServerError, errors.WithStack(err))
			return
		}
		for _, x := range rawAddresses {
			addresses = append(addresses, hex.EncodeToString(x))
		}

		_ = json.NewEncoder(w).Encode(GetAddressesResponse{Addresses: addresses})
		w.Header().Set("Content-Type", "application/json")
	}
}

// GetAddressesResponse Holds the response messages
//
// swagger:response GetAddressesResponse
type GetAddressesResponse struct {
	// in: body
	Addresses []string `json:"addresses"`
}

// GetAddressesRequest body
// swagger:parameters GetAddresses
type GetAddressesRequest struct {
	// Network to use when finding addresses.
	//
	// enum: mainnet,goerli,ropsten,rinkeby,local
	// in: query
	// required: true
	// example: goerli
	Network string `json:"network"`

	// Protocol to use when finding addresses.
	//
	// enum: ethereum, substrate
	// in: query
	// required: true
	// example: ethereum
	Protocol string `json:"protocol"`
}
