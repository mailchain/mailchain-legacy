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

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/http/params"
	"github.com/mailchain/mailchain/errs"
	"github.com/mailchain/mailchain/internal/mailbox"
	"github.com/pkg/errors"
)

// GetPublicKey returns a handler get spec
func GetPublicKey(finders map[string]mailbox.PubKeyFinder) func(w http.ResponseWriter, r *http.Request) {
	// Get swagger:route GET /ethereum/{network}/address/{address}/public-key PublicKey Ethereum GetPublicKey
	//
	// Get public key from an address.
	//
	// Get the public key.
	//
	// Responses:
	//   200: GetPublicKeyResponse
	//   404: NotFoundError
	//   422: ValidationError
	return func(w http.ResponseWriter, hr *http.Request) {
		ctx := hr.Context()
		address, network, err := parseGetPublicKey(hr)
		if err != nil {
			errs.JSONWriter(w, http.StatusUnprocessableEntity, errors.WithStack(err))
			return
		}
		finder, ok := finders[fmt.Sprintf("ethereum.%s", network)]
		if !ok {
			errs.JSONWriter(w, http.StatusUnprocessableEntity, errors.Errorf("no public key finder for chain.network configured"))
			return
		}

		publicKey, err := finder.PublicKeyFromAddress(ctx, network, address)
		if mailbox.IsNetworkNotSupportedError(err) {
			errs.JSONWriter(w, http.StatusNotAcceptable, errors.Errorf("network %q not supported", network))
			return
		}
		if err != nil {
			errs.JSONWriter(w, http.StatusInternalServerError, errors.WithStack(err))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(GetPublicKeyResponseBody{
			PublicKey: hexutil.Encode(publicKey),
		})
	}
}

// GetPublicKey pubic key from address request
// swagger:parameters GetPublicKey
type GetPublicKeyRequest struct {
	// address to query to get public key for
	//
	// in: path
	// required: true
	// example: 0xde0b295669a9fd93d5f28d9ec85e40f4cb697bae
	// pattern: 0x[a-fA-F0-9]{40}
	Address string `json:"address"`

	// Network for the message to send
	//
	// enum: mainnet,ropsten,rinkeby,local
	// in: path
	// required: true
	// example: ropsten
	Network string `json:"network"`
}

// parseGetPublicKey get all the details for the get request
func parseGetPublicKey(r *http.Request) (address []byte, network string, err error) {
	addr, err := params.PathAddress(r)
	if err != nil {
		return nil, "", errors.WithStack(err)
	}

	return addr, params.PathNetwork(r), nil
}

// GetPublicKeyResponse public key from address response
//
// swagger:response GetPublicKeyResponse
type GetPublicKeyResponse struct {
	// in: body
	Body GetPublicKeyResponseBody
}

// GetBody body response
//
// swagger:model GetPublicKeyResponseBody
type GetPublicKeyResponseBody struct {
	// The public key
	//
	// Required: true
	// nolint: lll
	// example: 0x79964e63752465973b6b3c610d8ac773fc7ce04f5d1ba599ba8768fb44cef525176f81d3c7603d5a2e466bc96da7b2443bef01b78059a98f45d5c440ca379463
	PublicKey string `json:"public_key"`
}
