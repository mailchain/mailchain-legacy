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

	keys "github.com/mailchain/mailchain/crypto"
	crypto "github.com/mailchain/mailchain/crypto/cipher"
	"github.com/mailchain/mailchain/crypto/cipher/aes256cbc"
	"github.com/mailchain/mailchain/crypto/ed25519/ed25519test"
	"github.com/mailchain/mailchain/crypto/secp256k1/secp256k1test"
	"github.com/stretchr/testify/assert"
)

func TestGetEnrypter(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		encrypter string
		pubKey    keys.PublicKey
	}

	testKeys := struct {
		secp256k1Key keys.PublicKey
		ed25519Key   keys.PublicKey
	}{
		secp256k1Key: secp256k1test.CharlottePublicKey,
		ed25519Key:   ed25519test.CharlottePublicKey,
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
				testKeys.ed25519Key,
			},
			nil,
			true,
		},
		{
			"empty",
			args{
				"",
				testKeys.ed25519Key,
			},
			nil,
			true,
		},
		{
			"aes",
			args{
				"aes256cbc",
				testKeys.secp256k1Key,
			},
			func() crypto.Encrypter {
				encrypter, _ := aes256cbc.NewEncrypter(testKeys.secp256k1Key)
				return encrypter
			}(),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetEncrypter(tt.args.encrypter, tt.args.pubKey)
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
