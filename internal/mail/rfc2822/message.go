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

package rfc2822

import (
	"bytes"
	"fmt"
	"io"
	"mime/quotedprintable"
	nm "net/mail"
	"time"

	"github.com/mailchain/mailchain/internal/mail"
	"github.com/pkg/errors"
)

func messageIDHeaderValue(id mail.ID) string {
	return fmt.Sprintf("%s@mailchain", id.HexString())
}

// EncodeNewMessage Mailchain message in a `[]byte` based on rfc2822 standard.
func EncodeNewMessage(message *mail.Message) ([]byte, error) {
	// TODO: implement base64 encoding too
	// headers["Content-Transfer-Encoding"] = "base64"

	// Setup message
	if message == nil {
		return nil, errors.Errorf("nil message")
	}
	// TODO: Check if valid
	headers := ""
	headers += fmt.Sprintf("Date: %s\r\n", message.Headers.Date.Format(time.RFC1123))
	headers += fmt.Sprintf("Message-ID: %s\r\n", messageIDHeaderValue(message.ID))
	headers += fmt.Sprintf("Subject: %s\r\n", message.Headers.Subject)
	headers += fmt.Sprintf("From: %s\r\n", message.Headers.From.String())
	headers += fmt.Sprintf("To: %s\r\n", message.Headers.To.String())
	if message.Headers.ReplyTo != nil {
		headers += fmt.Sprintf("Reply-To: %s\r\n", message.Headers.ReplyTo.String())
	}
	headers += fmt.Sprintf("Content-Type: %s\r\n", message.Headers.ContentType)
	headers += "Content-Transfer-Encoding: quoted-printable\r\n"
	// 	// Thread-Topic TODO:
	var ac bytes.Buffer
	w := quotedprintable.NewWriter(&ac)

	if _, err := w.Write(message.Body); err != nil {
		return nil, errors.Wrap(err, "could not create writer")
	}
	if err := w.Close(); err != nil {
		return nil, errors.Wrap(err, "could not close writer")
	}
	// TODO: this does not look efficient
	concat := []byte(headers + "\r\n" + ac.String() + "\r\n")
	ret := make([]byte, len(concat))
	copy(ret, concat)
	return ret, nil
}

// DecodeNewMessage from reader into a Mailchain message based on rfc2822 standard.
func DecodeNewMessage(r io.Reader) (*mail.Message, error) {
	netMsg, err := nm.ReadMessage(r)
	if err != nil {
		return nil, errors.Wrap(err, "could not read message")
	}
	h, err := parseHeaders(netMsg.Header)
	if err != nil {
		return nil, errors.Wrap(err, "could not parse headers")
	}
	id, err := parseID(netMsg.Header)
	if err != nil {
		return nil, errors.Wrap(err, "could not parse message id")
	}
	quotedReader := quotedprintable.NewReader(netMsg.Body)
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(quotedReader); err != nil {
		return nil, err
	}
	msg := &mail.Message{
		Body:    buf.Bytes(),
		Headers: h,
		ID:      id,
	}
	return msg, nil
}
