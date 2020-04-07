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

	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/cipher"
	"github.com/mailchain/mailchain/crypto/cipher/aes256cbc"
	"github.com/mailchain/mailchain/crypto/cipher/nacl"
	"github.com/mailchain/mailchain/crypto/cipher/noop"
	"github.com/mailchain/mailchain/crypto/secp256k1/secp256k1test"
	"github.com/stretchr/testify/assert"
)

func TestGetEncrypter(t *testing.T) {
	type args struct {
		encryption string
		pubKey     crypto.PublicKey
	}
	tests := []struct {
		name    string
		args    args
		want    cipher.Encrypter
		wantErr bool
	}{
		{
			"aes256cbc",
			args{
				"aes256cbc",
				secp256k1test.SofiaPublicKey,
			},
			func() cipher.Encrypter {
				encrypter, _ := aes256cbc.NewEncrypter(secp256k1test.SofiaPublicKey)
				return encrypter
			}(),
			false,
		},
		{
			"nacl-ecdh",
			args{
				"nacl-ecdh",
				secp256k1test.SofiaPublicKey,
			},
			func() cipher.Encrypter {
				encrypter, _ := nacl.NewEncrypter(secp256k1test.SofiaPublicKey)
				return encrypter
			}(),
			false,
		},
		{
			"noop",
			args{
				"noop",
				secp256k1test.SofiaPublicKey,
			},
			func() cipher.Encrypter {
				encrypter, _ := noop.NewEncrypter(secp256k1test.SofiaPublicKey)
				return encrypter
			}(),
			false,
		},
		{
			"err-empty",
			args{
				"",
				secp256k1test.SofiaPublicKey,
			},
			nil,
			true,
		},
		{
			"err-invalid",
			args{
				"invalid",
				secp256k1test.SofiaPublicKey,
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetEncrypter(tt.args.encryption, tt.args.pubKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetEncrypter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("GetEncrypter() = %v, want %v", got, tt.want)
			}
		})
	}
}
