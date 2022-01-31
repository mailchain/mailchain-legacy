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

package encoding

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_EncodeHexZeroX(t *testing.T) {
	type args struct {
		in []byte
	}
	tests := []struct {
		name        string
		args        args
		wantEncoded string
	}{
		{
			"success",
			args{
				[]byte{0x56, 0x2, 0xea, 0x95, 0x54, 0xb, 0xee, 0x46, 0xd0, 0x3b, 0xa3, 0x35, 0xee, 0xd6, 0xf4, 0x9d, 0x11, 0x7e, 0xab, 0x95, 0xc8, 0xab, 0x8b, 0x71, 0xba, 0xe2, 0xcd, 0xd1, 0xe5, 0x64, 0xa7, 0x61},
			},
			"0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotEncoded := EncodeHexZeroX(tt.args.in)
			if !assert.Equal(t, tt.wantEncoded, gotEncoded) {
				t.Errorf("EncodeZeroX() gotEncoded = %v, want %v", gotEncoded, tt.wantEncoded)
			}
		})
	}
}

func Test_DecodeHexZeroX(t *testing.T) {
	type args struct {
		in string
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
				"0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761",
			},
			[]byte{0x56, 0x2, 0xea, 0x95, 0x54, 0xb, 0xee, 0x46, 0xd0, 0x3b, 0xa3, 0x35, 0xee, 0xd6, 0xf4, 0x9d, 0x11, 0x7e, 0xab, 0x95, 0xc8, 0xab, 0x8b, 0x71, 0xba, 0xe2, 0xcd, 0xd1, 0xe5, 0x64, 0xa7, 0x61},
			false,
		},
		{
			"err-missing-prefix",
			args{
				"5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761",
			},
			nil,
			true,
		},
		{
			"err-empty",
			args{
				"",
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeHexZeroX(tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("decodeZeroX() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("decodeZeroX() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_EncodeHex(t *testing.T) {
	type args struct {
		in []byte
	}
	tests := []struct {
		name         string
		args         args
		wantEncoding string
	}{
		{
			"success",
			args{
				[]byte{0x56, 0x2, 0xea, 0x95, 0x54, 0xb, 0xee, 0x46, 0xd0, 0x3b, 0xa3, 0x35, 0xee, 0xd6, 0xf4, 0x9d, 0x11, 0x7e, 0xab, 0x95, 0xc8, 0xab, 0x8b, 0x71, 0xba, 0xe2, 0xcd, 0xd1, 0xe5, 0x64, 0xa7, 0x61},
			},
			"5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotEncoding := EncodeHex(tt.args.in)
			if !assert.Equal(t, tt.wantEncoding, gotEncoding) {
				t.Errorf("EncodeHex() gotEncoding = %v, want %v", gotEncoding, tt.wantEncoding)
			}
		})
	}
}

func Test_DecodeHex(t *testing.T) {
	type args struct {
		in string
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
				"5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761",
			},
			[]byte{86, 2, 234, 149, 84, 11, 238, 70, 208, 59, 163, 53, 238, 214, 244, 157, 17, 126, 171, 149, 200, 171, 139, 113, 186, 226, 205, 209, 229, 100, 167, 97},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeHex(tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeHex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DecodeHex() = %v, want %v", got, tt.want)
			}
		})
	}
}
