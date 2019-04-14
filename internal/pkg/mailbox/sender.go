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

	"github.com/ethereum/go-ethereum/common"
	"github.com/gogo/protobuf/proto"
	"github.com/mailchain/mailchain/internal/pkg/crypto"
	"github.com/mailchain/mailchain/internal/pkg/crypto/cipher/aes256cbc"
	"github.com/mailchain/mailchain/internal/pkg/crypto/keys"
	"github.com/mailchain/mailchain/internal/pkg/encoding"
	"github.com/mailchain/mailchain/internal/pkg/mail"
	"github.com/mailchain/mailchain/internal/pkg/mail/rfc2822"
	"github.com/mailchain/mailchain/internal/pkg/stores"
	"github.com/pkg/errors"
)

// Sender signs a transaction the sends it
type Sender interface {
	Send(ctx context.Context, to []byte, from []byte, data []byte, signer Signer, opts SenderOpts) (err error)
}

// SenderOpts options for sending a message
type SenderOpts interface{}

// SendMessage performs all the actions required to send a message.
// - Create a hash of encoded message
// - Encrypt message
// - Store sent message
// - Encrypt message location
// - Create transaction data with encrypted location and message hash
// - Send transaction
func SendMessage(ctx context.Context, msg *mail.Message, recipientKey keys.PublicKey,
	sender Sender, sent stores.Sent, signer Signer) error {
	encodedMsg, err := rfc2822.EncodeNewMessage(msg)
	if err != nil {
		return errors.WithMessage(err, "could not encode message")
	}
	msgHash, err := crypto.CreateMessageHash(encodedMsg)
	if err != nil {
		return errors.WithMessage(err, "could not create message hash")
	}

	encrypted, err := encryptMailMessage(recipientKey, encodedMsg)
	if err != nil {
		return errors.WithStack(err)
	}

	location, err := stores.PutMessage(sent, msg.ID, encrypted)
	if err != nil {
		return errors.WithStack(err)
	}
	encryptedLocation, err := encryptLocation(recipientKey, location)
	if err != nil {
		return errors.WithMessage(err, "could not create transaction data")
	}

	data := &mail.Data{
		EncryptedLocation: encryptedLocation,
		Hash:              msgHash,
	}
	encodedData, err := prefixedBytes(data)
	if err != nil {
		return errors.WithMessage(err, "could not encode transaction data")
	}
	transactonData := append(encoding.DataPrefix(), encodedData...)
	//TODO: should not use common to parse address
	to := common.FromHex(msg.Headers.To.ChainAddress)
	from := common.FromHex(msg.Headers.From.ChainAddress)
	if err := sender.Send(ctx, to, from, transactonData, signer, nil); err != nil {
		return errors.WithMessage(err, "could not send transaction")
	}

	return nil
}

func prefixedBytes(data proto.Message) ([]byte, error) {
	protoData, err := proto.Marshal(data)
	if err != nil {
		return nil, errors.WithMessage(err, "could not marshal data")
	}

	prefixedProto := make([]byte, len(protoData)+1)
	prefixedProto[0] = encoding.Protobuf
	copy(prefixedProto[1:], protoData)

	return prefixedProto, nil
}

// encryptLocation is encrypted with supplied public key and location string
func encryptLocation(pk keys.PublicKey, location string) ([]byte, error) {
	// TODO: encryptLocation hard coded to aes256cbc
	encryptedLocation, err := aes256cbc.Encrypt(pk, []byte(location))
	if err != nil {
		return nil, errors.WithMessage(err, "could not encrypt data")
	}
	return encryptedLocation, nil
}

// encryptMailMessage is encrypted with supplied public key and location string
func encryptMailMessage(pk keys.PublicKey, encodedMsg []byte) ([]byte, error) {
	// TODO: encryptMailMessage hard coded to aes256cbc
	encryptedData, err := aes256cbc.Encrypt(pk, encodedMsg)
	if err != nil {
		return nil, errors.WithMessage(err, "could not encrypt message")
	}

	return encryptedData, nil
}
