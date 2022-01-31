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

	"github.com/mailchain/mailchain/stores"
	"github.com/mailchain/mailchain/stores/pinata"
	"github.com/mailchain/mailchain/stores/s3store"
)

// GetEnvelope returns the available envelopes.
func GetEnvelope(sent stores.Sent) func(w http.ResponseWriter, r *http.Request) {
	zeroX01 := GetEnvelopeResponseBodyElement{
		Type:        "0x01",
		Description: "Private Message Stored with MLI",
	}
	zeroX02 := GetEnvelopeResponseBodyElement{
		Type:        "0x02",
		Description: "Private Message Stored on IPFS",
	}
	// Get swagger:route GET /envelope Envelope GetEnvelope
	//
	// Get Mailchain envelope
	//
	// Get envelope
	// This method returns the available envelope types
	//
	// Responses:
	//   200: GetEnvelopeResponse
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		body := []GetEnvelopeResponseBodyElement{}
		switch sent.(type) {
		case *s3store.Sent:
			body = append(body, zeroX01)
		case *stores.SentStore:
			body = append(body, zeroX01)
		case *pinata.Sent:
			body = append(body, zeroX02)
		}

		_ = json.NewEncoder(w).Encode(body)
	}
}

// GetEnvelopeResponse envelope response
//
// swagger:response GetEnvelopeResponse
type GetEnvelopeResponse struct {
	// in: body
	Body []GetEnvelopeResponseBodyElement
}

// GetEnvelopeResponseBodyElement response
//
// swagger:model GetEnvelopeResponseBodyElement
type GetEnvelopeResponseBodyElement struct {
	// The envelope type
	// Required: true
	// example: 0x01
	Type string `json:"type"`
	// The envelope description
	// Required: true
	// example: Private Message Stored with MLI
	Description string `json:"description"`
}
