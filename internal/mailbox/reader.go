// Copyright 2022 Mailchain Ltd.
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

	"github.com/mailchain/mailchain/crypto/cipher"
	"github.com/mailchain/mailchain/encoding"
	"github.com/mailchain/mailchain/internal/envelope"
	"github.com/mailchain/mailchain/internal/hash"
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
func ReadMessage(txData []byte, decrypter cipher.Decrypter, cache stores.Cache) (*mail.Message, error) {
	data, err := envelope.Unmarshal(txData)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to unmarshal")
	}

	url, err := data.URL(decrypter)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get URL")
	}

	integrityHash, err := data.IntegrityHash(decrypter)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get integrityHash")
	}

	toDecrypt, err := stores.GetMessage(url.String(), integrityHash, cache)
	if err != nil {
		return nil, errors.WithMessagef(err, "could not get message from %q", url.String())
	}

	rawMsg, err := decrypter.Decrypt(toDecrypt)
	if err != nil {
		return nil, errors.WithMessage(err, "could not decrypt message")
	}

	contentsHash, err := data.ContentsHash(decrypter)
	if err != nil {
		return nil, errors.WithMessage(err, "could not get hash")
	}

	if len(contentsHash) != 0 {
		messageHash := hash.CreateMessageHash(rawMsg)
		if !bytes.Equal(messageHash, contentsHash) {
			return nil, errors.Errorf("contents-hash invalid: message-hash = %v contents-hash = %v",
				encoding.EncodeHex(messageHash), encoding.EncodeHex(contentsHash))
		}
	}

	return rfc2822.DecodeNewMessage(bytes.NewReader(rawMsg))
}
