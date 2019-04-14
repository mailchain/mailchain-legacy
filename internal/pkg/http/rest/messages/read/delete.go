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
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mailchain/mailchain/internal/pkg/http/rest/errs"
	"github.com/mailchain/mailchain/internal/pkg/mail"
	"github.com/mailchain/mailchain/internal/pkg/stores"
	"github.com/pkg/errors"
)

// DeleteRequest open api documentation
// swagger:parameters DeleteRead
type DeleteRequest struct {
	// Unique id of the message
	//
	// in: path
	// required: true
	MessageID string `json:"message_id"`
}

// Delete returns a handler
func Delete(store stores.Inbox) func(w http.ResponseWriter, r *http.Request) {
	// Delete swagger:route Delete /messages/{message_id}/read Messages DeleteRead
	//
	// Mark message as unread
	//
	// Responses:
	//   200: StatusOK
	//   404: NotFoundError
	//   422: ValidationError
	errHandler := errs.JSONHandler
	return func(w http.ResponseWriter, r *http.Request) {
		messageID, err := mail.FromHexString(mux.Vars(r)["message_id"])
		if err != nil {
			errHandler(w, http.StatusNotAcceptable, errors.WithMessage(err, "invalid `message_id`"))
			return
		}

		if err := store.DeleteMessageRead(messageID); err != nil {
			errHandler(w, http.StatusUnprocessableEntity, errors.WithStack(err))
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
