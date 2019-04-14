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
	"bytes"
	"net/http"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"
	"github.com/mailchain/mailchain/internal/pkg/crypto/keys"
	"github.com/mailchain/mailchain/internal/pkg/crypto/keys/secp256k1"
	"github.com/mailchain/mailchain/internal/pkg/encoding"
	"github.com/mailchain/mailchain/internal/pkg/mail"
	"github.com/pkg/errors"
)

// swagger:model SendMessageRequestBody
type PostRequestBody struct {
	// required: true
	Message   PostMessage `json:"message"`
	to        *mail.Address
	from      *mail.Address
	replyTo   *mail.Address
	publicKey keys.PublicKey
	network   string
}

func (p *PostRequestBody) asMessage() (*mail.Message, error) {
	return mail.NewMessage(time.Now(), *p.from, *p.to, p.replyTo, p.Message.Subject, []byte(p.Message.Body))
}

func (p *PostRequestBody) checkForEmpties() error {
	if p.Message.Headers == nil {
		return errors.Errorf("headers must not be nil")
	}
	if p.Message.Body == "" {
		return errors.Errorf("`body` can not be empty")
	}
	if p.Message.Subject == "" {
		return errors.Errorf("`subject` can not be empty")
	}

	if p.Message.PublicKey == "" {
		return errors.Errorf("`public-key` can not be empty")
	}

	return nil
}

func (p *PostRequestBody) isValid(r *http.Request) error {
	if err := p.checkForEmpties(); err != nil {
		return err
	}
	var err error
	p.network = strings.ToLower(mux.Vars(r)["network"])
	chain := encoding.Ethereum

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

	if bytes.Equal(p.publicKey.Address(), common.HexToAddress(p.to.ChainAddress).Bytes()) {
		return errors.WithMessage(err, "`public-key` does not match to address")
	}

	return nil
}
