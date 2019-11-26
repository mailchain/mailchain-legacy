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

	"github.com/mailchain/mailchain/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func Test_EncodeZeroX(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		in []byte
	}
	tests := []struct {
		name         string
		args         args
		wantEncoded  string
	}{
		{
			"success",
			args{
				testutil.MustHexDecodeString("5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761"),
			},
			"0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotEncoded := EncodeZeroX(tt.args.in)
			if !assert.Equal(tt.wantEncoded, gotEncoded) {
				t.Errorf("EncodeZeroX() gotEncoded = %v, want %v", gotEncoded, tt.wantEncoded)
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
			testutil.MustHexDecodeString("5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761"),
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
