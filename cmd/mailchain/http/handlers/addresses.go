// Copyright 2022 Mailchain Ltd.
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
	"net/http"

	"github.com/mailchain/mailchain/cmd/internal/http/params"
	"github.com/mailchain/mailchain/encoding"
	"github.com/mailchain/mailchain/errs"
	"github.com/mailchain/mailchain/internal/addressing"
	"github.com/mailchain/mailchain/internal/keystore"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// GetAddresses returns a handler get spec.
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
		protocol := params.QueryOptionalProtocol(r)
		network := params.QueryOptionalNetwork(r)

		rawAddresses, err := ks.GetAddresses(protocol, network)
		if err != nil {
			errs.JSONWriter(w, r, http.StatusInternalServerError, errors.WithStack(err), log.Logger)

			return
		}

		flattenAddresses := keystore.FlattenAddressesMap(rawAddresses)
		addresses := []GetAddressesItem{}

		for _, x := range flattenAddresses {
			value, addressEncoding, err := addressing.EncodeByProtocol(x.Address, x.Protocol)
			if err != nil {
				errs.JSONWriter(w, r, http.StatusInternalServerError, errors.WithStack(err), log.Logger.With().Str("query-protocol", protocol).Str("address", encoding.EncodeHexZeroX(x.Address)).Logger())

				return
			}

			addresses = append(addresses, GetAddressesItem{
				Value:    value,
				Encoding: addressEncoding,
				Protocol: x.Protocol,
				Network:  x.Network,
			})
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
	Addresses []GetAddressesItem `json:"addresses"`
}

type GetAddressesItem struct {
	// Address value
	//
	// Required: true
	// example: 0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761
	Value string `json:"value"`
	// Encoding method used for encoding the `address`
	//
	// Required: true
	// example: hex/0x-prefix
	Encoding string `json:"encoding"`
	// Protocol `address` is available on
	//
	// Required: true
	// example: ethereum
	Protocol string `json:"protocol"`
	// Network `address` is available on
	//
	// Required: true
	// example: mainnet
	Network string `json:"network"`
}

// GetAddressesRequest bod
// swagger:parameters GetAddresses
type GetAddressesRequest struct {
	// Network to use when finding addresses.
	//
	// enum: mainnet,goerli,ropsten,rinkeby,local
	// in: query
	// required: false
	// example: goerli
	Network string `json:"network"`

	// Protocol to use when finding addresses.
	//
	// enum: ethereum, substrate, algorand
	// in: query
	// required: false
	// example: ethereum
	Protocol string `json:"protocol"`
}
