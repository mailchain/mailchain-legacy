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

package crypto

import (
	"encoding/hex"
	"testing"

	"github.com/mailchain/mailchain/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestCreateLocationHash(t *testing.T) {
	assert := assert.New(t)
	cases := []struct {
		name     string
		original []byte
		expected []byte
		err      error
	}{
		{"2c8432ca28ce929b86a47f2d40413d161f591f8985229060491573d83f82f292f4dc68f918446332837aa57cd5145235cc40702d962cbb53ac27fb2246fb6cba",
			testutil.MustHexDecodeString("2c8432ca28ce929b86a47f2d40413d161f591f8985229060491573d83f82f292f4dc68f918446332837aa57cd5145235cc40702d962cbb53ac27fb2246fb6cba"),
			testutil.MustHexDecodeString("2204abd5fcd4"),
			nil,
		},
		{"022c8432ca28ce929b86a47f2d40413d161f591f8985229060491573d83f82f292",
			testutil.MustHexDecodeString("022c8432ca28ce929b86a47f2d40413d161f591f8985229060491573d83f82f292"),
			testutil.MustHexDecodeString("2204be6f4863"),
			nil,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := createLocationHash(tc.original)
			assert.EqualValues(hex.EncodeToString(tc.expected), hex.EncodeToString(actual))
			assert.Equal(tc.err, err)
		})
	}
}
