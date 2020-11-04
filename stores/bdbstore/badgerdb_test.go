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
	"time"

	"github.com/mailchain/mailchain/internal/mail"
	"github.com/mailchain/mailchain/stores"

	"github.com/stretchr/testify/assert"
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

func TestDatabase_GetReadStatus(t *testing.T) {
	db, teardown, _ := setupDB(path.Join(getTempDir(), "bdb_test_read_status"))
	defer teardown()

	type args struct {
		messageID mail.ID
	}
	tests := []struct {
		name           string
		args           args
		wantStatusRead bool
		wantErr        bool
	}{
		{
			"success-key-0",
			args{
				[]byte{0x47, 0xec, 0xa0, 0x11, 0xe3, 0x2b, 0x52, 0xc7, 0x10, 0x5, 0xad, 0x8a, 0x8f, 0x75, 0xe1, 0xb4, 0x4c, 0x92, 0xc9, 0x9f, 0xd1, 0x2e, 0x43, 0xbc, 0xcf, 0xe5, 0x71, 0xe3, 0xc2, 0xd1, 0x3d, 0x2e, 0x9a, 0x82, 0x6a, 0x55, 0xf, 0x5f, 0xf6, 0x3b, 0x24, 0x7a, 0xf4, 0x71},
			},
			true,
			false,
		},
		{
			"success-key-1",
			args{
				[]byte{0x48, 0xed, 0xaf, 0x12, 0xe3, 0x2b, 0x52, 0xc7, 0x10, 0x5, 0xad, 0x8a, 0x8f, 0x75, 0xe1, 0xb4, 0x4c, 0x92, 0xc9, 0x9f, 0xd1, 0x2e, 0x43, 0xbc, 0xcf, 0xe5, 0x71, 0xe3, 0xc2, 0xd1, 0x3d, 0x2e, 0x9a, 0x82, 0x6a, 0x55, 0xf, 0x5f, 0xf6, 0x3b, 0x24, 0x7a, 0xf4, 0x71},
			},
			true,
			false,
		},
		{
			"err-invalid-key",
			args{
				[]byte{0x49, 0xef, 0xaf, 0x12, 0xe3, 0x2b, 0x52, 0xc7, 0x10, 0x5, 0xad, 0x8a, 0x8f, 0x75, 0xe1, 0xb4, 0x4c, 0x92, 0xc9, 0x9f, 0xd1, 0x2e, 0x43, 0xbc, 0xcf, 0xe5, 0x71, 0xe3, 0xc2, 0xd1, 0x3d, 0x2e, 0x9a, 0x82, 0x6a, 0x55, 0xf, 0x5f, 0xf6, 0x3b, 0x24, 0x7a, 0xf4, 0x71},
			},
			false,
			true,
		},
		{
			"err-invalid-value-length",
			args{
				[]byte{0x78, 0xed, 0xaf, 0x12, 0xe3, 0x2b, 0x52, 0xc7, 0x10, 0x5, 0xad, 0x8a, 0x8f, 0x75, 0xe1, 0xb4, 0x4c, 0x92, 0xc9, 0x9f, 0xd1, 0x2e, 0x43, 0xbc, 0xcf, 0xe5, 0x71, 0xe3, 0xc2, 0xd1, 0x3d, 0x2e, 0x9a, 0x82, 0x6a, 0x55, 0xf, 0x5f, 0xf6, 0x3b, 0x24, 0x7a, 0xf4, 0x71},
			},
			false,
			true,
		},
		{
			"success-key-3",
			args{
				nil,
			},
			true,
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, err := db.GetReadStatus(tt.args.messageID)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetReadStatus() err = %v, want err %v", err, tt.wantErr)
			}

			if status != tt.wantStatusRead {
				t.Errorf("GetReadStatus() status = %v, want %v", status, tt.wantStatusRead)
			}
		})
	}
}

