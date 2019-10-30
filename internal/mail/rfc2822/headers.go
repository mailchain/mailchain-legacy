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
	nm "net/mail"
	"time"

	"github.com/mailchain/mailchain/internal/mail"
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
		Date:        *date,
		Subject:     subject,
		To:          *to,
		From:        *from,
		ContentType: parseContentType(h),
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

func parseContentType(h nm.Header) string {
	sources, ok := h["Content-Type"]
	if !ok || len(sources) == 0 || sources[0] == "" {
		return mail.DefaultContentType
	}

	return sources[0]
}
