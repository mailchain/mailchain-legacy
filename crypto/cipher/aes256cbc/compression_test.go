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

package aes256cbc

import (
	"encoding/hex"
	"testing"

	"github.com/mailchain/mailchain/internal/encoding/encodingtest"
	"github.com/stretchr/testify/assert"
)

func TestCompress(t *testing.T) {
	assert := assert.New(t)
	cases := []struct {
		name     string
		original []byte
		want     []byte
		wantErr  bool
	}{
		{"no prefix:022c8432ca28ce929b86a47f2d40413d161f591f8985229060491573d83f82f292",
			encodingtest.MustDecodeHex("2c8432ca28ce929b86a47f2d40413d161f591f8985229060491573d83f82f292f4dc68f918446332837aa57cd5145235cc40702d962cbb53ac27fb2246fb6cba"),
			encodingtest.MustDecodeHex("022c8432ca28ce929b86a47f2d40413d161f591f8985229060491573d83f82f292"),
			false,
		},
		{"with prefix:022c8432ca28ce929b86a47f2d40413d161f591f8985229060491573d83f82f292",
			encodingtest.MustDecodeHex("042c8432ca28ce929b86a47f2d40413d161f591f8985229060491573d83f82f292f4dc68f918446332837aa57cd5145235cc40702d962cbb53ac27fb2246fb6cba"),
			encodingtest.MustDecodeHex("022c8432ca28ce929b86a47f2d40413d161f591f8985229060491573d83f82f292"),
			false,
		},
		{"err-invalid-key-length",
			encodingtest.MustDecodeHex("042c8432ca28ce929b86a47f2d40413d161f1f8985229060491573d83f82f292f4dc68f918446332837aa57cd5145235cc40702d962cbb53ac27fb2246fb6cba"),
			nil,
			true,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := compress(tt.original)
			if (err != nil) != tt.wantErr {
				t.Errorf("compress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(tt.want, got) {
				t.Errorf("compress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecompress(t *testing.T) {
	assert := assert.New(t)
	cases := []struct {
		name     string
		original []byte
		expected []byte
	}{
		{"022c8432ca28ce929b86a47f2d40413d161f591f8985229060491573d83f82f292",
			encodingtest.MustDecodeHex("022c8432ca28ce929b86a47f2d40413d161f591f8985229060491573d83f82f292"),
			encodingtest.MustDecodeHex("042c8432ca28ce929b86a47f2d40413d161f591f8985229060491573d83f82f292f4dc68f918446332837aa57cd5145235cc40702d962cbb53ac27fb2246fb6cba"),
		},
		{"03a34d6aef3eb42335fb3cacb59478c0b44c0bbeb8bb4ca427dbc7044157a5d24b",
			encodingtest.MustDecodeHex("03a34d6aef3eb42335fb3cacb59478c0b44c0bbeb8bb4ca427dbc7044157a5d24b"),
			encodingtest.MustDecodeHex("04a34d6aef3eb42335fb3cacb59478c0b44c0bbeb8bb4ca427dbc7044157a5d24b4adf14868d8449c9b3e50d3d6338f3e5a2d3445abe679cddbe75cb893475806f"),
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			actual := decompress(tt.original)
			assert.EqualValues(hex.EncodeToString(tt.expected), hex.EncodeToString(actual))
		})
	}
}
