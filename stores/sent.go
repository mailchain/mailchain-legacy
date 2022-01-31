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

package stores

import (
	"github.com/mailchain/mailchain/internal/mail"
)

//go:generate mockgen -source=sent.go -package=storestest -destination=./storestest/sent_mock.go

// The Sent saves the message. This should not be used directly but as the first argument of storing.PutMessage.
type Sent interface {
	// PutMessage should write the message contents to the underlying storage service. Return the final location information or any error.
	// Returns:
	// + address - url of the stored message. E.g. for https://mcx.mx/QmeEYqAaK7UpsffH9fDyPCwB9YH5jf36zDYgeqxhuCvHhh it is `https://mcx.mx/QmeEYqAaK7UpsffH9fDyPCwB9YH5jf36zDYgeqxhuCvHhh`
	// + resource - name of the stored message. E.g. for https://mcx.mx/QmeEYqAaK7UpsffH9fDyPCwB9YH5jf36zDYgeqxhuCvHhh it is `QmeEYqAaK7UpsffH9fDyPCwB9YH5jf36zDYgeqxhuCvHhh`
	// + mli - ID linked to the protocol + domain name, see `internal/envelope/mli.go`. E.g. for https://mcx.mx/QmeEYqAaK7UpsffH9fDyPCwB9YH5jf36zDYgeqxhuCvHhh is 1
	PutMessage(messageID mail.ID, contentsHash, msg []byte, headers map[string]string) (address, resource string, mli uint64, err error)
	Key(messageID mail.ID, contentsHash, msg []byte) string
}

type PutResponse struct {
	address           string
	resource          string
	mli               uint64
	contentsHashKind  int
	integrityHashKind int
}
