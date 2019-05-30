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

package keystore

import (
	"testing"

	"github.com/mailchain/mailchain/crypto/cipher"
	"github.com/mailchain/mailchain/crypto/cipher/aes256cbc"
	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/internal/encoding"
	"github.com/mailchain/mailchain/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestDecrypter(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		cipherType byte
		pk         crypto.PrivateKey
	}
	tests := []struct {
		name    string
		args    args
		want    cipher.Decrypter
		wantErr bool
	}{
		{
			"ethereum",
			args{
				encoding.AES256CBC,
				testutil.CharlottePrivateKey,
			},
			aes256cbc.NewDecrypter(testutil.CharlottePrivateKey),
			false,
		},
		{
			"err",
			args{
				0xFF,
				testutil.CharlottePrivateKey,
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Decrypter(tt.args.cipherType, tt.args.pk)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decrypter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(tt.want, got) {
				t.Errorf("Decrypter() = %v, want %v", got, tt.want)
			}
		})
	}
}
