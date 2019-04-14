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

package send

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"
	"github.com/mailchain/mailchain/internal/pkg/http/rest/errs"
	"github.com/mailchain/mailchain/internal/pkg/keystore"
	"github.com/mailchain/mailchain/internal/pkg/keystore/kdf/multi"
	"github.com/mailchain/mailchain/internal/pkg/mailbox"
	"github.com/mailchain/mailchain/internal/pkg/stores"
	"github.com/pkg/errors"
)

// Post handler http
func Post(store stores.Sender, senders map[string]mailbox.Sender, ks keystore.Store,
	deriveKeyOptions multi.OptionsBuilders) func(w http.ResponseWriter, r *http.Request) {
	// Post swagger:route POST /ethereum/{network}/messages/send Send Ethereum SendMessage
	//
	// Send message.
	//
	// Securely send message to ethereum address that can only be discovered and de-cryted by the private key holder.
	//
	// - Create mailchain message
	// - Encrypt content with public key
	// - Store message
	// - Encrypt location
	// - Store encrypted location on the blockchain.
	//
	// Responses:
	//   200: StatusOK
	//   404: NotFoundError
	//   422: ValidationError
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		req, err := parsePostRequest(r)
		if err != nil {
			errs.JSONWriter(w, http.StatusUnprocessableEntity, errors.WithStack(err))
			return
		}

		if !ks.HasAddress(common.HexToAddress(req.from.ChainAddress).Bytes()) {
			errs.JSONWriter(w, http.StatusNotAcceptable, errors.Errorf("no private key found for `%s` from address", req.Message.Headers.From))
			return
		}
		network := strings.ToLower(mux.Vars(r)["network"])
		sender, ok := senders[fmt.Sprintf("ethereum.%s", network)]
		if !ok {
			errs.JSONWriter(w, http.StatusUnprocessableEntity, errors.Errorf("no sender for chain.network configured"))
			return
		}

		msg, err := req.asMessage()
		if err != nil {
			errs.JSONWriter(w, http.StatusUnprocessableEntity, errors.WithStack(err))
			return
		}
		// TODO: signer is hard coded to ethereum
		signer, err := ks.GetSigner(common.FromHex(msg.Headers.From.ChainAddress), "ethereum", deriveKeyOptions)
		if err != nil {
			errs.JSONWriter(w, http.StatusUnprocessableEntity, errors.WithStack(errors.WithMessage(err, "could not get `signer`")))
			return
		}

		if err := mailbox.SendMessage(ctx, msg, req.publicKey, sender, store, signer); err != nil {
			errs.JSONWriter(w, http.StatusInternalServerError, errors.WithMessage(err, "could not send message"))
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

// PostRequest get mailchain inputs
// swagger:parameters SendMessage
type PostRequest struct {
	// Network
	//
	// enum: mainnet,ropsten,rinkeby,local
	// in: path
	// required: true
	// example: ropsten
	Network string `json:"network"`

	// Message to send
	// in: body
	// required: true
	PostRequestBody PostRequestBody
}

// parsePostRequest post all the details for the message
func parsePostRequest(r *http.Request) (*PostRequestBody, error) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var req PostRequestBody
	if err := decoder.Decode(&req); err != nil {
		return nil, errors.WithMessage(err, "'message' is invalid")
	}

	return &req, req.isValid(r)
}

// swagger:model PostMessagesResponseHeaders
type PostHeaders struct {
	// The sender of the message
	// required: true
	// example: Charlotte <5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761@ropsten.ethereum>
	From string `json:"from"`
	// The recipient of the message
	// required: true
	// To: <4cb0a77b76667dac586c40cc9523ace73b5d772bd503c63ed0ca596eae1658b2@ropsten.ethereum>
	To string `json:"to"`
	// Reply to if the reply address is different to the from address.
	ReplyTo string `json:"reply-to"`
}

// swagger:model PostMessagesResponseMessage
type PostMessage struct {
	// Headers
	// required: true
	// in: body
	Headers *PostHeaders `json:"headers"`
	// Body of the mail message
	// required: true
	// example: Lorem ipsum dolor sit amet, consectetur adipiscing elit. Curabitur maximus metus ante,
	// sit amet ullamcorper dui hendrerit ac. Sed vestibulum dui lectus, quis eleifend urna mollis eu.
	// Integer dictum metus ut sem rutrum aliquet.
	Body string `json:"body"`
	// Subject of the mail message
	// required: true
	// example: Hello world
	Subject string `json:"subject"`
	// Public key of the recipient to encrypt with
	// required: true
	PublicKey string `json:"public-key"`
}
