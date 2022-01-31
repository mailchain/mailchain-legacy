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

package envelope

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshal(t *testing.T) {
	type args struct {
		data Data
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			"ZeroX01",
			args{
				&ZeroX01{
					UIBEncryptedLocationHash: []byte{0x10, 0xe1, 0xd3},
				},
			},
			[]byte{0x1, 0xa, 0x3, 0x10, 0xe1, 0xd3},
			false,
		},
		{
			"ZeroX02",
			args{
				&ZeroX02{
					UIBEncryptedLocationHash: []byte{0x10, 0xe1, 0xd3},
					DecryptedHash:            []byte{0x10, 0x11, 0x12},
				},
			},
			[]byte{0x2, 0xa, 0x3, 0x10, 0xe1, 0xd3, 0x12, 0x3, 0x10, 0x11, 0x12},
			false,
		},
		{
			"err-unknown",
			args{
				&ZeroX50{},
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Marshal(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("Marshal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnmarshal(t *testing.T) {
	type args struct {
		buf []byte
	}
	tests := []struct {
		name    string
		args    args
		want    Data
		wantErr bool
	}{
		{
			"success-0x01",
			args{
				[]byte{0x1, 0xa, 0x3, 0x10, 0xe1, 0xd3},
			},
			&ZeroX01{
				UIBEncryptedLocationHash: []byte{0x10, 0xe1, 0xd3},
			},
			false,
		},
		{
			"success-0x02",
			args{
				[]byte{0x2, 0xa, 0x3, 0x10, 0xe1, 0xd3, 0x12, 0x3, 0x10, 0x11, 0x12},
			},
			&ZeroX02{
				UIBEncryptedLocationHash: []byte{0x10, 0xe1, 0xd3},
				DecryptedHash:            []byte{0x10, 0x11, 0x12},
			},
			false,
		},
		{
			"success-0x50",
			args{
				[]byte{0x50, 0x12, 0xd, 0x65, 0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x65, 0x64, 0x2d, 0x75, 0x72, 0x6c, 0x1a, 0xe, 0x64, 0x65, 0x63, 0x72, 0x79, 0x70, 0x74, 0x65, 0x64, 0x2d, 0x68, 0x61, 0x73, 0x68},
			},
			&ZeroX50{
				Version:       0,
				EncryptedURL:  []uint8{0x65, 0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x65, 0x64, 0x2d, 0x75, 0x72, 0x6c},
				DecryptedHash: []uint8{0x64, 0x65, 0x63, 0x72, 0x79, 0x70, 0x74, 0x65, 0x64, 0x2d, 0x68, 0x61, 0x73, 0x68},
			},
			false,
		},
		{
			"err-invalid",
			args{
				[]byte{0x0, 0xa, 0x3, 0x10, 0xe1, 0xd3},
			},
			nil,
			true,
		},
		{
			"err-empty",
			args{
				[]byte{},
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Unmarshal(tt.args.buf)
			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("Unmarshal() = %v, want %v", got, tt.want)
			}
		})
	}
}
