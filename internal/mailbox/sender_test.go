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

package mailbox

import (
	"testing"

	"github.com/mailchain/mailchain/crypto/cipher/aes256cbc"
	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func Test_encryptLocation(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		pk       crypto.PublicKey
		location string
	}
	type val struct {
		encryptedLen   int
		pk             crypto.PrivateKey
		wantDecryptErr bool
		location       string
	}
	tests := []struct {
		name    string
		args    args
		val     val
		wantErr bool
	}{
		{
			"testutil-charlotte",
			args{
				testutil.CharlottePublicKey,
				"http://test.com/location",
			},
			val{
				114,
				testutil.CharlottePrivateKey,
				false,
				"http://test.com/location",
			},
			false,
		},
		{
			"testutil-charlotte-incorrect-private-key",
			args{
				testutil.SofiaPublicKey,
				"http://test.com/location",
			},
			val{
				114,
				testutil.CharlottePrivateKey,
				true,
				"",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := encryptLocation(tt.args.pk, tt.args.location)
			if (err != nil) != tt.wantErr {
				t.Errorf("encryptLocation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(tt.val.encryptedLen, len(got)) {
				t.Errorf("len(encryptLocation()) = %v, want %v", got, tt.val.encryptedLen)
			}
			if err == nil {
				decrypter := aes256cbc.NewDecrypter(tt.val.pk)
				loc, err := decrypter.Decrypt(got)
				if (err != nil) != tt.val.wantDecryptErr {
					t.Errorf("decrypter.Decrypt() error = %v, wantDecryptErr %v", err, tt.val.wantDecryptErr)
					return
				}
				if !assert.Equal(tt.val.location, string(loc)) {
					t.Errorf("decryptedLocation = %v, want %v", string(loc), tt.val.location)
				}
			}
		})
	}
}
