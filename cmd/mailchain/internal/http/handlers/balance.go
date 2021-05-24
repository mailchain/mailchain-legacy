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
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mailchain/mailchain/cmd/internal/http/params"
	"github.com/mailchain/mailchain/errs"
	"github.com/mailchain/mailchain/internal/address"
	"github.com/mailchain/mailchain/internal/mailbox"
	"github.com/pkg/errors"
)

// GetBalance returns the balance of a user.
func GetBalance(balanceFinder map[string]mailbox.BalanceFinder) func(w http.ResponseWriter, r *http.Request) {
	// Get swagger:route GET /balance Balance GetBalance
	//
	// Get balance.
	//
	// Get the  balance of the user. The balance is used to send or receive messages.
	// Responses:
	//   200: GetBalanceResponse
	//   404: NotFoundError
	//   422: ValidationError
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		req, err := parseGetBalance(r)
		if err != nil {
			errs.JSONWriter(w, r, http.StatusUnprocessableEntity, errors.WithStack(err))
			return
		}
		fmt.Println(balanceFinder)
		balanceFinder, ok := balanceFinder[fmt.Sprintf("%s/%s", req.Protocol, req.Network)]
		if !ok {
			errs.JSONWriter(w, r, http.StatusUnprocessableEntity, errors.Errorf("balance not supported on \"%s/%s\"", req.Protocol, req.Network))
			return
		}

		if balanceFinder == nil {
			errs.JSONWriter(w, r, http.StatusUnprocessableEntity, errors.Errorf("no balance finder configured for \"%s/%s\"", req.Protocol, req.Network))
			return
		}

		balance, err := balanceFinder.GetBalance(ctx, req.Protocol, req.Network, req.addressBytes)
		if mailbox.IsNetworkNotSupportedError(err) {
			errs.JSONWriter(w, r, http.StatusNotAcceptable, errors.Errorf("network %q not supported", req.Network))
			return
		}

		if err != nil {
			errs.JSONWriter(w, r, http.StatusInternalServerError, errors.WithStack(err))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(GetBalanceResponseBody{
			Balance: balance,
			Unit:    "wei",
		})
	}
}

// GetBalanaceRequest pubic key from address request
// swagger:parameters GetBalance
type GetBalanaceRequest struct {
	// Address to to use when performing public key lookup.
	//
	// in: query
	// required: true
	// example: 0xde0b295669a9fd93d5f28d9ec85e40f4cb697bae
	// pattern: 0x[a-fA-F0-9]{40}
	Address      string `json:"address"`
	addressBytes []byte

	// Network to use when performing public key lookup.
	//
	// enum: mainnet,goerli,ropsten,rinkeby,local
	// in: query
	// required: true
	// example: goerli
	Network string `json:"network"`

	// Protocol to use when performing public key lookup.
	//
	// enum: ethereum
	// in: query
	// required: true
	// example: ethereum
	Protocol string `json:"protocol"`
}

// parseGetBalance get all the details for the get request
func parseGetBalance(r *http.Request) (*GetBalanaceRequest, error) {
	protocol, err := params.QueryRequireProtocol(r)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	network, err := params.QueryRequireNetwork(r)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	addr := mux.Vars(r)["address"]

	addressBytes, err := address.DecodeByProtocol(addr, protocol)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to decode address")
	}

	return &GetBalanaceRequest{
		Address:      addr,
		addressBytes: addressBytes,
		Network:      network,
		Protocol:     protocol,
	}, nil
}

// GetBalanceResponse public key from address response
//
// swagger:response GetBalanceResponse
type GetBalanceResponse struct {
	// in: body
	Body GetBalanceResponseBody
}

// GetBalanceResponseBody body response
//
// swagger:model GetBalanceBody
type GetBalanceResponseBody struct {
	// The public key encoded as per `public-key-encoding`
	//
	// Required: true
	// example: 0x79964e63752465973b6b3c610d8ac773fc7ce04f5d1ba599ba8768fb44cef525176f81d3c7603d5a2e466bc96da7b2443bef01b78059a98f45d5c440ca379463
	Balance uint64 `json:"balance"`

	// Encoding method used for encoding the `public-key`
	//
	// Required: true
	// example: hex/0x-prefix
	Unit string `json:"unit"`
}
