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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeBase32(t *testing.T) {
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
				"BGDMM4KDVWLG7JLZZENTBRT7SXTUXTHDYXWLSXEWX62YFB3FMTSJYCB7DQ",
			},
			[]byte{0x9, 0x86, 0xc6, 0x71, 0x43, 0xad, 0x96, 0x6f, 0xa5, 0x79, 0xc9, 0x1b, 0x30, 0xc6, 0x7f, 0x95, 0xe7, 0x4b, 0xcc, 0xe3, 0xc5, 0xec, 0xb9, 0x5c, 0x96, 0xbf, 0xb5, 0x82, 0x87, 0x65, 0x64, 0xe4, 0x9c, 0x8, 0x3f, 0x1c},
			false,
		},
		{
			"algorand-address",
			args{
				"C7Z4NNMIMOGZW56JCILF6DVY4MBZJMHXUQ67W2WKVE6U5QJSIDPYUEAXQU",
			},
			[]byte{0x17, 0xf3, 0xc6, 0xb5, 0x88, 0x63, 0x8d, 0x9b, 0x77, 0xc9, 0x12, 0x16, 0x5f, 0xe, 0xb8, 0xe3, 0x3, 0x94, 0xb0, 0xf7, 0xa4, 0x3d, 0xfb, 0x6a, 0xca, 0xa9, 0x3d, 0x4e, 0xc1, 0x32, 0x40, 0xdf, 0x8a, 0x10, 0x17, 0x85},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeBase32(tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeBase32() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("DecodeBase32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncodeBase32(t *testing.T) {
	type args struct {
		in []byte
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"success",
			args{
				[]byte{0x9, 0x86, 0xc6, 0x71, 0x43, 0xad, 0x96, 0x6f, 0xa5, 0x79, 0xc9, 0x1b, 0x30, 0xc6, 0x7f, 0x95, 0xe7, 0x4b, 0xcc, 0xe3, 0xc5, 0xec, 0xb9, 0x5c, 0x96, 0xbf, 0xb5, 0x82, 0x87, 0x65, 0x64, 0xe4, 0x9c, 0x8, 0x3f, 0x1c},
			},
			"BGDMM4KDVWLG7JLZZENTBRT7SXTUXTHDYXWLSXEWX62YFB3FMTSJYCB7DQ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EncodeBase32(tt.args.in)
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("EncodeBase32() = %v, want %v", got, tt.want)
			}
		})
	}
}
