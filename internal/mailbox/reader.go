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
	"bytes"

	"github.com/gogo/protobuf/proto"
	"github.com/mailchain/mailchain/internal/crypto"
	"github.com/mailchain/mailchain/internal/crypto/cipher"
	"github.com/mailchain/mailchain/internal/encoding"
	"github.com/mailchain/mailchain/internal/mail"
	"github.com/mailchain/mailchain/internal/mail/rfc2822"
	"github.com/mailchain/mailchain/internal/stores"
	"github.com/pkg/errors"
)

// ReadMessage gets the messages, decrypts and checks to see if it's valid
// - Check transaction data
// - Decrypt location
// - Get message
// - Decrypt message
// - Check hash
func ReadMessage(txData []byte, decrypter cipher.Decrypter) (*mail.Message, error) {
	if txData[0] != encoding.Protobuf {
		return nil, errors.Errorf("invalid encoding prefix")
	}
	var data mail.Data
	if err := proto.Unmarshal(txData[1:], &data); err != nil {
		return nil, errors.WithMessage(err, "could not unmarshal to data")
	}

	decryptedLocation, err := decryptLocation(&data, decrypter)
	if err != nil {
		return nil, err
	}
	toDecrypt, err := stores.GetMessage(decryptedLocation)
	if err != nil {
		return nil, errors.WithMessage(err, "could not get message from `location`")
	}
	rawMsg, err := decrypter.Decrypt(toDecrypt)
	if err != nil {
		return nil, errors.WithMessage(err, "could not decrypt message")
	}
	messageHash, err := crypto.CreateMessageHash(rawMsg)
	if err != nil {
		return nil, errors.WithMessage(err, "could not create message hash")
	}
	if !bytes.Equal(messageHash, data.Hash) {
		return nil, errors.Errorf("message-hash invalid")
	}
	return rfc2822.DecodeNewMessage(bytes.NewReader(rawMsg))
}

// decryptLocation return the location in readable form
func decryptLocation(d *mail.Data, decrypter cipher.Decrypter) (string, error) {
	decryptedLocation, err := decrypter.Decrypt(d.EncryptedLocation)
	if err != nil {
		return "", errors.WithMessage(err, "could not decrypt location")
	}
	return string(decryptedLocation), nil
}
