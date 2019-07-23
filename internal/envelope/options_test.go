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

package envelope

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithKind(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		kind byte
	}
	tests := []struct {
		name string
		args args
		want *CreateOpts
	}{
		{
			"success",
			args{0x01},
			&CreateOpts{Kind: 0x1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WithKind(tt.args.kind)
			opts := &CreateOpts{}
			got(opts)
			if !assert.Equal(tt.want, opts) {
				t.Errorf("WithKind() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithURL(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		address string
	}
	tests := []struct {
		name string
		args args
		want *CreateOpts
	}{
		{
			"success",
			args{"https://address.com"},
			&CreateOpts{URL: "https://address.com"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WithURL(tt.args.address)
			opts := &CreateOpts{}
			got(opts)
			if !assert.Equal(tt.want, opts) {
				t.Errorf("WithKind() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithResource(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		resource string
	}
	tests := []struct {
		name string
		args args
		want *CreateOpts
	}{
		{
			"success",
			args{"resource"},
			&CreateOpts{Resource: "resource"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WithResource(tt.args.resource)
			opts := &CreateOpts{}
			got(opts)
			if !assert.Equal(tt.want, opts) {
				t.Errorf("WithKind() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithMessageLocationIdentifier(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		mli uint64
	}
	tests := []struct {
		name    string
		args    args
		want    *CreateOpts
		wantErr bool
	}{
		{
			"success-0",
			args{0},
			&CreateOpts{Location: 0},
			false,
		},
		{
			"success-1",
			args{1},
			&CreateOpts{Location: 1},
			false,
		},
		{
			"invalid",
			args{99999999},
			&CreateOpts{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := WithMessageLocationIdentifier(tt.args.mli)
			if (err != nil) != tt.wantErr {
				t.Errorf("WithMessageLocationIdentifier() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			opts := &CreateOpts{}
			got(opts)
			if !assert.Equal(tt.want, opts) {
				t.Errorf("WithKind() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithDecryptedHash(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		decryptedHash []byte
	}
	tests := []struct {
		name string
		args args
		want *CreateOpts
	}{
		{
			"success",
			args{[]byte("decrypted-hash")},
			&CreateOpts{DecryptedHash: []byte("decrypted-hash")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WithDecryptedHash(tt.args.decryptedHash)
			opts := &CreateOpts{}
			got(opts)
			if !assert.Equal(tt.want, opts) {
				t.Errorf("WithKind() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithEncryptedHash(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		encryptedHash []byte
	}
	tests := []struct {
		name string
		args args
		want *CreateOpts
	}{
		{
			"success",
			args{[]byte("encrypted-hash")},
			&CreateOpts{EncryptedHash: []byte("encrypted-hash")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WithEncryptedHash(tt.args.encryptedHash)
			opts := &CreateOpts{}
			got(opts)
			if !assert.Equal(tt.want, opts) {
				t.Errorf("WithEncryptedHash() = %v, want %v", got, tt.want)
			}
		})
	}
}
