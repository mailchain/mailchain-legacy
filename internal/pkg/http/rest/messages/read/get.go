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

package read

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mailchain/mailchain/internal/pkg/http/rest/errs"
	"github.com/mailchain/mailchain/internal/pkg/mail"
	"github.com/mailchain/mailchain/internal/pkg/stores"
	"github.com/pkg/errors"
)

// Get returns a handler get spec
func Get(store stores.Inbox) func(w http.ResponseWriter, r *http.Request) {
	// Get swagger:route GET /messages/{message_id}/read Messages GetRead
	//
	// Get message read status.
	//
	// Responses:
	//   200: GetReadResponse
	//   404: NotFoundError
	//   422: ValidationError
	return func(w http.ResponseWriter, r *http.Request) {
		messageID, err := mail.FromHexString(mux.Vars(r)["message_id"])
		if err != nil {
			errs.JSONWriter(w, http.StatusNotAcceptable, errors.WithMessage(err, "invalid `message_id`"))
			return
		}
		read, err := store.GetReadStatus(messageID)
		if stores.IsNotFoundError(err) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if err != nil {
			errs.JSONWriter(w, http.StatusInternalServerError, err)
			return
		}

		js, err := json.Marshal(GetBody{
			Read: read,
		})
		if err != nil {
			errs.JSONWriter(w, http.StatusInternalServerError, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(js)
	}
}

// swagger:response GetReadResponse
type GetReadResponse struct {
	// in: body
	Body GetBody
}

// swagger:model GetReadResponseBody
type GetBody struct {
	// Read
	//
	// Required: true
	// example: true
	Read bool `json:"read"`
}
