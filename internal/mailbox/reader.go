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

	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/cipher"
	"github.com/mailchain/mailchain/internal/envelope"
	"github.com/mailchain/mailchain/internal/mail"
	"github.com/mailchain/mailchain/internal/mail/rfc2822"
	"github.com/mailchain/mailchain/stores"
	"github.com/pkg/errors"
)

// ReadMessage gets the messages, decrypts and checks to see if it's valid
// - Check transaction data
// - Decrypt location
// - Get message
// - Decrypt message
// - Check hash
func ReadMessage(txData []byte, decrypter cipher.Decrypter) (*mail.Message, error) {
	data, err := envelope.Unmarshal(txData)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to unmarshal")
	}
	url, err := data.URL(decrypter)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get URL")
	}

	toDecrypt, err := stores.GetMessage(url.String())
	if err != nil {
		return nil, errors.WithMessagef(err, "could not get message from %q", url.String())
	}
	rawMsg, err := decrypter.Decrypt(toDecrypt)
	if err != nil {
		return nil, errors.WithMessage(err, "could not decrypt message")
	}
	hash, err := data.ContentsHash(decrypter)
	if err != nil {
		return nil, errors.WithMessage(err, "could not get hash")
	}
	if len(hash) != 0 {
		messageHash := crypto.CreateMessageHash(rawMsg)
		if !bytes.Equal(messageHash, hash) {
			return nil, errors.Errorf("message-hash invalid")
		}
	}

	return rfc2822.DecodeNewMessage(bytes.NewReader(rawMsg))
}
