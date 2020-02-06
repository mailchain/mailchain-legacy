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

// Test PublicKeyEncoding
func TestPublicKeyEncoding(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want bool
	}{
		{
			"PublicKeyEncoding-empty-false",
			"",
			false,
		},
		{
			"PublicKeyEncoding-True",
			KindHex,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkEncod := PublicKeyEncoding()
			_, got := pkEncod[tt.in]

			if !assert.Equal(t, tt.want, got) {
				t.Errorf("PublicKeyEncoding() = %v, want %v", got, tt.want)
			}
		})
	}
}
