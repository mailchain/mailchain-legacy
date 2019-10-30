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

package ethereum

import (
	"reflect"
	"testing"

	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/ed25519/ed25519test"
	"github.com/mailchain/mailchain/crypto/secp256k1/secp256k1test"
)

func TestAddress(t *testing.T) {
	type args struct {
		pubKey crypto.PublicKey
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			"sofia",
			args{
				secp256k1test.SofiaPublicKey,
			},
			[]byte{0xd5, 0xab, 0x4c, 0xe3, 0x60, 0x5c, 0xd5, 0x90, 0xdb, 0x60, 0x9b, 0x6b, 0x5c, 0x89, 0x1, 0xfd, 0xb2, 0xef, 0x7f, 0xe6},
			false,
		},
		{
			"charlotte",
			args{
				secp256k1test.CharlottePublicKey,
			},
			[]byte{0x92, 0xd8, 0xf1, 0x2, 0x48, 0xc6, 0xa3, 0x95, 0x3c, 0xc3, 0x69, 0x2a, 0x89, 0x46, 0x55, 0xad, 0x5, 0xd6, 0x1e, 0xfb},
			false,
		},
		{
			"err-ed25519",
			args{
				ed25519test.CharlottePublicKey,
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Address(tt.args.pubKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("Address() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Address() = %v, want %v", got, tt.want)
			}
		})
	}
}
