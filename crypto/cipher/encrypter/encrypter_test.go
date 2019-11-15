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

package encrypter

import (
	"testing"

	crypto "github.com/mailchain/mailchain/crypto/cipher"
	"github.com/mailchain/mailchain/crypto/cipher/aes256cbc"
	"github.com/mailchain/mailchain/crypto/cipher/nacl"
	"github.com/stretchr/testify/assert"
)

func TestGetEnrypter(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		encrypter string
	}
	tests := []struct {
		name    string
		args    args
		want    crypto.Encrypter
		wantErr bool
	}{
		{
			"invalid",
			args{
				"test-invalid",
			},
			nil,
			true,
		},
		{
			"empty",
			args{
				"",
			},
			nil,
			true,
		},
		{
			"aes",
			args{
				"aes256cbc",
			},
			aes256cbc.NewEncrypter(),
			false,
		},
		{
			"nacl",
			args{
				"nacl",
			},
			nacl.NewEncrypter(),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetEncrypter(tt.args.encrypter)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetEncrypter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.IsType(tt.want, got) {
				t.Errorf("GetEncrypter() = %v, want %v", got, tt.want)
			}
		})
	}
}
