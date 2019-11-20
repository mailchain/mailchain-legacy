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

package handlers //nolint: lll

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mailchain/mailchain/cmd/mailchain/internal/http/params"
	"github.com/mailchain/mailchain/errs"
	"github.com/mailchain/mailchain/internal/address"
	"github.com/mailchain/mailchain/internal/mailbox"
	"github.com/mailchain/mailchain/internal/pubkey"
	"github.com/pkg/errors"
)

// GetPublicKey returns a handler get spec
func GetPublicKey(finders map[string]mailbox.PubKeyFinder) func(w http.ResponseWriter, r *http.Request) {
	// Get swagger:route GET /public-key PublicKey GetPublicKey
	//
	// Public key from address.
	//
	// This method will get the public key to use when encrypting messages and envelopes.
	// Protocols and networks have different methods for retrieving or calculating a public key from an address.
	//
	// Responses:
	//   200: GetPublicKeyResponse
	//   404: NotFoundError
	//   422: ValidationError
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		req, err := parseGetPublicKey(r)
		if err != nil {
			errs.JSONWriter(w, http.StatusUnprocessableEntity, errors.WithStack(err))
			return
		}
		finder, ok := finders[fmt.Sprintf("%s/%s", req.Protocol, req.Network)]
		if !ok {
			errs.JSONWriter(w, http.StatusUnprocessableEntity, errors.Errorf("public key finder not supported on \"%s/%s\"", req.Protocol, req.Network))
			return
		}
		if finder == nil {
			errs.JSONWriter(w, http.StatusUnprocessableEntity, errors.Errorf("no public key finder configured for \"%s/%s\"", req.Protocol, req.Network))
			return
		}

		publicKey, err := finder.PublicKeyFromAddress(ctx, req.Protocol, req.Network, req.addressBytes)
		if mailbox.IsNetworkNotSupportedError(err) {
			errs.JSONWriter(w, http.StatusNotAcceptable, errors.Errorf("network %q not supported", req.Network))
			return
		}
		if err != nil {
			errs.JSONWriter(w, http.StatusInternalServerError, errors.WithStack(err))
			return
		}
		encodedKey, encodingType, err := pubkey.EncodeByProtocol(publicKey, req.Protocol)
		if err != nil {
			errs.JSONWriter(w, http.StatusInternalServerError, errors.WithStack(err))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(GetPublicKeyResponseBody{
			PublicKey:         encodedKey,
			PublicKeyEncoding: encodingType,
		})
	}
}

// GetPublicKeyRequest pubic key from address request
// swagger:parameters GetPublicKey
type GetPublicKeyRequest struct {
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

// parseGetPublicKey get all the details for the get request
func parseGetPublicKey(r *http.Request) (*GetPublicKeyRequest, error) {
	protocol, err := params.QueryRequireProtocol(r)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	network, err := params.QueryRequireNetwork(r)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	addr, err := params.QueryRequireAddress(r)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	addressBytes, err := address.DecodeByProtocol(addr, protocol)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to decode public key")
	}

	return &GetPublicKeyRequest{
		Address:      addr,
		addressBytes: addressBytes,
		Network:      network,
		Protocol:     protocol,
	}, nil
}

// GetPublicKeyResponse public key from address response
//
// swagger:response GetPublicKeyResponse
type GetPublicKeyResponse struct {
	// in: body
	Body GetPublicKeyResponseBody
}

// GetPublicKeyResponseBody body response
//
// swagger:model GetPublicKeyResponseBody
type GetPublicKeyResponseBody struct {
	// The public key encoded as per `public_key_encoding`
	//
	// Required: true
	// example: 0x79964e63752465973b6b3c610d8ac773fc7ce04f5d1ba599ba8768fb44cef525176f81d3c7603d5a2e466bc96da7b2443bef01b78059a98f45d5c440ca379463
	PublicKey string `json:"public_key"`

	// Encoding method used for encoding the `public_key`
	//
	// Required: true
	// example: hex/0x-prefix
	PublicKeyEncoding string `json:"public_key_encoding"`
}
