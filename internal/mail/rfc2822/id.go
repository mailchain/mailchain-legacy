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

package rfc2822

import (
	"encoding/hex"
	nm "net/mail"
	"strings"

	"github.com/mailchain/mailchain/internal/mail"
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

	decoded, err := hex.DecodeString(strings.TrimSuffix(messageID, "@mailchain"))
	if err != nil {
		return nil, errors.WithMessage(err, "failed to decode")
	}
	return mail.ID(decoded), nil
}
