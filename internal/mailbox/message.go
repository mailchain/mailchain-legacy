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

package mailbox

import (
	"context"

	"github.com/mailchain/mailchain/crypto/cipher"
	"github.com/mailchain/mailchain/encoding"
	"github.com/mailchain/mailchain/internal/address"
	"github.com/mailchain/mailchain/internal/envelope"
	"github.com/mailchain/mailchain/internal/hash"
	"github.com/mailchain/mailchain/internal/mail"
	"github.com/mailchain/mailchain/internal/mail/rfc2822"
	"github.com/mailchain/mailchain/internal/mailbox/signer"
	"github.com/mailchain/mailchain/sender"
	"github.com/mailchain/mailchain/stores"
	"github.com/pkg/errors"
)

// SendMessage performs all the actions required to send a message.
// - Create a hash of encoded message
// - Encrypt message
// - Store sent message
// - Encrypt message location
// - Create transaction data with encrypted location and message hash
// - Send transaction
func SendMessage(ctx context.Context, protocol, network string, msg *mail.Message, encrypter cipher.Encrypter,
	msgSender sender.Message, sent stores.Sent, msgSigner signer.Signer, envelopeKind byte) error {
	encodedMsg, err := rfc2822.EncodeNewMessage(msg)
	if err != nil {
		return errors.WithMessage(err, "could not encode message")
	}

	encrypted, err := encrypter.Encrypt(encodedMsg)
	if err != nil {
		return errors.WithMessage(err, "could not encrypt mail message")
	}

	msgAddress, resource, mli, err := sent.PutMessage(msg.ID, hash.CreateMessageHash(encodedMsg), encrypted, nil)
	if err != nil {
		return errors.WithMessage(err, "failed to store message")
	}

	locOpt, err := envelope.WithMessageLocationIdentifier(mli)
	if err != nil {
		return errors.WithStack(err)
	}

	opts := []envelope.CreateOptionsBuilder{
		envelope.WithKind(envelopeKind),
		envelope.WithURL(msgAddress),
		envelope.WithResource(resource),
		envelope.WithEncryptedContents(encrypted),
		envelope.WithDecryptedHash(hash.CreateMessageHash(encodedMsg)),
		envelope.WithEncryptedHash(hash.CreateIntegrityHash(encrypted)),
		locOpt,
	}

	env, err := envelope.NewEnvelope(encrypter, opts)
	if err != nil {
		return errors.WithMessage(err, "could not create envelope")
	}

	encodedData, err := envelope.Marshal(env)
	if err != nil {
		return errors.WithMessage(err, "could not marshal envelope")
	}

	transactonData := append(encoding.DataPrefix(), encodedData...)
	to, err := address.DecodeByProtocol(msg.Headers.To.ChainAddress, protocol)

	if err != nil {
		return errors.WithMessage(err, "could not decode to address")
	}

	from, err := address.DecodeByProtocol(msg.Headers.From.ChainAddress, protocol)
	if err != nil {
		return errors.WithMessage(err, "could not decode from address")
	}

	if err := msgSender.Send(ctx, network, to, from, transactonData, msgSigner, nil); err != nil {
		return errors.WithMessage(err, "could not send transaction")
	}

	return nil
}
