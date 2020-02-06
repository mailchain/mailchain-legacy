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

package crypto_test

import (
	"testing"

	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/encoding"
	"github.com/mailchain/mailchain/encoding/encodingtest"
	"github.com/stretchr/testify/assert"
)

func TestCreateIntegrityHash(t *testing.T) {
	cases := []struct {
		name     string
		original []byte
		expected []byte
	}{
		{"2c8432ca28ce929b86a47f2d40413d161f591f8985229060491573d83f82f292f4dc68f918446332837aa57cd5145235cc40702d962cbb53ac27fb2246fb6cba",
			encodingtest.MustDecodeHex("2c8432ca28ce929b86a47f2d40413d161f591f8985229060491573d83f82f292f4dc68f918446332837aa57cd5145235cc40702d962cbb53ac27fb2246fb6cba"),
			encodingtest.MustDecodeHex("2204abd5fcd4"),
		},
		{"022c8432ca28ce929b86a47f2d40413d161f591f8985229060491573d83f82f292",
			encodingtest.MustDecodeHex("022c8432ca28ce929b86a47f2d40413d161f591f8985229060491573d83f82f292"),
			encodingtest.MustDecodeHex("2204be6f4863"),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual := crypto.CreateIntegrityHash(tc.original)
			assert.EqualValues(t, encoding.EncodeHex(tc.expected), encoding.EncodeHex(actual))
		})
	}
}

func TestCreateMessageHash(t *testing.T) {
	cases := []struct {
		name     string
		original []byte
		expected []byte
	}{
		{"2c8432ca28ce929b86a47f2d40413d161f591f8985229060491573d83f82f292f4dc68f918446332837aa57cd5145235cc40702d962cbb53ac27fb2246fb6cba",
			encodingtest.MustDecodeHex("2c8432ca28ce929b86a47f2d40413d161f591f8985229060491573d83f82f292f4dc68f918446332837aa57cd5145235cc40702d962cbb53ac27fb2246fb6cba"),
			encodingtest.MustDecodeHex("16202b3cde1b72727d0b38daa592efae7117b86e7c2f5646543e2ae0f86f64b2922a"),
		},
		{"022c8432ca28ce929b86a47f2d40413d161f591f8985229060491573d83f82f292",
			encodingtest.MustDecodeHex("022c8432ca28ce929b86a47f2d40413d161f591f8985229060491573d83f82f292"),
			encodingtest.MustDecodeHex("1620671f6f840e08b9c6b3e2125e0381dd5da5578a698eb97a357f1015552263aec6"),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual := crypto.CreateMessageHash(tc.original)
			assert.EqualValues(t, encoding.EncodeHex(tc.expected), encoding.EncodeHex(actual))
		})
	}
}
