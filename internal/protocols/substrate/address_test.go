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

package substrate

import (
	"testing"

	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/ed25519/ed25519test"
	"github.com/mailchain/mailchain/crypto/secp256k1/secp256k1test"
	"github.com/mailchain/mailchain/crypto/sr25519/sr25519test"
	"github.com/mailchain/mailchain/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func Test_prefixWithNetwork(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		network   string
		publicKey crypto.PublicKey
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			"edgeware-testnet",
			args{
				"edgeware-testnet",
				ed25519test.SofiaPublicKey,
			},
			[]byte{0x2a, 0x72, 0x3c, 0xaa, 0x23, 0xa5, 0xb5, 0x11, 0xaf, 0x5a, 0xd7, 0xb7, 0xef, 0x60, 0x76, 0xe4, 0x14, 0xab, 0x7e, 0x75, 0xa9, 0xdc, 0x91, 0xe, 0xa6, 0xe, 0x41, 0x7a, 0x2b, 0x77, 0xa, 0x56, 0x71},
			false,
		},
		{
			"invalid",
			args{
				"invalid",
				ed25519test.SofiaPublicKey,
			},
			nil,
			true,
		},
		{
			"polkadot-testnet",
			args{
				"polkadot-testnet",
				ed25519test.SofiaPublicKey,
			},
			[]byte{0x00, 0x72, 0x3c, 0xaa, 0x23, 0xa5, 0xb5, 0x11, 0xaf, 0x5a, 0xd7, 0xb7, 0xef, 0x60, 0x76, 0xe4, 0x14, 0xab, 0x7e, 0x75, 0xa9, 0xdc, 0x91, 0xe, 0xa6, 0xe, 0x41, 0x7a, 0x2b, 0x77, 0xa, 0x56, 0x71},
			false,
		},
		{
			"invalid",
			args{
				"invalid",
				ed25519test.SofiaPublicKey,
			},
			nil,
			true,
		},
		{
			"kusama-testnet",
			args{
				"kusama-testnet",
				ed25519test.SofiaPublicKey,
			},
			[]byte{0x02, 0x72, 0x3c, 0xaa, 0x23, 0xa5, 0xb5, 0x11, 0xaf, 0x5a, 0xd7, 0xb7, 0xef, 0x60, 0x76, 0xe4, 0x14, 0xab, 0x7e, 0x75, 0xa9, 0xdc, 0x91, 0xe, 0xa6, 0xe, 0x41, 0x7a, 0x2b, 0x77, 0xa, 0x56, 0x71},
			false,
		},
		{
			"invalid",
			args{
				"invalid",
				ed25519test.SofiaPublicKey,
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := prefixWithNetwork(tt.args.network, tt.args.publicKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("prefixWithNetwork() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(tt.want, got) {
				t.Errorf("prefixWithNetwork() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_addSS58Prefix(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		pubKey []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			"success",
			args{
				testutil.MustHexDecodeString("b14d"),
			},
			[]byte{0x53, 0x53, 0x35, 0x38, 0x50, 0x52, 0x45, 0xb1, 0x4d},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := addSS58Prefix(tt.args.pubKey); !assert.Equal(tt.want, got) {
				t.Errorf("addSS58Prefix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSS58AddressFormat(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		network   string
		publicKey crypto.PublicKey
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			"success-sofia-edgeware-testnet",
			args{
				"edgeware-testnet",
				ed25519test.SofiaPublicKey,
			},
			[]byte{0x2a, 0x72, 0x3c, 0xaa, 0x23, 0xa5, 0xb5, 0x11, 0xaf, 0x5a, 0xd7, 0xb7, 0xef, 0x60, 0x76, 0xe4, 0x14, 0xab, 0x7e, 0x75, 0xa9, 0xdc, 0x91, 0xe, 0xa6, 0xe, 0x41, 0x7a, 0x2b, 0x77, 0xa, 0x56, 0x71, 0x63, 0x83},
			false,
		},
		{
			"success-charlotte-edgeware-testnet",
			args{
				"edgeware-testnet",
				ed25519test.CharlottePublicKey,
			},
			[]byte{0x2a, 0x2e, 0x32, 0x2f, 0x87, 0x40, 0xc6, 0x1, 0x72, 0x11, 0x1a, 0xc8, 0xea, 0xdc, 0xdd, 0xa2, 0x51, 0x2f, 0x90, 0xd0, 0x6d, 0xe, 0x50, 0x3e, 0xf1, 0x89, 0x97, 0x9a, 0x15, 0x9b, 0xec, 0xe1, 0xe8, 0x6d, 0x48},
			false,
		},
		{
			"err-network",
			args{
				"invalid",
				ed25519test.SofiaPublicKey,
			},
			nil,
			true,
		},
		{
			"err-key-length",
			args{
				"edgeware-testnet",
				nil,
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SS58AddressFormat(tt.args.network, tt.args.publicKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("SS58AddressFormat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(tt.want, got) {
				t.Errorf("SS58AddressFormat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validPublicKeyType(t *testing.T) {
	type args struct {
		pubKey crypto.PublicKey
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"success-ed25519",
			args{
				ed25519test.SofiaPublicKey,
			},
			false,
		},
		{
			"success-sr25519",
			args{
				sr25519test.SofiaPublicKey,
			},
			false,
		},
		{
			"err-secp256k1",
			args{
				secp256k1test.SofiaPublicKey,
			},
			true,
		},
		{
			"err-nil",
			args{
				nil,
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validPublicKeyType(tt.args.pubKey); (err != nil) != tt.wantErr {
				t.Errorf("validPublicKeyType() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
