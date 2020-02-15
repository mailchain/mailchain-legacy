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

package address

import (
	"testing"

	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/ed25519/ed25519test"
	"github.com/mailchain/mailchain/crypto/secp256k1/secp256k1test"
	"github.com/mailchain/mailchain/internal/protocols/substrate"
	"github.com/stretchr/testify/assert"
)

func TestFromPublicKey(t *testing.T) {
	type args struct {
		pubKey   crypto.PublicKey
		protocol string
		network  string
	}
	tests := []struct {
		name        string
		args        args
		wantAddress []byte
		wantErr     bool
	}{
		{
			"ethereum",
			args{
				secp256k1test.SofiaPublicKey,
				"ethereum",
				"mainnet",
			},
			[]byte{0xd5, 0xab, 0x4c, 0xe3, 0x60, 0x5c, 0xd5, 0x90, 0xdb, 0x60, 0x9b, 0x6b, 0x5c, 0x89, 0x1, 0xfd, 0xb2, 0xef, 0x7f, 0xe6},
			false,
		},
		{
			"substrate",
			args{
				ed25519test.SofiaPublicKey,
				"substrate",
				substrate.EdgewareMainnet,
			},
			[]byte{0x7, 0x72, 0x3c, 0xaa, 0x23, 0xa5, 0xb5, 0x11, 0xaf, 0x5a, 0xd7, 0xb7, 0xef, 0x60, 0x76, 0xe4, 0x14, 0xab, 0x7e, 0x75, 0xa9, 0xdc, 0x91, 0xe, 0xa6, 0xe, 0x41, 0x7a, 0x2b, 0x77, 0xa, 0x56, 0x71, 0xda, 0xb},
			false,
		},
		{
			"err-invalid-protocol",
			args{
				secp256k1test.SofiaPublicKey,
				"invalid",
				"mainnet",
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAddress, err := FromPublicKey(tt.args.pubKey, tt.args.protocol, tt.args.network)
			if (err != nil) != tt.wantErr {
				t.Errorf("FromPublicKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.wantAddress, gotAddress) {
				t.Errorf("FromPublicKey() = %v, want %v", gotAddress, tt.wantAddress)
			}
		})
	}
}
