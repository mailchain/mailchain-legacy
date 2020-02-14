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

package bdbstore

import (
	"context"
	stderrors "errors"
	"fmt"
	"time"

	"github.com/dgraph-io/badger/v2"
	"github.com/mailchain/mailchain/internal/mail"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus" //nolint:depguard
)

const (
	// recommended discard ratio for the badgerDB GC.
	// ref: https://github.com/dgraph-io/badger/blob/master/db.go#L1107
	discardRatio = 0.5

	// interval at which the BadgerDB GC will be called.
	gcInterval = 10 * time.Minute
)

func newBadgerDB(opts *badger.Options) (*Database, error) {
	db, err := badger.Open(*opts)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	bdb := &Database{
		db:     db,
		ctx:    ctx,
		cancel: cancel,
	}

	go bdb.runGC()

	return bdb, nil
}

// Database is a persistent key-value store. Apart from
// basic data store functionality it also supports iterating
// over the key space in byte-wise lexicographical order,
// setting TTL on Keys and other functionality
// which can be found here: https://github.com/dgraph-io/badger
type Database struct {
	db *badger.DB

	// badgerDB GC
	ctx    context.Context
	cancel context.CancelFunc
}

// New returns a wrapped BadgerDB object with default options.
func New(dir string) (*Database, error) {
	opts := badger.DefaultOptions(dir)
	return newBadgerDB(&opts)
}

// NewWithOptions returns a wrapped BadgerDB object
// with the given options used.
func NewWithOptions(opts *badger.Options) (*Database, error) {
	return newBadgerDB(opts)
}

// Close flushes any pending updates to disk and closes
// the underlying key-value store.
func (db *Database) Close() error {
	db.cancel()
	return db.db.Close()
}

// PutMessageRead mark message as read.
func (db *Database) PutMessageRead(messageID mail.ID) error {
	return db.db.Update(func(txn *badger.Txn) error {
		return txn.Set(db.messageReadKey(messageID), []byte{1})
	})
}

// DeleteMessageRead mark message as unread.
func (db *Database) DeleteMessageRead(messageID mail.ID) error {
	return db.db.Update(func(txn *badger.Txn) error {
		return txn.Delete(db.messageReadKey(messageID))
	})
}

// GetReadStatus returns true if message is read.
func (db *Database) GetReadStatus(messageID mail.ID) (bool, error) {
	var (
		val []byte
		err error
	)

	err = db.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(db.messageReadKey(messageID))
		if err != nil {
			return err
		}

		// copy the value as the item is only valid
		// while the txn is open.
		val, err = item.ValueCopy(nil)
		return err
	})

	if err != nil {
		return false, err
	}

	if len(val) != 1 {
		return false, errors.Errorf("invalid read status length")
	}

	return val[0] == 1, nil
}

func (db *Database) messageReadKey(messageID mail.ID) []byte {
	return []byte(fmt.Sprintf("message.%s.read", messageID.HexString()))
}

// it implements a method Is(error) bool such that Is(target) returns true.
func Is(err, target error) bool { return stderrors.Is(err, target) }

func (db *Database) runGC() {
	ticker := time.NewTicker(gcInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := db.db.RunValueLogGC(discardRatio); err != nil {
				if Is(err, badger.ErrNoRewrite) {
					logrus.Debugf("BadgerDB GC call ended with no rewrites: %v", err)
				} else {
					logrus.Errorf("BadgerDB GC call failed: %v", err)
				}
			}
		case <-db.ctx.Done():
			return
		}
	}
}
