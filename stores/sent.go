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

//go:generate mockgen -source=sent.go -package=mocks -destination=$PACKAGE_PATH/internal/testutil/mocks/sent.go

package stores

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/mailchain/mailchain/internal/crypto"
	"github.com/mailchain/mailchain/internal/mail"
	"github.com/pkg/errors"
)

// The Sent saves the message. This should not be used directly but as the first argument of storing.PutMessage.
type Sent interface {
	// PutMessage should write the message contents to the underlying storage service. Return the final location or any error.
	PutMessage(path string, msg io.Reader, headers map[string]string) (string, error)
}

// PutMessage does the pre work before saving the message as implemented by store.
func PutMessage(sent Sent, messageID mail.ID, msg io.Reader) (location string, err error) {
	if sent == nil {
		return "", errors.Errorf("'sent' must not be nil")
	}
	if msg == nil {
		return "", errors.Errorf("'msg' must not be nil")
	}
	contents, err := ioutil.ReadAll(msg)
	if err != nil {
		return "", err
	}
	hash := crypto.CreateLocationHash(contents)
	location, err = sent.PutMessage(fmt.Sprintf("%s-%s", messageID.HexString(), hash.String()), msg, nil)
	if err != nil {
		return "", errors.Wrap(err, "could not store message")
	}
	return location, nil
}
