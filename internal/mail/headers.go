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
)

// DefaultContentType value.
const DefaultContentType = "text/plain; charset=\"UTF-8\""

// NewHeaders create the headers for sending a new message
func NewHeaders(date time.Time, from, to Address, replyTo *Address, subject, contentType string) *Headers {
	return &Headers{
		From:        from,
		To:          to,
		ReplyTo:     replyTo,
		Subject:     subject,
		Date:        date,
		ContentType: contentType,
	}
}

// Headers for the message
type Headers struct {
	From        Address
	To          Address
	Date        time.Time
	Subject     string
	ReplyTo     *Address
	ContentType string
}
