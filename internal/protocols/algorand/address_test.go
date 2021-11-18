// Copyright 2021 Finobo
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

package algorand

import (
	"testing"

	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/ed25519/ed25519test"
	"github.com/mailchain/mailchain/crypto/secp256k1/secp256k1test"
	"github.com/mailchain/mailchain/crypto/sr25519/sr25519test"
	"github.com/mailchain/mailchain/encoding"
	"github.com/stretchr/testify/assert"
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
			"alice",
			args{
				ed25519test.AlicePublicKey,
			},
			[]byte{0x72, 0x3c, 0xaa, 0x23, 0xa5, 0xb5, 0x11, 0xaf, 0x5a, 0xd7, 0xb7, 0xef, 0x60, 0x76, 0xe4, 0x14, 0xab, 0x7e, 0x75, 0xa9, 0xdc, 0x91, 0xe, 0xa6, 0xe, 0x41, 0x7a, 0x2b, 0x77, 0xa, 0x56, 0x71, 0x64, 0xea, 0x58, 0x50},
			false,
		},
		{
			"bob",
			args{
				ed25519test.BobPublicKey,
			},
			[]byte{0x2e, 0x32, 0x2f, 0x87, 0x40, 0xc6, 0x1, 0x72, 0x11, 0x1a, 0xc8, 0xea, 0xdc, 0xdd, 0xa2, 0x51, 0x2f, 0x90, 0xd0, 0x6d, 0xe, 0x50, 0x3e, 0xf1, 0x89, 0x97, 0x9a, 0x15, 0x9b, 0xec, 0xe1, 0xe8, 0x96, 0xc8, 0xe9, 0x2c},
			false,
		},
		{
			"err-secp256k1",
			args{
				secp256k1test.BobPublicKey,
			},
			nil,
			true,
		},
		{
			"err-sr25519",
			args{
				sr25519test.AlicePublicKey,
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

			if !assert.Equal(t, tt.want, got) {
				t.Errorf("Address() = %v, want %v", got, tt.want)
			}

			if err != nil {
				return
			}

			if !assert.Len(t, encoding.EncodeBase32(got), 58) {
				t.Error("Encoded address must have a length of 58")
			}
		})
	}
}
