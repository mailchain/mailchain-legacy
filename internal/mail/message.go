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
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// NewMessage create a new message used when sending a new message
func NewMessage(date time.Time, from, to Address, replyTo *Address, subject string, body []byte) (*Message, error) {
	id, err := NewID()

	return &Message{
		ID:      id,
		Headers: NewHeaders(date, from, to, replyTo, subject, detectContentType(body)),
		Body:    body,
	}, errors.WithMessage(err, "could not create ID")
}

func detectContentType(body []byte) string {
	contentType := http.DetectContentType(body)
	result := strings.Split(contentType, ";")

	if len(result) == 1 {
		return result[0]
	} else if len(result) == 2 {
		encodingParts := strings.Split(result[1], "=")
		encoding := fmt.Sprintf("%s=\"%s\"", encodingParts[0], strings.ToUpper(encodingParts[1]))
		return fmt.Sprintf("%s;%s", result[0], encoding)
	}

	return DefaultContentType
}

// Message Mailchain message.
type Message struct {
	Headers *Headers
	ID      ID
	Body    []byte
}
