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

package pubkey

import (
	"reflect"
	"testing"

	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/cipher/encrypter"
	"github.com/mailchain/mailchain/crypto/ed25519/ed25519test"
	"github.com/mailchain/mailchain/crypto/secp256k1/secp256k1test"
	"github.com/mailchain/mailchain/crypto/sr25519/sr25519test"
)

func TestEncryptionMethods(t *testing.T) {
	type args struct {
		key crypto.PublicKey
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			"ed25519",
			args{
				ed25519test.AlicePublicKey,
			},
			[]string{encrypter.NACLECDH, encrypter.NoOperation},
			false,
		},
		{
			"sr25519",
			args{
				sr25519test.AlicePublicKey,
			},
			[]string{encrypter.NACLECDH, encrypter.NoOperation},
			false,
		},
		{
			"secp256k1",
			args{
				secp256k1test.AlicePublicKey,
			},
			[]string{encrypter.AES256CBC, encrypter.NoOperation},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EncryptionMethods(tt.args.key)

			if (err != nil) != tt.wantErr {
				t.Errorf("EncryptionMethods() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EncryptionMethods() gotTypes = %v, want %v", got, tt.want)
			}
		})
	}
}
