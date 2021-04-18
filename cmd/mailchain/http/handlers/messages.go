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
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/mailchain/mailchain/cmd/internal/http/params"
	"github.com/mailchain/mailchain/encoding"
	"github.com/mailchain/mailchain/errs"
	"github.com/mailchain/mailchain/internal/addressing"
	"github.com/mailchain/mailchain/internal/envelope"
	"github.com/mailchain/mailchain/internal/keystore"
	"github.com/mailchain/mailchain/internal/keystore/kdf/multi"
	"github.com/mailchain/mailchain/internal/mailbox"
	"github.com/mailchain/mailchain/stores"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// GetMessages returns a handler get spec.
func GetMessages(receivers map[string]mailbox.Receiver, inbox stores.State, cache stores.Cache, ks keystore.Store, deriveKeyOptions multi.OptionsBuilders) func(w http.ResponseWriter, r *http.Request) { //nolint: funlen, gocyclo
	// Get swagger:route GET /messages Messages GetMessages
	//
	// Get Mailchain messages.
	//
	// Check the protocol, network, address combination for Mailchain messages.
	// Responses:
	//   200: GetMessagesResponse
	//   422: ValidationError
	return func(w http.ResponseWriter, r *http.Request) {
		logger := log.Logger.With().Str("component", "http").Logger()

		req, err := parseGetMessagesRequest(r)
		if err != nil {
			errs.JSONWriter(w, r, http.StatusUnprocessableEntity, errors.WithStack(err), logger)
			return
		}

		if !ks.HasAddress(req.addressBytes, req.Protocol, req.Network) {
			errs.JSONWriter(w, r, http.StatusNotAcceptable, errors.Errorf("no private key found for address"), logger)
			return
		}

		if req.Fetch {
			if err := fetchMessages(r.Context(), req.Protocol, req.Network, req.addressBytes, receivers, inbox); err != nil {
				errs.JSONWriter(w, r, http.StatusNotAcceptable, errors.WithMessage(err, "fetch messages failed"), logger)
				return
			}
		}

		txs, err := inbox.GetTransactions(req.Protocol, req.Network, req.addressBytes)
		if err != nil {
			errs.JSONWriter(w, r, http.StatusInternalServerError, errors.WithStack(err), logger)
			return
		}

		messages := make([]getMessage, len(txs))

		for i, tx := range txs {
			buf := make([]byte, binary.MaxVarintLen64)
			n := binary.PutVarint(buf, tx.BlockNumber)
			blockID := encoding.EncodeHexZeroX(buf[:n])
			blockIDEncoding := encoding.KindHex0XPrefix

			env, err := envelope.Unmarshal(tx.EnvelopeData)
			if err != nil {
				messages[i] = getMessage{Status: errors.WithMessagef(err, "failed to unmarshal envelope, %s", encoding.EncodeHexZeroX(tx.EnvelopeData)).Error(), BlockID: blockID, BlockIDEncoding: blockIDEncoding}

				continue
			}

			decrypterKind, err := env.DecrypterKind()
			if err != nil {
				messages[i] = getMessage{Status: errors.WithMessage(err, "failed to find decrypter type").Error(), BlockID: blockID, BlockIDEncoding: blockIDEncoding}

				continue
			}

			decrypter, err := ks.GetDecrypter(req.addressBytes, req.Protocol, req.Network, decrypterKind, deriveKeyOptions)
			if err != nil {
				messages[i] = getMessage{Status: errors.WithMessage(err, "could not get `decrypter`").Error(), BlockID: blockID, BlockIDEncoding: blockIDEncoding}

				continue
			}

			message, err := mailbox.ReadMessage(tx.EnvelopeData, decrypter, cache)
			if err != nil {
				messages[i] = getMessage{Status: errors.WithMessage(err, "could not read message").Error(), BlockID: blockID, BlockIDEncoding: blockIDEncoding}

				continue
			}

			readStatus, _ := inbox.GetReadStatus(message.ID)

			mailStore := &stores.Message{
				Body: string(message.Body),
				Headers: stores.Header{
					To:          message.Headers.To.String(),
					From:        message.Headers.From.String(),
					Date:        message.Headers.Date,
					MessageID:   message.ID.HexString(),
					ContentType: message.Headers.ContentType,
				},
				Read:                    readStatus,
				Subject:                 message.Headers.Subject,
				Status:                  "ok",
				BlockID:                 blockID,
				BlockIDEncoding:         blockIDEncoding,
				TransactionHash:         encoding.EncodeHexZeroX(tx.Hash),
				TransactionHashEncoding: encoding.KindHex0XPrefix,
			}

			if tx.RekeyAddress != nil {
				encodedAddress, _, _ := addressing.EncodeByProtocol(tx.RekeyAddress, req.Protocol)
				mailStore.Headers.RekeyTo = encodedAddress
			}

			messages[i] = convertStoreMessageToGetMessage(mailStore)
		}

		if err := json.NewEncoder(w).Encode(getResponse{Messages: messages}); err != nil {
			errs.JSONWriter(w, r, http.StatusInternalServerError, err, logger)
			return
		}

		w.Header().Set("Content-Type", "application/json")
	}
}