func TestDatabase_DeleteMessageRead(t *testing.T) {
	db, teardown, _ := setupDB(path.Join(getTempDir(), "bdb_test_delete_message"))
	defer teardown()

	type args struct {
		messageID mail.ID
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"success-key-0",
			args{
				[]byte{0x47, 0xec, 0xa0, 0x11, 0xe3, 0x2b, 0x52, 0xc7, 0x10, 0x5, 0xad, 0x8a, 0x8f, 0x75, 0xe1, 0xb4, 0x4c, 0x92, 0xc9, 0x9f, 0xd1, 0x2e, 0x43, 0xbc, 0xcf, 0xe5, 0x71, 0xe3, 0xc2, 0xd1, 0x3d, 0x2e, 0x9a, 0x82, 0x6a, 0x55, 0xf, 0x5f, 0xf6, 0x3b, 0x24, 0x7a, 0xf4, 0x71},
			},
			false,
		},
		{
			"success-key-1",
			args{
				[]byte{0x48, 0xed, 0xaf, 0x12, 0xe3, 0x2b, 0x52, 0xc7, 0x10, 0x5, 0xad, 0x8a, 0x8f, 0x75, 0xe1, 0xb4, 0x4c, 0x92, 0xc9, 0x9f, 0xd1, 0x2e, 0x43, 0xbc, 0xcf, 0xe5, 0x71, 0xe3, 0xc2, 0xd1, 0x3d, 0x2e, 0x9a, 0x82, 0x6a, 0x55, 0xf, 0x5f, 0xf6, 0x3b, 0x24, 0x7a, 0xf4, 0x71},
			},
			false,
		},
		{
			"success-key-2",
			args{
				[]byte{0x49, 0xef, 0xaf, 0x12, 0xe3, 0x2b, 0x52, 0xc7, 0x10, 0x5, 0xad, 0x8a, 0x8f, 0x75, 0xe1, 0xb4, 0x4c, 0x92, 0xc9, 0x9f, 0xd1, 0x2e, 0x43, 0xbc, 0xcf, 0xe5, 0x71, 0xe3, 0xc2, 0xd1, 0x3d, 0x2e, 0x9a, 0x82, 0x6a, 0x55, 0xf, 0x5f, 0xf6, 0x3b, 0x24, 0x7a, 0xf4, 0x71},
			},
			false,
		},
		{
			"success-key-3",
			args{
				nil,
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := db.DeleteMessageRead(tt.args.messageID); (err != nil) != tt.wantErr {
				t.Errorf("DeleteMessageRead() err = %v, want err %v", err, tt.wantErr)
			}
		})
	}
}

