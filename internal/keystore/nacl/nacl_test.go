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

package nacl

import (
	"bytes"
	"io"
	"testing"
	"testing/iotest"

	"github.com/mailchain/mailchain/crypto/ed25519/ed25519test"
	"github.com/mailchain/mailchain/crypto/sr25519/sr25519test"
	"github.com/stretchr/testify/assert"
)

func Test_easySeal(t *testing.T) {
	type args struct {
		message []byte
		key     []byte
		rand    io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			"success-bob",
			args{
				[]byte("message"),
				ed25519test.BobPublicKey.Bytes(),
				bytes.NewReader([]byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")),
			},
			[]byte{0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4a, 0x4b, 0x4c, 0x4d, 0x4e, 0x4f, 0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x5b, 0x19, 0x83, 0xe5, 0x6e, 0x7f, 0xed, 0xfe, 0xbb, 0xd0, 0x70, 0x34, 0xce, 0x25, 0x49, 0x76, 0xa3, 0x50, 0x78, 0x91, 0x18, 0xe6, 0xe3},
			false,
		},
		{
			"success-alice",
			args{
				[]byte("message"),
				sr25519test.AlicePublicKey.Bytes(),
				bytes.NewReader([]byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")),
			},
			[]byte{0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4a, 0x4b, 0x4c, 0x4d, 0x4e, 0x4f, 0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0xa3, 0x81, 0x27, 0x6f, 0xdb, 0x97, 0x17, 0x51, 0x10, 0x1e, 0x17, 0x2c, 0xec, 0x5b, 0xae, 0xdc, 0x7, 0x26, 0xea, 0x16, 0xe4, 0xc7, 0xde},
			false,
		},
		{
			"err-rand",
			args{
				[]byte("message"),
				ed25519test.BobPublicKey.Bytes(),
				iotest.DataErrReader(bytes.NewReader(nil)),
			},
			nil,
			true,
		},
		{
			"err-rand-schnorrkel",
			args{
				[]byte("message"),
				sr25519test.AlicePublicKey.Bytes(),
				iotest.DataErrReader(bytes.NewReader(nil)),
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := easySeal(tt.args.message, tt.args.key, tt.args.rand)
			if (err != nil) != tt.wantErr {
				t.Errorf("easySeal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("easySeal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_easyOpen(t *testing.T) {

	type args struct {
		box []byte
		key []byte
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
				[]byte{0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4a, 0x4b, 0x4c, 0x4d, 0x4e, 0x4f, 0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x5b, 0x19, 0x83, 0xe5, 0x6e, 0x7f, 0xed, 0xfe, 0xbb, 0xd0, 0x70, 0x34, 0xce, 0x25, 0x49, 0x76, 0xa3, 0x50, 0x78, 0x91, 0x18, 0xe6, 0xe3},
				ed25519test.BobPrivateKey.Bytes()[32:],
			},
			[]byte("message"),
			false,
		},
		{
			"err-key-size",
			args{
				[]byte{0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4a, 0x4b, 0x4c, 0x4d, 0x4e, 0x4f, 0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x5b, 0x19, 0x83, 0xe5, 0x6e, 0x7f, 0xed, 0xfe, 0xbb, 0xd0, 0x70, 0x34, 0xce, 0x25, 0x49, 0x76, 0xa3, 0x50, 0x78, 0x91, 0x18, 0xe6, 0xe3},
				ed25519test.BobPrivateKey.Bytes()[16:],
			},
			nil,
			true,
		},
		{
			"err-nonce-size",
			args{
				[]byte{0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47},
				ed25519test.BobPrivateKey.Bytes(),
			},
			nil,
			true,
		},
		{
			"err-input",
			args{
				[]byte{0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4a, 0x4b, 0x4c, 0x4d, 0x4e, 0x4f, 0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x5b, 0x19, 0x83, 0xe5, 0x6e, 0x7f, 0xed, 0xfe, 0xbb, 0xd0, 0x70, 0x34, 0xce, 0x25, 0x49, 0x76, 0xa3, 0x50, 0x78, 0x91, 0x18, 0xe6, 0xe3},
				ed25519test.AlicePrivateKey.Bytes(),
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := easyOpen(tt.args.box, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("easyOpen() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("easyOpen() = %v, want %v", got, tt.want)
			}
		})
	}
}
