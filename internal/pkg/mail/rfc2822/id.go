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
	"encoding/hex"
	nm "net/mail"
	"strings"

	"github.com/mailchain/mailchain/internal/pkg/mail"
	"github.com/multiformats/go-multihash"
	"github.com/pkg/errors"
)

func parseID(h nm.Header) (mail.ID, error) {
	sources, ok := h["Message-Id"]
	if !ok {
		return nil, errors.Errorf("header missing")
	}
	if len(sources) == 0 {
		return nil, errors.Errorf("empty header")
	}
	messageID := sources[0]
	if !strings.HasSuffix(messageID, "@mailchain") {
		return nil, errors.Errorf("invalid suffix")
	}
	messageID = strings.TrimRight(messageID, "@mailchain")
	decoded, err := hex.DecodeString(messageID)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to decode")
	}
	i, err := multihash.Cast(decoded)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to cast")
	}
	return mail.ID(i), nil
}
