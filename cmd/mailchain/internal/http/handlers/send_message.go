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
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/mailchain/mailchain/cmd/mailchain/internal/http/params"
	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/cipher/aes256cbc"
	"github.com/mailchain/mailchain/crypto/secp256k1"
	"github.com/mailchain/mailchain/errs"
	"github.com/mailchain/mailchain/internal/address"
	"github.com/mailchain/mailchain/internal/chains"
	"github.com/mailchain/mailchain/internal/envelope"
	"github.com/mailchain/mailchain/internal/keystore"
	"github.com/mailchain/mailchain/internal/keystore/kdf/multi"
	"github.com/mailchain/mailchain/internal/mail"
	"github.com/mailchain/mailchain/internal/mailbox"
	"github.com/mailchain/mailchain/sender"
	"github.com/mailchain/mailchain/stores"
	"github.com/pkg/errors"
)

// SendMessage handler http
func SendMessage(sent stores.Sent, senders map[string]sender.Message, ks keystore.Store, deriveKeyOptions multi.OptionsBuilders) func(
	w http.ResponseWriter, r *http.Request) {
	encrypter := aes256cbc.NewEncrypter()
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
		protocol := "ethereum"
		ctx := r.Context()
		req, err := parsePostRequest(r)
		if err != nil {
			errs.JSONWriter(w, http.StatusUnprocessableEntity, errors.WithStack(err))
			return
		}
		from, err := address.DecodeByProtocol(req.from.ChainAddress, protocol)
		if err != nil {
			errs.JSONWriter(w, http.StatusInternalServerError, errors.WithMessage(err, "failed to decode address"))
			return
		}
		if !ks.HasAddress(from) {
			errs.JSONWriter(w, http.StatusNotAcceptable, errors.Errorf("no private key found for `%s` from address", req.Message.Headers.From))
			return
		}
		sender, ok := senders[fmt.Sprintf("ethereum/%s", req.network)]
		if !ok {
			errs.JSONWriter(w, http.StatusUnprocessableEntity, errors.Errorf("no sender for ethereum/%s configured", req.network))
			return
		}

		msg, err := mail.NewMessage(time.Now(), *req.from, *req.to, req.replyTo, req.Message.Subject, []byte(req.Message.Body))
		if err != nil {
			errs.JSONWriter(w, http.StatusUnprocessableEntity, errors.WithStack(err))
			return
		}
		// TODO: signer is hard coded to ethereum
		signer, err := ks.GetSigner(from, protocol, deriveKeyOptions)
		if err != nil {
			errs.JSONWriter(w, http.StatusUnprocessableEntity, errors.WithStack(errors.WithMessage(err, "could not get `signer`")))
			return
		}

		if err := mailbox.SendMessage(ctx, protocol, req.network, msg, req.publicKey, encrypter, sender, sent, signer, envelope.Kind0x01); err != nil {
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

	return &req, isValid(&req, params.PathNetwork(r))
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

// swagger:model SendMessageRequestBody
type PostRequestBody struct {
	// required: true
	Message   PostMessage `json:"message"`
	to        *mail.Address
	from      *mail.Address
	replyTo   *mail.Address
	publicKey crypto.PublicKey
	network   string
}

func checkForEmpties(msg PostMessage) error {
	if msg.Headers == nil {
		return errors.Errorf("headers must not be nil")
	}
	if msg.Body == "" {
		return errors.Errorf("`body` can not be empty")
	}
	if msg.Subject == "" {
		return errors.Errorf("`subject` can not be empty")
	}
	if msg.PublicKey == "" {
		return errors.Errorf("`public-key` can not be empty")
	}
	return nil
}

func isValid(p *PostRequestBody, protocol string) error {
	if p == nil {
		return errors.New("PostRequestBody must not be nil")
	}
	if err := checkForEmpties(p.Message); err != nil {
		return err
	}
	var err error
	p.network = protocol
	chain := chains.Ethereum

	p.to, err = mail.ParseAddress(p.Message.Headers.To, chain, p.network)
	if err != nil {
		return errors.WithMessage(err, "`to` is invalid")
	}
	// TODO: figure this out
	// if !ethereup.IsAddressValid(p.to.ChainAddress) {
	// 	return errors.Errorf("'address' is invalid")
	// }
	p.from, err = mail.ParseAddress(p.Message.Headers.From, chain, p.network)
	if err != nil {
		return errors.WithMessage(err, "`from` is invalid")
	}

	if p.Message.Headers.ReplyTo != "" {
		p.replyTo, err = mail.ParseAddress(p.Message.Headers.ReplyTo, chain, p.network)
		if err != nil {
			return errors.WithMessage(err, "`reply-to` is invalid")
		}
	}

	// TODO: be more general when getting key from hex
	p.publicKey, err = secp256k1.PublicKeyFromHex(p.Message.PublicKey)
	if err != nil {
		return errors.WithMessage(err, "invalid `public-key`")
	}
	pkAddress := p.publicKey.Address()
	toAddress, err := address.DecodeByProtocol(p.to.ChainAddress, protocol)
	if err != nil {
		return errors.WithMessage(err, "failed to decode address")
	}

	if !bytes.Equal(pkAddress, toAddress) {
		return errors.Errorf("`public-key` does not match to address")
	}

	return nil
}
