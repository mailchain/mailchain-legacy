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

package multikey

import (
	"testing"

	"github.com/mailchain/mailchain/crypto/ed25519/ed25519test"
	"github.com/mailchain/mailchain/crypto/secp256k1/secp256k1test"
	"github.com/mailchain/mailchain/crypto/sr25519/sr25519test"
	"github.com/stretchr/testify/assert"
)

func TestPublicKeyFromBytes(t *testing.T) {
	type args struct {
		hex     string
		keyType []byte
	}
	tests := []struct {
		name      string
		args      args
		wantBytes []byte
		wantErr   bool
	}{
		{
			"secp256k1",
			args{
				"secp256k1",
				secp256k1test.AlicePublicKey.Bytes(),
			},
			secp256k1test.AlicePublicKey.Bytes(),
			false,
		},
		{
			"ed25519",
			args{
				"ed25519",
				ed25519test.AlicePublicKey.Bytes(),
			},
			ed25519test.AlicePublicKey.Bytes(),
			false,
		},
		{
			"sr25519",
			args{
				"sr25519",
				sr25519test.AlicePublicKey.Bytes(),
			},
			sr25519test.AlicePublicKey.Bytes(),
			false,
		},
		{
			"err",
			args{
				"unknown",
				secp256k1test.AlicePublicKey.Bytes(),
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PublicKeyFromBytes(tt.args.hex, tt.args.keyType)
			if (err != nil) != tt.wantErr {
				t.Errorf("PublicKeyFromBytes() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != nil {
				if !assert.EqualValues(t, tt.wantBytes, got.Bytes()) {
					t.Errorf("PublicKeyFromBytes() = %v, want %v", got, tt.wantBytes)
				}
			}
			if got == nil {
				if !assert.Nil(t, tt.wantBytes) {
					t.Errorf("PublicKeyFromBytes() = %v, want %v", got, tt.wantBytes)
				}
			}

		})
	}
}