func TestDatabase_PutMessageRead(t *testing.T) {
	db, teardown, _ := setupDB(path.Join(getTempDir(), "bdb_test_put_message_read"))
	defer teardown()

	type args struct {
		messageID mail.ID
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"success-key-0",
			args{
				[]byte{0x57, 0xac, 0xa0, 0x11, 0xe3, 0x2b, 0x52, 0xc7, 0x10, 0x5, 0xad, 0x8a, 0x8f, 0x75, 0xe1, 0xb4, 0x4c, 0x92, 0xc9, 0x9f, 0xd1, 0x2e, 0x43, 0xbc, 0xcf, 0xe5, 0x71, 0xe3, 0xc2, 0xd1, 0x3d, 0x2e, 0x9a, 0x82, 0x6a, 0x55, 0xf, 0x5f, 0xf6, 0x3b, 0x24, 0x7a, 0xf4, 0x71},
			},
			false,
		},
		{
			"success-key-1",
			args{
				[]byte{0x7a, 0xad, 0xff, 0x12, 0xe3, 0x2b, 0x52, 0xc7, 0x10, 0x5, 0xad, 0x8a, 0x8f, 0x75, 0xe1, 0xb4, 0x4c, 0x92, 0xc9, 0x9f, 0xd1, 0x2e, 0x43, 0xbc, 0xcf, 0xe5, 0x71, 0xe3, 0xc2, 0xd1, 0x3d, 0x2e, 0x9a, 0x82, 0x6a, 0x55, 0xf, 0x5f, 0xf6, 0x3b, 0x24, 0x7a, 0xf4, 0x71},
			},
			false,
		},
		{
			"success-key-2",
			args{
				[]byte{0x29, 0x1f, 0xef, 0xff, 0xe3, 0x2b, 0x52, 0xc7, 0x10, 0x5, 0xad, 0x8a, 0x8f, 0x75, 0xe1, 0xb4, 0x4c, 0x92, 0xc9, 0x9f, 0xd1, 0x2e, 0x43, 0xbc, 0xcf, 0xe5, 0x71, 0xe3, 0xc2, 0xd1, 0x3d, 0x2e, 0x9a, 0x82, 0x6a, 0x55, 0xf, 0x5f, 0xf6, 0x3b, 0x24, 0x7a, 0xf4, 0x71},
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := db.PutMessageRead(tt.args.messageID); (err != nil) != tt.wantErr {
				t.Errorf("PutMessageRead() err = %v, want err %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetMessages(t *testing.T) {
	db, teardown, _ := setupDB(path.Join(getTempDir(), "bdb_test_get_messages"))
	defer teardown()

	type args struct {
		protocol string
		network  string
		address  string
		message  stores.Message
	}

	tests := []struct {
		name            string
		args            []args
		expectedIDOrder []string
		queryProtocol   string
		queryNetwork    string
		queryAddress    string
		wantErr         bool
	}{
		{
			name: "successfully get messages",
			args: []args{
				{
					protocol: "test",
					network:  "test",
					address:  "test",
					message: stores.Message{
						Headers: stores.Header{
							Date:        time.Now().Add(-1 * time.Hour),
							From:        "",
							To:          "",
							ReplyTo:     "",
							MessageID:   "0xc0302c62ced73abdf1d034553fe059bb596f2ef2",
							ContentType: "",
						},
						Body: "test1234",
					},
				},
				{
					protocol: "test",
					network:  "test",
					address:  "test",
					message: stores.Message{
						Headers: stores.Header{
							Date:        time.Now().Add(1 * time.Hour),
							From:        "",
							To:          "",
							ReplyTo:     "",
							MessageID:   "0x92d8f10248c6a3953cc3692a894655ad05d61efb",
							ContentType: "",
						},
						Body: "test1234",
					},
				},
				{
					protocol: "test",
					network:  "test",
					address:  "test",
					message: stores.Message{
						Headers: stores.Header{
							Date:        time.Now(),
							From:        "",
							To:          "",
							ReplyTo:     "",
							MessageID:   "0x5602ea95540bee46d03ba335eed6f49d117eab95c",
							ContentType: "",
						},
						Body: "test1234",
					},
				},
				{
					protocol: "test",
					network:  "test",
					address:  "test",
					message: stores.Message{
						Headers: stores.Header{
							Date:        time.Now().Add(-2 * time.Hour),
							From:        "",
							To:          "",
							ReplyTo:     "",
							MessageID:   "0xd2c574543459bf6704174fa869df4974220b71f673",
							ContentType: "",
						},
						Body: "test1234",
					},
				},
				{
					protocol: "test1",
					network:  "test2",
					address:  "test",
					message: stores.Message{
						Headers: stores.Header{
							Date:        time.Now().Add(-2 * time.Hour),
							From:        "",
							To:          "",
							ReplyTo:     "",
							MessageID:   "0x134574543459bf6704174fa869df4974220b123",
							ContentType: "",
						},
						Body: "test1234",
					},
				},
			},
			queryAddress:  "test",
			queryNetwork:  "test",
			queryProtocol: "test",
			// Higher to lower
			expectedIDOrder: []string{"0x92d8f10248c6a3953cc3692a894655ad05d61efb", "0x5602ea95540bee46d03ba335eed6f49d117eab95c", "0xc0302c62ced73abdf1d034553fe059bb596f2ef2", "0xd2c574543459bf6704174fa869df4974220b71f673"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, args := range tt.args {
				if err := db.PutMessage(args.protocol, args.network, args.address, args.message); err != nil {
					t.Errorf("PutMessage() returned an error %v", err)
				}
			}

			messages, err := db.GetMessages(tt.queryProtocol, tt.queryNetwork, tt.queryAddress)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMessages() err = %v, want err %v", err, tt.wantErr)
				return
			}

			var IDOrder []string
			for _, message := range messages {
				IDOrder = append(IDOrder, message.Headers.MessageID)
			}

			assert.Equal(t, tt.expectedIDOrder, IDOrder, "expected ID order %v, but got %v", tt.expectedIDOrder, IDOrder)
		})
	}
}

func TestPutMessage(t *testing.T) {
	db, teardown, _ := setupDB(path.Join(getTempDir(), "bdb_test_put_message"))
	defer teardown()

	type args struct {
		protocol string
		network  string
		address  string
		message  stores.Message
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "successful insert 1",
			args: args{
				protocol: "test",
				network:  "test",
				address:  "test",
				message: stores.Message{Headers: stores.Header{
					Date:        time.Now(),
					From:        "",
					To:          "",
					ReplyTo:     "",
					MessageID:   "0xc0302c62ced73abdf1d034553fe059bb596f2ef2",
					ContentType: "",
				},
					Body: "test1234",
				},
			},
		},
		{
			name: "successful insert 1",
			args: args{
				protocol: "test",
				network:  "test",
				address:  "test",
				message: stores.Message{Headers: stores.Header{
					Date:      time.Now(),
					MessageID: "0x92d8f10248c6a3953cc3692a894655ad05d61efb",
				},
					Body: "test1234",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := db.PutMessage(tt.args.protocol, tt.args.network, tt.args.address, tt.args.message); (err != nil) != tt.wantErr {
				t.Errorf("PutMessage() err = %v, want err %v", err, tt.wantErr)
			}
		})
	}
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
