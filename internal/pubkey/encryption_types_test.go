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
)

func TestEncryptionMethods(t *testing.T) {
	type args struct {
		kind string
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
				crypto.ED25519,
			},
			[]string{encrypter.NACL, encrypter.NoOperation},
			false,
		},
		{
			"secp256k1",
			args{
				crypto.SECP256K1,
			},
			[]string{encrypter.AES256CBC, encrypter.NoOperation},
			false,
		},
		{
			"unknown",
			args{
				"unknown",
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EncryptionMethods(tt.args.kind)

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
