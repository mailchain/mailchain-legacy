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

	"github.com/ethereum/go-ethereum/common"
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

func isValid(p *PostRequestBody, network string) error {
	if p == nil {
		return errors.New("PostRequestBody must not be nil")
	}
	if err := checkForEmpties(p.Message); err != nil {
		return err
	}
	var err error
	p.network = network
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

	// TODO: be more general when getting key from hex
	p.publicKey, err = secp256k1.PublicKeyFromHex(p.Message.PublicKey)
	if err != nil {
		return errors.WithMessage(err, "invalid `public-key`")
	}
	pkAddress := p.publicKey.Address()
	toAddress := common.HexToAddress(p.to.ChainAddress).Bytes()
	if !bytes.Equal(pkAddress, toAddress) {
		return errors.Errorf("`public-key` does not match to address")
	}

	return nil
}
