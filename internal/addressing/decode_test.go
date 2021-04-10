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

package addressing

import (
	"testing"

	"github.com/mailchain/mailchain/encoding/encodingtest"
	"github.com/stretchr/testify/assert"
)

func TestDecodeByProtocol(t *testing.T) {
	type args struct {
		in       string
		protocol string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			"ethereum",
			args{
				"0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761",
				"ethereum",
			},
			encodingtest.MustDecodeHex("5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761"),
			false,
		},
		{
			"substrate",
			args{
				"5DJJhV3tVzsWG1jZfL157azn8iRyDC7HyNG1yh8v2nQYd994",
				"substrate",
			},
			encodingtest.MustDecodeBase58("5DJJhV3tVzsWG1jZfL157azn8iRyDC7HyNG1yh8v2nQYd994"),
			false,
		},
		{
			"algorand",
			args{
				"C7Z4NNMIMOGZW56JCILF6DVY4MBZJMHXUQ67W2WKVE6U5QJSIDPYUEAXQU",
				"algorand",
			},
			[]byte{0x17, 0xf3, 0xc6, 0xb5, 0x88, 0x63, 0x8d, 0x9b, 0x77, 0xc9, 0x12, 0x16, 0x5f, 0xe, 0xb8, 0xe3, 0x3, 0x94, 0xb0, 0xf7, 0xa4, 0x3d, 0xfb, 0x6a, 0xca, 0xa9, 0x3d, 0x4e, 0xc1, 0x32, 0x40, 0xdf, 0x8a, 0x10, 0x17, 0x85},
			false,
		},
		{
			"err",
			args{
				"0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761",
				"invalid",
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeByProtocol(tt.args.in, tt.args.protocol)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeByProtocol() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("DecodeByProtocol() = %v, want %v", got, tt.want)
			}
		})
	}
}
