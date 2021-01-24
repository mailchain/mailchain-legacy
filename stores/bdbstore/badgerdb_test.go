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

// +build !windows

package bdbstore

import (
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"testing"
)

type kv struct {
	key []byte
	val []byte
}

func getTempDir() string {
	switch runtime.GOOS {
	case "windows":
		panic("getTempDir() unsupported for windows")
	default:
		return "/tmp"
	}
}

func setupDB(path string) (*Database, func(), error) {
	db, err := New(path)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating BadgerDB: %v", err)
	}

	teardown := func() {
		if err := db.Close(); err != nil {
			log.Fatalf("error closing BadgerDB: %v", err)
		}

		if err := os.RemoveAll(path); err != nil {
			log.Fatalf("error removing test directories: %v", err)
		}
	}

	w := db.db.NewWriteBatch()
	defer w.Cancel()

	for _, kv := range []kv{
		{
			[]byte{0x47, 0xec, 0xa0, 0x11, 0xe3, 0x2b, 0x52, 0xc7, 0x10, 0x5, 0xad, 0x8a, 0x8f, 0x75, 0xe1, 0xb4, 0x4c, 0x92, 0xc9, 0x9f, 0xd1, 0x2e, 0x43, 0xbc, 0xcf, 0xe5, 0x71, 0xe3, 0xc2, 0xd1, 0x3d, 0x2e, 0x9a, 0x82, 0x6a, 0x55, 0xf, 0x5f, 0xf6, 0x3b, 0x24, 0x7a, 0xf4, 0x71},
			[]byte{1},
		},
		{
			[]byte{0x48, 0xed, 0xaf, 0x12, 0xe3, 0x2b, 0x52, 0xc7, 0x10, 0x5, 0xad, 0x8a, 0x8f, 0x75, 0xe1, 0xb4, 0x4c, 0x92, 0xc9, 0x9f, 0xd1, 0x2e, 0x43, 0xbc, 0xcf, 0xe5, 0x71, 0xe3, 0xc2, 0xd1, 0x3d, 0x2e, 0x9a, 0x82, 0x6a, 0x55, 0xf, 0x5f, 0xf6, 0x3b, 0x24, 0x7a, 0xf4, 0x71},
			[]byte{1},
		},
		{
			[]byte{0x78, 0xed, 0xaf, 0x12, 0xe3, 0x2b, 0x52, 0xc7, 0x10, 0x5, 0xad, 0x8a, 0x8f, 0x75, 0xe1, 0xb4, 0x4c, 0x92, 0xc9, 0x9f, 0xd1, 0x2e, 0x43, 0xbc, 0xcf, 0xe5, 0x71, 0xe3, 0xc2, 0xd1, 0x3d, 0x2e, 0x9a, 0x82, 0x6a, 0x55, 0xf, 0x5f, 0xf6, 0x3b, 0x24, 0x7a, 0xf4, 0x71},
			[]byte{1, 2},
		},
		{
			nil,
			[]byte{1},
		},
	} {
		if err := w.Set(db.messageReadKey(kv.key), kv.val); err != nil {
			return nil, nil, fmt.Errorf("error preparing BadgerDB: %v", err)
		}
	}

	if err := w.Flush(); err != nil {
		return nil, nil, fmt.Errorf("error preparing BadgerDB: %v", err)
	}

	return db, teardown, nil
}

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		dir     string
		wantErr bool
	}{
		{"success-new-1", path.Join(getTempDir(), "bdb_test_new_1"), false},
		{"success-new-2", path.Join(getTempDir(), "bdb_test_new_2"), false},
		{"err-new-3", path.Join(getTempDir(), "uncreated_test_dir", "bdb_test_new_3"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := New(tt.dir)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() failed to open badger at path %q, err: %v", tt.dir, err)
				return
			}

			if db != nil {
				if err := db.Close(); err != nil {
					t.Errorf("error closing BadgerDB: %v", err)
				}

				if err := os.RemoveAll(tt.dir); err != nil {
					log.Fatalf("error removing test directories: %v", err)
				}
			}
		})
	}
}
