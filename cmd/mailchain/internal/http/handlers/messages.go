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
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/http/params"
	"github.com/mailchain/mailchain/crypto/cipher"
	"github.com/mailchain/mailchain/errs"
	"github.com/mailchain/mailchain/internal/address"
	"github.com/mailchain/mailchain/internal/keystore"
	"github.com/mailchain/mailchain/internal/keystore/kdf/multi"
	"github.com/mailchain/mailchain/internal/mailbox"
	"github.com/mailchain/mailchain/stores"
	"github.com/pkg/errors"
)

// GetMessages returns a handler get spec
func GetMessages(inbox stores.State, receivers map[string]mailbox.Receiver, ks keystore.Store,
	deriveKeyOptions multi.OptionsBuilders) func(w http.ResponseWriter, r *http.Request) {
	// Get swagger:route GET /ethereum/{network}/address/{address}/messages Messages Ethereum GetMessages
	//
	// Get Messages.
	//
	// Get mailchain messages.
	// Responses:
	//   200: GetMessagesResponse
	//   404: NotFoundError
	//   422: ValidationError
	return func(w http.ResponseWriter, r *http.Request) {
		protocol := "ethereum"
		ctx := r.Context()
		req, err := parseGetMessagesRequest(r)
		if err != nil {
			errs.JSONWriter(w, http.StatusUnprocessableEntity, errors.WithStack(err))
			return
		}
		receiver, ok := receivers[fmt.Sprintf("%s/%s", protocol, req.Network)]
		if !ok {
			errs.JSONWriter(w, http.StatusUnprocessableEntity, errors.Errorf("no receiver for chain.network configured"))
			return
		}

		addr, err := address.DecodeByProtocol(req.Address, protocol)
		if err != nil {
			errs.JSONWriter(w, http.StatusUnprocessableEntity, errors.WithStack(err))
			return
		}
		if !ks.HasAddress(addr) {
			errs.JSONWriter(w, http.StatusNotAcceptable, errors.Errorf("no private key found for address"))
			return
		}
		encryptedSlice, err := receiver.Receive(ctx, req.Network, addr)
		if mailbox.IsNetworkNotSupportedError(err) {
			errs.JSONWriter(w, http.StatusNotAcceptable, errors.Errorf("network `%s` does not have etherscan client configured", req.Network))
			return
		}
		if err != nil {
			errs.JSONWriter(w, http.StatusInternalServerError, errors.WithStack(err))
			return
		}
		decrypter, err := ks.GetDecrypter(addr, cipher.AES256CBC, deriveKeyOptions)
		if err != nil {
			errs.JSONWriter(w, http.StatusInternalServerError, errors.WithMessage(err, "could not get `decrypter`"))
			return
		}
		messages := []getMessage{}
		for _, transactionData := range encryptedSlice { // TODO: thats an arbitrary limit
			message, err := mailbox.ReadMessage(transactionData, decrypter)
			if err != nil {
				messages = append(messages, getMessage{
					Status: err.Error(),
				})
				continue
			}
			readStatus, _ := inbox.GetReadStatus(message.ID)
			messages = append(messages, getMessage{
				Body: string(message.Body),
				Headers: &getHeaders{
					To:        message.Headers.To.String(),
					From:      message.Headers.From.String(),
					Date:      message.Headers.Date,
					MessageID: message.ID.HexString(),
				},
				Read:    readStatus,
				Subject: message.Headers.Subject,
				Status:  "ok",
			})
		}

		if err := json.NewEncoder(w).Encode(getResponse{Messages: messages}); err != nil {
			errs.JSONWriter(w, http.StatusInternalServerError, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
	}
}

// GetMessagesRequest get mailchain messages
// swagger:parameters GetMessages
type GetMessagesRequest struct {
	// address to query
	//
	// in: path
	// required: true
	// example: 0xde0b295669a9fd93d5f28d9ec85e40f4cb697bae
	// pattern: 0x[a-fA-F0-9]{40}
	Address string `json:"address"`

	// Network
	//
	// enum: mainnet,ropsten,rinkeby,local
	// in: path
	// required: true
	// example: ropsten
	Network string `json:"network"`
}

// ParseGetRequest get all the details for the get request
func parseGetMessagesRequest(r *http.Request) (*GetMessagesRequest, error) {
	addr := strings.ToLower(mux.Vars(r)["address"])
	if addr == "" {
		return nil, errors.Errorf("'address' must not be empty")
	}
	// TODO: validate address
	// if !ethereum.IsAddressValid(addr) {
	// 	return nil, errors.Errorf("'address' is invalid")
	// }

	req := &GetMessagesRequest{
		Address: addr,
		Network: params.PathNetwork(r),
	}
	return req, nil
}

// GetResponse Holds the response messages
//
// swagger:response GetMessagesResponse
type getResponse struct {
	// in: body
	Messages []getMessage `json:"messages"`
}

// swagger:model GetMessagesResponseMessage
type getMessage struct {
	// Headers
	// readOnly: true
	Headers *getHeaders `json:"headers,omitempty"`
	// Body of the mail message
	// readOnly: true
	// example: Lorem ipsum dolor sit amet, consectetur adipiscing elit. Curabitur maximus metus ante, sit amet ullamcorper dui hendrerit ac.
	// Sed vestibulum dui lectus, quis eleifend urna mollis eu. Integer dictum metus ut sem rutrum aliquet.
	Body string `json:"body,omitempty"`
	// Subject of the mail message
	// readOnly: true
	// example: Hello world
	Subject string `json:"subject,omitempty"`
	// readOnly: true
	Status string `json:"status"`
	// readOnly: true
	StatusCode string `json:"status-code"`
	// Read status of the message
	// readOnly: true
	// example: true
	Read bool `json:"read"`
}

// swagger:model GetMessagesResponseHeaders
type getHeaders struct {
	// When the message was created, this can be different to the transaction data of the message.
	// readOnly: true
	// example: 12 Mar 19 20:23 UTC
	Date time.Time `json:"date"`
	// The sender of the message
	// readOnly: true
	// example: Charlotte <5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761@ropsten.ethereum>
	From string `json:"from"`
	// The recipient of the message
	// readOnly: true
	// To: <4cb0a77b76667dac586c40cc9523ace73b5d772bd503c63ed0ca596eae1658b2@ropsten.ethereum>
	To string `json:"to"`
	// Reply to if the reply address is different to the from address.
	// readOnly: true
	ReplyTo string `json:"reply-to,omitempty"`
	// Unique identifier of the message
	// example: 47eca011e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471@mailchain
	// readOnly: true
	MessageID string `json:"message-id"`
}
