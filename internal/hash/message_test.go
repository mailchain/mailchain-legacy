// Copyright 2022 Mailchain Ltd
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

package hash

import (
	"reflect"
	"testing"

	"github.com/mailchain/mailchain/encoding/encodingtest"
	"github.com/multiformats/go-multihash"
)

func TestCreateMessageHash(t *testing.T) {
	type args struct {
		encodedData []byte
	}
	tests := []struct {
		name string
		args args
		want multihash.Multihash
	}{
		{
			"2c8432ca28ce929b86a47f2d40413d161f591f8985229060491573d83f82f292f4dc68f918446332837aa57cd5145235cc40702d962cbb53ac27fb2246fb6cba",
			args{
				encodingtest.MustDecodeHex("2c8432ca28ce929b86a47f2d40413d161f591f8985229060491573d83f82f292f4dc68f918446332837aa57cd5145235cc40702d962cbb53ac27fb2246fb6cba"),
			},
			encodingtest.MustDecodeHex("16202b3cde1b72727d0b38daa592efae7117b86e7c2f5646543e2ae0f86f64b2922a"),
		},
		{
			"022c8432ca28ce929b86a47f2d40413d161f591f8985229060491573d83f82f292",
			args{
				encodingtest.MustDecodeHex("022c8432ca28ce929b86a47f2d40413d161f591f8985229060491573d83f82f292"),
			},
			encodingtest.MustDecodeHex("1620671f6f840e08b9c6b3e2125e0381dd5da5578a698eb97a357f1015552263aec6"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateMessageHash(tt.args.encodedData); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateMessageHash() = %v, want %v", got, tt.want)
			}
		})
	}
}
