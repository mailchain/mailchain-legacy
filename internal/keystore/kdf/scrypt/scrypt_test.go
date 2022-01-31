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

package scrypt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApplyDefault(t *testing.T) {
	opts := &DeriveOpts{}
	apply(opts, []DeriveOptionsBuilder{DefaultDeriveOptions()})
	assert.Equal(t, 32, opts.Len)
	assert.Equal(t, 262144, opts.N)
	assert.Equal(t, 1, opts.P)
	assert.Equal(t, "", opts.Passphrase)
	assert.Equal(t, 8, opts.R)
	assert.Nil(t, opts.Salt)
}

func TestApplyDefaultAndPassword(t *testing.T) {
	randomSalt, err := RandomSalt()
	if err != nil {
		t.Fail()
	}

	opts := &DeriveOpts{}
	apply(opts, []DeriveOptionsBuilder{DefaultDeriveOptions(), WithPassphrase("test"), randomSalt})
	assert.Equal(t, 32, opts.Len)
	assert.Equal(t, 262144, opts.N)
	assert.Equal(t, 1, opts.P)
	assert.Equal(t, "test", opts.Passphrase)
	assert.Equal(t, 8, opts.R)
	assert.Equal(t, 32, len(opts.Salt))
}

func TestApplyFromEncryptedKey(t *testing.T) {

	opts := &DeriveOpts{}
	apply(opts, []DeriveOptionsBuilder{DefaultDeriveOptions(), FromEncryptedKey(32,
		1<<18,
		1,
		8,
		[]byte("salt-value"))})
	assert.Equal(t, 32, opts.Len)
	assert.Equal(t, 262144, opts.N)
	assert.Equal(t, 1, opts.P)
	assert.Equal(t, 8, opts.R)
	assert.Equal(t, []byte("salt-value"), opts.Salt)
}

func TestCreateOptions(t *testing.T) {
	type args struct {
		o []DeriveOptionsBuilder
	}
	tests := []struct {
		name string
		args args
		want *DeriveOpts
	}{
		{
			"success",
			args{
				[]DeriveOptionsBuilder{WithPassphrase("passphrase"), DefaultDeriveOptions()},
			},
			&DeriveOpts{Len: 32, N: 262144, P: 1, R: 8, Salt: []uint8(nil), Passphrase: "passphrase"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateOptions(tt.args.o); !assert.Equal(t, tt.want, got) {
				t.Errorf("CreateOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeriveKey(t *testing.T) {
	type args struct {
		o []DeriveOptionsBuilder
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			"success-passphrase",
			args{
				[]DeriveOptionsBuilder{WithPassphrase("passphrase"), DefaultDeriveOptions()},
			},
			[]byte{0x14, 0x64, 0xbb, 0x68, 0x5d, 0xd1, 0x6a, 0xef, 0xe9, 0x1a, 0xe1, 0x34, 0x6d, 0x3d, 0x9a, 0x60, 0x2, 0x43, 0x1a, 0x5a, 0x1b, 0xc2, 0xde, 0xa6, 0x23, 0x35, 0xef, 0xc8, 0xad, 0x1f, 0x60, 0x1b},
			false,
		},
		{
			"success-no-passphrase",
			args{
				[]DeriveOptionsBuilder{DefaultDeriveOptions()},
			},
			[]byte{0xae, 0x51, 0xe1, 0x6a, 0x25, 0x28, 0x52, 0x2f, 0x1c, 0x61, 0x85, 0xcc, 0xed, 0xb4, 0xba, 0xb, 0x96, 0xbc, 0x3, 0x32, 0x86, 0x17, 0x2, 0x65, 0x9d, 0x26, 0xf2, 0x6c, 0x54, 0x83, 0xc2, 0xe6},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DeriveKey(tt.args.o)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeriveKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, got, tt.want) {
				t.Errorf("DeriveKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