func fetchMessages(ctx context.Context, protocol, network string, address []byte, receivers map[string]mailbox.Receiver, inbox stores.State) error {
	receiver, ok := receivers[fmt.Sprintf("%s/%s", protocol, network)]
	if !ok {
		return errors.Errorf("receiver not supported on \"%s/%s\"", protocol, network)
	}

	if receiver == nil {
		return errors.Errorf("no receiver configured for \"%s/%s\"", protocol, network)
	}

	transactions, err := receiver.Receive(ctx, protocol, network, address)
	if mailbox.IsNetworkNotSupportedError(err) {
		return errors.Errorf("network `%s` does not have etherscan client configured", network)
	} else if err != nil {
		return errors.WithStack(err)
	}

	for i := range transactions {
		tx := transactions[i]
		if err := inbox.PutTransaction(protocol, network, address, tx); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

// GetMessagesRequest get mailchain messages
// swagger:parameters GetMessages
type GetMessagesRequest struct {
	// Address to use when looking for messages.
	//
	// in: query
	// required: true
	// example: 0xde0b295669a9fd93d5f28d9ec85e40f4cb697bae
	// pattern: 0x[a-fA-F0-9]{40}
	Address string `json:"address"`

	// Network to use when looking for messages.
	//
	// enum: mainnet,goerli,ropsten,rinkeby,local
	// in: query
	// required: true
	// example: goerli
	Network string `json:"network"`

	// Protocol to use when looking for messages.
	//
	// enum: algorand, ethereum, substrate
	// in: query
	// required: true
	// example: ethereum
	Protocol string `json:"protocol"`

	// Fetch go to the blockchain to retrieve messages
	//
	// in: query
	// example: true
	// default: false
	Fetch bool `json:"fetch"`

	addressBytes []byte
}

// ParseGetRequest get all the details for the get request.
func parseGetMessagesRequest(r *http.Request) (*GetMessagesRequest, error) {
	protocol, err := params.QueryRequireProtocol(r)
	if err != nil {
		return nil, err
	}

	network, err := params.QueryRequireNetwork(r)
	if err != nil {
		return nil, err
	}

	addr, err := params.QueryRequireAddress(r)
	if err != nil {
		return nil, err
	}

	addressBytes, err := addressing.DecodeByProtocol(addr, protocol)
	if err != nil {
		return nil, err
	}

	res := &GetMessagesRequest{
		Address:      addr,
		addressBytes: addressBytes,
		Network:      network,
		Protocol:     protocol,
	}

	fetch := r.URL.Query()["fetch"]
	if len(fetch) == 1 && fetch[0] == "true" {
		res.Fetch = true
	}

	return res, nil
}

func convertStoreMessageToGetMessage(message *stores.Message) getMessage {
	return getMessage{
		Body: message.Body,
		Headers: &getHeaders{
			To:          message.Headers.To,
			From:        message.Headers.From,
			Date:        message.Headers.Date,
			MessageID:   message.Headers.MessageID,
			ContentType: message.Headers.ContentType,
			RekeyTo:     message.Headers.RekeyTo,
		},
		Read:                    message.Read,
		Subject:                 message.Subject,
		Status:                  message.Status,
		BlockID:                 message.BlockID,
		BlockIDEncoding:         message.BlockIDEncoding,
		TransactionHash:         message.TransactionHash,
		TransactionHashEncoding: message.TransactionHashEncoding,
	}
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
	// Transaction's block number
	// readOnly: true
	BlockID string `json:"block-id,omitempty"`
	// Transaction's block number encoding type used by the specific protocol
	// readOnly: true
	BlockIDEncoding string `json:"block-id-encoding,omitempty"`
	// Transaction's hash
	// readOnly: true
	TransactionHash string `json:"transaction-hash,omitempty"`
	// Transaction's hash encoding type used by the specific protocol
	// readOnly: true
	TransactionHashEncoding string `json:"transaction-hash-encoding,omitempty"`
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
	// Rekey use this key when responding.
	// readOnly: true
	RekeyTo string `json:"rekey-to,omitempty"`
	// Unique identifier of the message
	// example: 47eca011e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471@mailchain
	// readOnly: true
	MessageID string `json:"message-id"`
	// The content type and the encoding of the message body
	// readOnly: true
	// example: text/plain; charset=\"UTF-8\",
	// 			text/html; charset=\"UTF-8\"
	ContentType string `json:"content-type"`
}
