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

package scrypt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApplyDefault(t *testing.T) {
	assert := assert.New(t)
	opts := &DeriveOpts{}
	apply(opts, []DeriveOptionsBuilder{DefaultDeriveOptions()})
	assert.Equal(32, opts.Len)
	assert.Equal(262144, opts.N)
	assert.Equal(1, opts.P)
	assert.Equal("", opts.Passphrase)
	assert.Equal(8, opts.R)
	assert.Nil(opts.Salt)
}

func TestApplyDefaultAndPassword(t *testing.T) {
	assert := assert.New(t)
	randomSalt, err := RandomSalt()
	if err != nil {
		t.Fail()
	}

	opts := &DeriveOpts{}
	apply(opts, []DeriveOptionsBuilder{DefaultDeriveOptions(), WithPassphrase("test"), randomSalt})
	assert.Equal(32, opts.Len)
	assert.Equal(262144, opts.N)
	assert.Equal(1, opts.P)
	assert.Equal("test", opts.Passphrase)
	assert.Equal(8, opts.R)
	assert.Equal(32, len(opts.Salt))
}

func TestApplyFromEncryptedKey(t *testing.T) {
	assert := assert.New(t)

	opts := &DeriveOpts{}
	apply(opts, []DeriveOptionsBuilder{DefaultDeriveOptions(), FromEncryptedKey(32,
		1<<18,
		1,
		8,
		[]byte("salt-value"))})
	assert.Equal(32, opts.Len)
	assert.Equal(262144, opts.N)
	assert.Equal(1, opts.P)
	assert.Equal(8, opts.R)
	assert.Equal([]byte("salt-value"), opts.Salt)
}

func TestCreateOptions(t *testing.T) {
	assert := assert.New(t)
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
			if got := CreateOptions(tt.args.o); !assert.Equal(tt.want, got) {
				t.Errorf("CreateOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeriveKey(t *testing.T) {
	assert := assert.New(t)
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
			"success",
			args{
				[]DeriveOptionsBuilder{WithPassphrase("passphrase"), DefaultDeriveOptions()},
			},
			[]byte{0x14, 0x64, 0xbb, 0x68, 0x5d, 0xd1, 0x6a, 0xef, 0xe9, 0x1a, 0xe1, 0x34, 0x6d, 0x3d, 0x9a, 0x60, 0x2, 0x43, 0x1a, 0x5a, 0x1b, 0xc2, 0xde, 0xa6, 0x23, 0x35, 0xef, 0xc8, 0xad, 0x1f, 0x60, 0x1b},
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
			if !assert.Equal(got, tt.want) {
				t.Errorf("DeriveKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
