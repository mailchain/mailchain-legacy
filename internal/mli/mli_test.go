// Copyright 2022 Mailchain Ltd.
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

package mli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToAddress(t *testing.T) {
	tests := []struct {
		name string
		want map[uint64]string
	}{
		{
			"compatibility",
			map[uint64]string{
				0x1: "https://mcx.mx",
				0x2: "https://ipfs.io/ipfs",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToAddress(); !assert.Equal(t, tt.want, got) {
				t.Errorf("ToAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}
