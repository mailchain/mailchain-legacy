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

package mail

import (
	"time"
)

// NewHeaders create the headers for sending a new message
func NewHeaders(date time.Time, from Address, to Address, replyTo *Address, subject string) (*Headers, error) {
	return &Headers{
		From:    from,
		To:      to,
		ReplyTo: replyTo,
		Subject: subject,
		Date:    date,
	}, nil
}

// type Headers map[string][]string

// func (h Headers) Address() (*Address, error) {
// 	return nil, nil
// }

// // Date header value
// func (h Headers) Date() (*time.Time, error) {
// 	dateStrings, ok := h["date"]
// 	if !ok {
// 		return nil, errors.Errorf("data header missing")
// 	}
// 	if len(dateStrings) == 0 {
// 		return nil, errors.Errorf("data string is empty")
// 	}
// 	t, err := mail.ParseDate(dateStrings[0])
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &t, nil
// }

// Headers for the message
type Headers struct {
	From    Address
	To      Address
	Date    time.Time
	Subject string
	ReplyTo *Address
}
