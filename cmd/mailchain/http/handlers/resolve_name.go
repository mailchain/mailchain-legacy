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
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/mailchain/mailchain/cmd/internal/http/params"
	"github.com/mailchain/mailchain/errs"
	"github.com/mailchain/mailchain/internal/addressing"
	"github.com/mailchain/mailchain/internal/mailbox"
	"github.com/mailchain/mailchain/nameservice"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// GetResolveName returns a handler get spec
func GetResolveName(resolvers map[string]nameservice.ForwardLookup) func(w http.ResponseWriter, r *http.Request) {
	// Get swagger:route GET /nameservice/name/{domain-name}/resolve?network={network}&protocol={protocol} NameService GetResolveName
	//
	// Resolve Name Against Name Service
	//
	// Get address for supplied name. The name is typically a human-readable value that can be used in place of the address.
	// Resolve will query the protocol's name service to find the address for supplied human-readable name.
	//
	// Responses:
	//   200: GetResolveNameResponse
	//   404: NotFoundError
	//   422: ValidationError
	return func(w http.ResponseWriter, r *http.Request) {
		protocol, network, domainName, err := parseGetResolveNameRequest(r)
		if err != nil {
			errs.JSONWriter(w, r, http.StatusUnprocessableEntity, errors.WithStack(err), log.Logger)
			return
		}

		resolver, ok := resolvers[fmt.Sprintf("%s/%s", protocol, network)]
		if !ok {
			errs.JSONWriter(w, r, http.StatusUnprocessableEntity, errors.Errorf("nameserver not supported on \"%s/%s\"", protocol, network), log.Logger)
			return
		}

		if resolver == nil {
			errs.JSONWriter(w, r, http.StatusUnprocessableEntity, errors.Errorf("no nameserver configured for \"%s/%s\"", protocol, network), log.Logger)
			return
		}

		resolvedAddress, err := resolver.ResolveName(r.Context(), protocol, network, domainName)
		if mailbox.IsNetworkNotSupportedError(err) {
			errs.JSONWriter(w, r, http.StatusNotAcceptable, errors.Errorf("%q not supported", protocol+"/"+network), log.Logger)
			return
		}

		if nameservice.ErrorToRFC1035Status(err) > 0 {
			_ = json.NewEncoder(w).Encode(GetResolveNameResponseBody{Status: nameservice.ErrorToRFC1035Status(err)})
			return
		}

		if err != nil {
			errs.JSONWriter(w, r, http.StatusInternalServerError, errors.WithStack(err), log.Logger)
			return
		}

		encAddress, _, err := addressing.EncodeByProtocol(resolvedAddress, protocol)
		if err != nil {
			errs.JSONWriter(w, r, http.StatusInternalServerError, errors.WithMessage(err, "failed to encode address"), log.Logger)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(GetResolveNameResponseBody{
			Address: encAddress,
		})
	}
}

// GetResolveNameRequest pubic key from address request
// swagger:parameters GetResolveName
type GetResolveNameRequest struct {
	// name to query to get address for
	//
	// in: path
	// required: true
	// example: mailchain.eth
	Name string `json:"domain-name"`

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

// parseGetResolveNameRequest get all the details for the get request
func parseGetResolveNameRequest(r *http.Request) (protocol, network, domain string, err error) {
	protocol, err = params.QueryRequireProtocol(r)
	if err != nil {
		return "", "", "", err
	}
	network, err = params.QueryRequireNetwork(r)
	if err != nil {
		return "", "", "", err
	}
	domain = strings.ToLower(mux.Vars(r)["domain-name"])

	return protocol, network, domain, nil
}

// GetResolveNameResponse address of resolved name
//
// swagger:response GetResolveNameResponse
type GetResolveNameResponse struct {
	// in: body
	Body GetResolveNameResponseBody
}

// GetResolveNameResponseBody body response
//
// swagger:model GetResolveNameResponseBody
type GetResolveNameResponseBody struct {
	// The resolved address
	//
	// Required: true
	// example: 0x4ad2b251246aafc2f3bdf3b690de3bf906622c51
	Address string `json:"address"`

	// The rFC1035 status code describing the outcome of the lookup
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
