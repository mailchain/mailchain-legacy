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
	"github.com/mailchain/mailchain/encoding"
	"github.com/mailchain/mailchain/encoding/encodingtest"
	"github.com/stretchr/testify/assert"
)

func Test_prefixWithNetwork(t *testing.T) {
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
			"edgeware-mainnet",
			args{
				EdgewareMainnet,
				ed25519test.SofiaPublicKey,
			},
			[]byte{0x7, 0x72, 0x3c, 0xaa, 0x23, 0xa5, 0xb5, 0x11, 0xaf, 0x5a, 0xd7, 0xb7, 0xef, 0x60, 0x76, 0xe4, 0x14, 0xab, 0x7e, 0x75, 0xa9, 0xdc, 0x91, 0xe, 0xa6, 0xe, 0x41, 0x7a, 0x2b, 0x77, 0xa, 0x56, 0x71},
			false,
		},
		{
			"edgeware-beresheet",
			args{
				EdgewareBeresheet,
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := prefixWithNetwork(tt.args.network, tt.args.publicKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("prefixWithNetwork() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("prefixWithNetwork() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_addSS58Prefix(t *testing.T) {
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
				encodingtest.MustDecodeHex("b14d"),
			},
			[]byte{0x53, 0x53, 0x35, 0x38, 0x50, 0x52, 0x45, 0xb1, 0x4d},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := addSS58Prefix(tt.args.pubKey); !assert.Equal(t, tt.want, got) {
				t.Errorf("addSS58Prefix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSS58AddressFormat(t *testing.T) {
	type args struct {
		network   string
		publicKey crypto.PublicKey
	}
	tests := []struct {
		name              string
		args              args
		want              []byte
		wantBase58Encoded string
		wantErr           bool
	}{
		// Mainnet
		{
			"ed25519-sofia-edgeware-mainnet",
			args{
				EdgewareMainnet,
				ed25519test.SofiaPublicKey,
			},
			[]byte{0x7, 0x72, 0x3c, 0xaa, 0x23, 0xa5, 0xb5, 0x11, 0xaf, 0x5a, 0xd7, 0xb7, 0xef, 0x60, 0x76, 0xe4, 0x14, 0xab, 0x7e, 0x75, 0xa9, 0xdc, 0x91, 0xe, 0xa6, 0xe, 0x41, 0x7a, 0x2b, 0x77, 0xa, 0x56, 0x71, 0xda, 0xb},
			"k6QHnYGfD7oioVwk2jRwR2H5uSW7fV2yiBU1BdBsJ7ebqHC",
			false,
		},
		{
			"ed25519-charlotte-edgeware-mainnet",
			args{
				EdgewareMainnet,
				ed25519test.CharlottePublicKey,
			},
			[]byte{0x7, 0x2e, 0x32, 0x2f, 0x87, 0x40, 0xc6, 0x1, 0x72, 0x11, 0x1a, 0xc8, 0xea, 0xdc, 0xdd, 0xa2, 0x51, 0x2f, 0x90, 0xd0, 0x6d, 0xe, 0x50, 0x3e, 0xf1, 0x89, 0x97, 0x9a, 0x15, 0x9b, 0xec, 0xe1, 0xe8, 0x9b, 0x76},
			"iZBvP2dtRccehTbT7syTGzGSiRQaAXhwCuk5K6GGJWGSBBK",
			false,
		},
		{
			"sr25519-sofia-edgeware-mainnet",
			args{
				EdgewareMainnet,
				sr25519test.SofiaPublicKey,
			},
			[]byte{0x7, 0x16, 0x9a, 0x11, 0x72, 0x18, 0x51, 0xf5, 0xdf, 0xf3, 0x54, 0x1d, 0xd5, 0xc4, 0xb0, 0xb4, 0x78, 0xac, 0x1c, 0xd0, 0x92, 0xc9, 0xd5, 0x97, 0x6e, 0x83, 0xda, 0xa0, 0xd0, 0x3f, 0x26, 0x62, 0xc, 0x35, 0x20},
			"i2FdaX7B2pXU3AApfkCSYPP8NDwLstfsU1wgb7uPq5zAtis",
			false,
		},
		{
			"sr25519-charlotte-edgeware-mainnet",
			args{
				EdgewareMainnet,
				sr25519test.CharlottePublicKey,
			},
			[]byte{0x7, 0x84, 0x62, 0x3e, 0x72, 0x52, 0xe4, 0x11, 0x38, 0xaf, 0x69, 0x4, 0xe1, 0xb0, 0x23, 0x4, 0xc9, 0x41, 0x62, 0x5f, 0x39, 0xe5, 0x76, 0x25, 0x89, 0x12, 0x5d, 0xc1, 0xa2, 0xf2, 0xcf, 0x2e, 0x30, 0xf0, 0xec},
			"kWCKFhkg5m5XmGz2pbcUYhJK7RAf6jLj1wLM5J5junonVTH",
			false,
		},
		// beresheet
		{
			"ed25519-sofia-edgeware-beresheet",
			args{
				EdgewareBeresheet,
				ed25519test.SofiaPublicKey,
			},
			[]byte{0x2a, 0x72, 0x3c, 0xaa, 0x23, 0xa5, 0xb5, 0x11, 0xaf, 0x5a, 0xd7, 0xb7, 0xef, 0x60, 0x76, 0xe4, 0x14, 0xab, 0x7e, 0x75, 0xa9, 0xdc, 0x91, 0xe, 0xa6, 0xe, 0x41, 0x7a, 0x2b, 0x77, 0xa, 0x56, 0x71, 0x63, 0x83},
			"5EeVLWVP7PjFBrfBSfXG9xL9htSo35XFada4Y78j5DviSz9Q",
			false,
		},
		{
			"ed25519-charlotte-edgeware-beresheet",
			args{
				EdgewareBeresheet,
				ed25519test.CharlottePublicKey,
			},
			[]byte{0x2a, 0x2e, 0x32, 0x2f, 0x87, 0x40, 0xc6, 0x1, 0x72, 0x11, 0x1a, 0xc8, 0xea, 0xdc, 0xdd, 0xa2, 0x51, 0x2f, 0x90, 0xd0, 0x6d, 0xe, 0x50, 0x3e, 0xf1, 0x89, 0x97, 0x9a, 0x15, 0x9b, 0xec, 0xe1, 0xe8, 0x6d, 0x48},
			"5D7Gy6ykLcE47kcq9kfofpJ94hRhVaZvY8JLcEboUEKLHRYs",
			false,
		},
		{
			"sr25519-sofia-edgeware-beresheet",
			args{
				EdgewareBeresheet,
				sr25519test.SofiaPublicKey,
			},
			[]byte{0x2a, 0x16, 0x9a, 0x11, 0x72, 0x18, 0x51, 0xf5, 0xdf, 0xf3, 0x54, 0x1d, 0xd5, 0xc4, 0xb0, 0xb4, 0x78, 0xac, 0x1c, 0xd0, 0x92, 0xc9, 0xd5, 0x97, 0x6e, 0x83, 0xda, 0xa0, 0xd0, 0x3f, 0x26, 0x62, 0xc, 0x46, 0x4b},
			"5CaLgJUDdDRxw6KQXJY2f5hFkMEEGHvtUPQYDWdSbku42Dv2",
			false,
		},
		{
			"sr25519-charlotte-edgeware-beresheet",
			args{
				EdgewareBeresheet,
				sr25519test.CharlottePublicKey,
			},
			[]byte{0x2a, 0x84, 0x62, 0x3e, 0x72, 0x52, 0xe4, 0x11, 0x38, 0xaf, 0x69, 0x4, 0xe1, 0xb0, 0x23, 0x4, 0xc9, 0x41, 0x62, 0x5f, 0x39, 0xe5, 0x76, 0x25, 0x89, 0x12, 0x5d, 0xc1, 0xa2, 0xf2, 0xcf, 0x2e, 0x30, 0xe0, 0x2a},
			"5F4HMyes8GNWzpSDjTPSh61Aw6RTaWmZKwKvszocwqbsdn4h",
			false,
		},
		{
			"err-network",
			args{
				"invalid",
				sr25519test.SofiaPublicKey,
			},
			nil,
			"",
			true,
		},
		{
			"err-key-length",
			args{
				"edgeware-beresheet",
				nil,
			},
			nil,
			"",
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
			assert.Equal(t, tt.wantBase58Encoded, encoding.EncodeBase58(got))
			if !assert.Equal(t, tt.want, got) {
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
			"success-schnorrkel",
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
