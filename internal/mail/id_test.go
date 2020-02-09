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

package mail

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewID(t *testing.T) {
	id, err := NewID()
	assert.NoError(t, err)
	assert.Len(t, id, 44)
}

func TestFromHexString(t *testing.T) {
	type args struct {
		hex string
	}
	tests := []struct {
		name    string
		args    args
		want    ID
		wantErr bool
	}{
		{
			"success",
			args{
				"47eca011e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471",
			},
			ID{0x47, 0xec, 0xa0, 0x11, 0xe3, 0x2b, 0x52, 0xc7, 0x10, 0x5, 0xad, 0x8a, 0x8f, 0x75, 0xe1, 0xb4, 0x4c, 0x92, 0xc9, 0x9f, 0xd1, 0x2e, 0x43, 0xbc, 0xcf, 0xe5, 0x71, 0xe3, 0xc2, 0xd1, 0x3d, 0x2e, 0x9a, 0x82, 0x6a, 0x55, 0xf, 0x5f, 0xf6, 0x3b, 0x24, 0x7a, 0xf4, 0x71},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FromHexString(tt.args.hex)
			if (err != nil) != tt.wantErr {
				t.Errorf("FromHexString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("FromHexString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestID_HexString(t *testing.T) {
	tests := []struct {
		name string
		id   ID
		want string
	}{
		{
			"success",
			ID{0x47, 0xec, 0xa0, 0x11, 0xe3, 0x2b, 0x52, 0xc7, 0x10, 0x5, 0xad, 0x8a, 0x8f, 0x75, 0xe1, 0xb4, 0x4c, 0x92, 0xc9, 0x9f, 0xd1, 0x2e, 0x43, 0xbc, 0xcf, 0xe5, 0x71, 0xe3, 0xc2, 0xd1, 0x3d, 0x2e, 0x9a, 0x82, 0x6a, 0x55, 0xf, 0x5f, 0xf6, 0x3b, 0x24, 0x7a, 0xf4, 0x71},
			"47eca011e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.id.HexString(); got != tt.want {
				t.Errorf("ID.HexString() = %v, want %v", got, tt.want)
			}
		})
	}
}
