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

package ldbstore

import (
	"fmt"

	"github.com/mailchain/mailchain/internal/mail"
	"github.com/pkg/errors"
)

func messageReadKey(messageID mail.ID) []byte {
	return []byte(fmt.Sprintf("message.%s.read", messageID.HexString()))
}

// DeleteMessageRead mark message as unread
func (db Database) DeleteMessageRead(messageID mail.ID) error {
	return db.db.Delete(messageReadKey(messageID), nil)
}

// PutMessageRead mark message as read
func (db Database) PutMessageRead(messageID mail.ID) error {
	return db.db.Put(messageReadKey(messageID), []byte{1}, nil)
}

// GetReadStatus return if message is read
func (db Database) GetReadStatus(messageID mail.ID) (bool, error) {
	value, err := db.db.Get(messageReadKey(messageID), nil)
	if err != nil {
		return false, err
	}
	if len(value) != 1 {
		return false, errors.Errorf("invalid read status length")
	}
	return value[0] == 1, nil
}
