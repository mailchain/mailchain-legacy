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

package stores

import (
	"fmt"

	"github.com/mailchain/mailchain/internal/pkg/crypto"
	"github.com/mailchain/mailchain/internal/pkg/mail"
	"github.com/pkg/errors"
)

// The Sent saves the message. This should not be used directly but as the first argument of storing.PutMessage.
type Sent interface {
	// PutMessage should write the message contents to the underlying storage service. Return the final location or any error.
	PutMessage(path string, msg []byte) (location string, err error)
}

// PutMessage does the pre work before saving the message as implemented by store.
func PutMessage(sent Sent, messageID mail.ID, msg []byte) (location string, err error) {
	hash, err := crypto.CreateLocationHash(msg)
	if err != nil {
		return "", err
	}
	location, err = sent.PutMessage(fmt.Sprintf("%s-%s", messageID.HexString(), hash.String()), msg)
	if err != nil {
		return "", errors.Wrap(err, "could not store message")
	}
	return location, nil
}
