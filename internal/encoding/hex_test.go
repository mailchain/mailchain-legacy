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

package encoding

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/mailchain/mailchain/internal/encoding/encodingtest"
	"github.com/stretchr/testify/assert"
)

type encDecTest struct {
	enc string
	dec []byte
}

var encDecTests = []encDecTest{
	{"", []byte{}},
	{"0001020304050607", []byte{0, 1, 2, 3, 4, 5, 6, 7}},
	{"08090a0b0c0d0e0f", []byte{8, 9, 10, 11, 12, 13, 14, 15}},
	{"f0f1f2f3f4f5f6f7", []byte{0xf0, 0xf1, 0xf2, 0xf3, 0xf4, 0xf5, 0xf6, 0xf7}},
	{"f8f9fafbfcfdfeff", []byte{0xf8, 0xf9, 0xfa, 0xfb, 0xfc, 0xfd, 0xfe, 0xff}},
	{"67", []byte{'g'}},
	{"e3a1", []byte{0xe3, 0xa1}},
}

func Test_EncodeZeroX(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		in []byte
	}
	tests := []struct {
		name         string
		args         args
		wantEncoded  []byte
		wantEncoding string
	}{
		{
			"success",
			args{
				encodingtest.MustDecodeHex("5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761"),
			},
			[]byte{0x56, 0x2, 0xea, 0x95, 0x54, 0xb, 0xee, 0x46, 0xd0, 0x3b, 0xa3, 0x35, 0xee, 0xd6, 0xf4, 0x9d, 0x11, 0x7e, 0xab, 0x95, 0xc8, 0xab, 0x8b, 0x71, 0xba, 0xe2, 0xcd, 0xd1, 0xe5, 0x64, 0xa7, 0x61},
			TypeHex0XPrefix,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotEncoded, gotEncoding := EncodeZeroX(tt.args.in)
			if !assert.Equal(tt.wantEncoded, gotEncoded) {
				t.Errorf("EncodeZeroX() gotEncoded = %v, want %v", gotEncoded, tt.wantEncoded)
			}
			if !assert.Equal(tt.wantEncoding, gotEncoding) {
				t.Errorf("EncodeZeroX() gotEncoding = %v, want %v", gotEncoding, tt.wantEncoding)
			}
		})
	}
}

func Test_DecodeZeroX(t *testing.T) {
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
			got, err := DecodeZeroX(tt.args.in)
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
	for i, test := range encDecTests {
		s := EncodeHex(test.dec)
		if s != test.enc {
			t.Errorf("#%d got:%s want:%s", i, s, test.enc)
		}
	}

}

func Test_DecodeHex(t *testing.T) {
	for i, test := range encDecTests {
		dst, err := DecodeHex(test.enc)
		if err != nil {
			t.Errorf("#%d: unexpected err value: %s", i, err)
			continue
		}
		if !bytes.Equal(dst, test.dec) {
			t.Errorf("#%d: got: %#v want: #%v", i, dst, test.dec)
		}
	}
}
