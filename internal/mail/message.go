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

package mail

import (
	"time"

	"github.com/pkg/errors"
)

// NewMessage create a new message used when sending a new message
func NewMessage(date time.Time, from, to Address, replyTo *Address, subject string, body []byte, contentType string) (*Message, error) {
	id, err := NewID()

	return &Message{
		ID:      id,
		Headers: NewHeaders(date, from, to, replyTo, subject, contentTypeOrDefault(contentType)),
		Body:    body,
	}, errors.WithMessage(err, "could not create ID")
}

func contentTypeOrDefault(contentType string) string {
	switch contentType {
	case TextContentType, HTMLContentType:
		return contentType
	default:
		return DefaultContentType
	}
}

// Message Mailchain message.
type Message struct {
	Headers *Headers
	ID      ID
	Body    []byte
}
