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
	"io/ioutil"
	"net/url"
	"strings"

	"github.com/gogo/protobuf/proto"
	"github.com/mailchain/mailchain/internal/pkg/crypto"
	"github.com/mailchain/mailchain/internal/pkg/crypto/cipher"
	"github.com/mailchain/mailchain/internal/pkg/encoding"
	"github.com/mailchain/mailchain/internal/pkg/mail"
	"github.com/mailchain/mailchain/internal/pkg/mail/rfc2822"
	"github.com/pkg/errors"
	"gopkg.in/resty.v1"
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
	toDecrypt, err := getMessage(decryptedLocation)
	if err != nil {
		return nil, errors.WithMessage(err, "could not get message from `location`")
	}
	rawMsg, err := decrypter.Decrypt(toDecrypt)
	if err != nil {
		return nil, errors.WithMessage(err, "could not decrypt transaction data")
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

// getMessage get the message contents from the location and perform location hash check
func getMessage(location string) ([]byte, error) {
	msg, err := getAnyMessage(location)
	if err != nil {
		return nil, err
	}

	hash, err := crypto.CreateLocationHash(msg)
	if err != nil {
		return nil, err
	}
	parts := strings.Split(location, "-")
	if len(parts) < 1 {
		return nil, errors.Errorf("could not safely extract hash from location")
	}
	if hash.String() != parts[len(parts)-1] {
		return nil, errors.Errorf("hash does not match contents")
	}
	return msg, nil
}

func getAnyMessage(location string) ([]byte, error) {
	parsed, err := url.Parse(location)
	if err != nil {
		return nil, err
	}

	switch parsed.Scheme {
	case "http":
		return getHTTPMessage(location)
	case "file":
		return ioutil.ReadFile(parsed.Host + parsed.Path)
	case "test":
		return []byte(parsed.Host), nil
	default:
		return nil, errors.Errorf("unsupported scheme")
	}
}
func getHTTPMessage(location string) ([]byte, error) {
	res, err := resty.R().Get(location)
	if err != nil {
		return nil, errors.Wrap(err, "could not get message from `location`")
	}
	msg := res.Body()

	return msg, nil
}
