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
	"reflect"
	"testing"
)

func TestDecodeBase58(t *testing.T) {
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
				"5CLmNK8f16nagFeF2h3iNeeChaxPiAsJu7piNYJgdPpmaRzPD",
			},
			[]byte{0x9, 0x86, 0xc6, 0x71, 0x43, 0xad, 0x96, 0x6f, 0xa5, 0x79, 0xc9, 0x1b, 0x30, 0xc6, 0x7f, 0x95, 0xe7, 0x4b, 0xcc, 0xe3, 0xc5, 0xec, 0xb9, 0x5c, 0x96, 0xbf, 0xb5, 0x82, 0x87, 0x65, 0x64, 0xe4, 0x9c, 0x8, 0x3f, 0x1c},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeBase58(tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeBase58() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DecodeBase58() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncodeBase58(t *testing.T) {
	// assert := assert.New(t)
	type args struct {
		in []byte
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			"success",
			args{
				[]byte{0x9, 0x86, 0xc6, 0x71, 0x43, 0xad, 0x96, 0x6f, 0xa5, 0x79, 0xc9, 0x1b, 0x30, 0xc6, 0x7f, 0x95, 0xe7, 0x4b, 0xcc, 0xe3, 0xc5, 0xec, 0xb9, 0x5c, 0x96, 0xbf, 0xb5, 0x82, 0x87, 0x65, 0x64, 0xe4, 0x9c, 0x8, 0x3f, 0x1c},
			},
			"5CLmNK8f16nagFeF2h3iNeeChaxPiAsJu7piNYJgdPpmaRzPD",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EncodeBase58(tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncodeBase58() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EncodeBase58() = %v, want %v", got, tt.want)
			}
		})
	}
}
