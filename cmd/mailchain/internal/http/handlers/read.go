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

	"github.com/gorilla/mux"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/http/params"
	"github.com/mailchain/mailchain/errs"
	"github.com/mailchain/mailchain/internal/mail"
	"github.com/mailchain/mailchain/stores"
	"github.com/pkg/errors"
)

// DeleteRead returns a handler
func DeleteRead(store stores.State) func(w http.ResponseWriter, r *http.Request) {
	// DeleteRequest open api documentation
	// swagger:parameters DeleteRead
	type deleteRequest struct {
		// Unique id of the message
		//
		// in: path
		// required: true
		MessageID string `json:"message_id"`
	}
	// Delete swagger:route Delete /messages/{message_id}/read Messages DeleteRead
	//
	// Mark message as unread
	//
	// Responses:
	//   200: StatusOK
	//   404: NotFoundError
	//   422: ValidationError
	return func(w http.ResponseWriter, r *http.Request) {
		doRead(store.DeleteMessageRead, w, r)
	}
}

// GetRead returns a handler get spec
func GetRead(store stores.State) func(w http.ResponseWriter, r *http.Request) {
	// swagger:response GetReadResponse
	type getReadResponse struct {
		// in: body
		Body getBody
	}
	// Get swagger:route GET /messages/{message_id}/read Messages GetRead
	//
	// Message read status.
	//
	// Messages can be either read or unread.
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
		if err := json.NewEncoder(w).Encode(getBody{Read: read}); err != nil {
			errs.JSONWriter(w, http.StatusInternalServerError, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
	}
}

// PutRead returns a handler put spec
func PutRead(store stores.State) func(w http.ResponseWriter, r *http.Request) {
	// putRequest open api documentation
	// swagger:parameters PutRead
	type putRequest struct { //nolint: unused
		// Unique id of the message
		//
		// in: path
		// required: true
		MessageID string `json:"message_id"`
	}
	// Put swagger:route PUT /messages/{message_id}/read Messages PutRead
	//
	// PutRead.
	//
	// Mark message as read.
	// Responses:
	//   200: StatusOK
	//   404: NotFoundError
	//   422: ValidationError
	return func(w http.ResponseWriter, r *http.Request) {
		doRead(store.PutMessageRead, w, r)
	}
}

// swagger:model GetReadResponseBody
type getBody struct {
	// Read
	//
	// Required: true
	// example: true
	Read bool `json:"read"`
}

func doRead(inboxFunc func(messageID mail.ID) error, w http.ResponseWriter, r *http.Request) {
	messageID, err := params.PathMessageID(r)
	if err != nil {
		errs.JSONWriter(w, http.StatusNotAcceptable, errors.WithMessage(err, "invalid `message_id`"))
		return
	}

	if err := inboxFunc(messageID); err != nil {
		errs.JSONWriter(w, http.StatusUnprocessableEntity, errors.WithStack(err))
		return
	}

	w.WriteHeader(http.StatusOK)
}
