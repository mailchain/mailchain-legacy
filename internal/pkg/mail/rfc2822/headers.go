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
	nm "net/mail"
	"time"

	"github.com/mailchain/mailchain/internal/pkg/mail"
	"github.com/pkg/errors"
)

func parseHeaders(h nm.Header) (*mail.Headers, error) {
	date, err := parseDate(h)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to parse `date`")
	}
	subject, err := parseSubject(h)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to parse `subject`")
	}
	to, err := parseTo(h)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to parse `to`")
	}
	from, err := parseFrom(h)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to parse `from`")
	}
	return &mail.Headers{
		Date:    *date,
		Subject: subject,
		To:      *to,
		From:    *from,
	}, nil
}

func parseTo(h nm.Header) (*mail.Address, error) {
	sources, ok := h["To"]
	if !ok {
		return nil, errors.Errorf("header missing")
	}
	if len(sources) == 0 {
		return nil, errors.Errorf("empty header")
	}

	return mail.ParseAddress(sources[0], "", "")
}
func parseFrom(h nm.Header) (*mail.Address, error) {
	sources, ok := h["From"]
	if !ok {
		return nil, errors.Errorf("header missing")
	}
	if len(sources) == 0 {
		return nil, errors.Errorf("empty header")
	}

	return mail.ParseAddress(sources[0], "", "")
}
func parseDate(h nm.Header) (*time.Time, error) {
	dateStrings, ok := h["Date"]
	if !ok {
		return nil, errors.Errorf("header missing")
	}
	if len(dateStrings) == 0 {
		return nil, errors.Errorf("empty header")
	}
	t, err := nm.ParseDate(dateStrings[0])
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func parseSubject(h nm.Header) (string, error) {
	sources, ok := h["Subject"]
	if !ok {
		return "", errors.Errorf("header missing")
	}
	if len(sources) == 0 {
		return "", errors.Errorf("empty header")
	}

	return sources[0], nil
}
