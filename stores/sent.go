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

//go:generate mockgen -source=sent.go -package=storestest -destination=./storestest/sent_mock.go

package stores

import (
	"github.com/mailchain/mailchain/internal/mail"
)

// The Sent saves the message. This should not be used directly but as the first argument of storing.PutMessage.
type Sent interface {
	// PutMessage should write the message contents to the underlying storage service. Return the final location or any error.
	PutMessage(messageID mail.ID, msg []byte, headers map[string]string) (string, error)
	Key(messageID mail.ID, msg []byte) string
}

// func (s SentStore) PutMessage(path string, msg []byte, headers map[string]string) (string, error) {

// }
