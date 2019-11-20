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
	"net/http"
	"sort"

	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings"
)

// GetProtocols returns a handler get spec
func GetProtocols(base *settings.Root) func(w http.ResponseWriter, r *http.Request) {
	res := GetProtocolsResponse{
		Protocols: []GetProtocolsProtocol{},
	}

	for _, protocol := range base.Protocols {
		if protocol.Disabled.Get() {
			continue
		}

		networks := []string{}

		for _, network := range protocol.Networks {
			if !network.Disabled() {
				networks = append(networks, network.Kind())
			}
		}

		sort.Strings(networks)
		resP := GetProtocolsProtocol{
			Name:     protocol.Kind,
			Networks: networks,
		}
		res.Protocols = append(res.Protocols, resP)
	}

	// Get swagger:route GET /protocols protocols GetProtocols
	//
	// Get protocols and the networks.
	//
	// Get all networks for each protocol that is enabled.
	// Responses:
	//   200: GetProtocolsResponse
	//   404: NotFoundError
	//   422: ValidationError
	return func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(res)
		w.Header().Set("Content-Type", "application/json")
	}
}

// GetProtocolsResponse Holds the response messages
//
// swagger:response GetProtocolsResponse
type GetProtocolsResponse struct {
	// in: body
	Protocols []GetProtocolsProtocol `json:"protocols"`
}

// GetProtocolsProtocol body
type GetProtocolsProtocol struct {
	// in: body
	Name string `json:"name"`
	// in: body
	Networks []string `json:"networks"`
}
