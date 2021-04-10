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
	"io"
	"log"
	"time"

	"github.com/dgraph-io/badger/v2"
	"github.com/pkg/errors"
	"github.com/dgraph-io/badger/v2/options"

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
		return nil, errors.WithStack(err)
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
//
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
func New(dir string, logWriter io.Writer) (*Database, error) {
	opts := badger.DefaultOptions(dir)
	opts.Logger = newLogger(logWriter)
	opts.Truncate = true
	opts.ValueLogLoadingMode = options.FileIO
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

func (db *Database) runGC() {
	ticker := time.NewTicker(gcInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := db.db.RunValueLogGC(discardRatio); err != nil {
				if errors.Is(errors.Cause(err), badger.ErrNoRewrite) {
					log.Printf("BadgerDB GC call ended with no rewrites: %v\n", err)
				} else {
					log.Printf("BadgerDB GC call failed: %v\n", err)
				}
			}
		case <-db.ctx.Done():
			return
		}
	}
}
