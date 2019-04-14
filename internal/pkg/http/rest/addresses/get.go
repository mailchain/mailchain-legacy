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

package addresses

import (
	"encoding/hex"
	"encoding/json"
	"net/http"

	"github.com/mailchain/mailchain/internal/pkg/http/rest/errs"

	"github.com/mailchain/mailchain/internal/pkg/keystore"
	"github.com/pkg/errors"
)

// Get returns a handler get spec
func Get(ks keystore.Store) func(w http.ResponseWriter, r *http.Request) {
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
		addresses := []string{}
		rawAddresses, err := ks.GetAddresses()
		if err != nil {
			errs.JSONWriter(w, http.StatusInternalServerError, errors.WithStack(err))
			return
		}
		for _, x := range rawAddresses {
			addresses = append(addresses, hex.EncodeToString(x))
		}

		js, err := json.Marshal(GetResponse{Addresses: addresses})
		if err != nil {
			errs.JSONWriter(w, http.StatusInternalServerError, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(js)
	}

}

// GetResponse Holds the response messages
//
// swagger:response GetAddressesResponse
type GetResponse struct {
	// in: body
	Addresses []string `json:"addresses"`
}
