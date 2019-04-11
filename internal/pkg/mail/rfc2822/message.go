// Copyright (c) 2019 Finobo
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package rfc2822

import (
	"bytes"
	"fmt"
	"io"
	"mime/quotedprintable"
	nm "net/mail"
	"time"

	"github.com/mailchain/mailchain/internal/pkg/mail"
	"github.com/pkg/errors"
)

func messageIDHeaderValue(id mail.ID) string {
	return fmt.Sprintf("%s@mailchain", id.String())
}

func EncodeNewMessage(message *mail.Message) ([]byte, error) {
	// headers["Content-Type"] = "text/html; charset=\"UTF-8\""
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
	headers += "Content-Type: text/plain; charset=\"UTF-8\"\r\n"
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
