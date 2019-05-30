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

package secp256k1

import (
	"crypto/ecdsa"
	"encoding/hex"
	"reflect"
	"testing"

	"github.com/mailchain/mailchain/crypto/secp256k1/secp256k1test"
	"github.com/stretchr/testify/assert"
)

func TestPublicKeyFromHex(t *testing.T) {
	type args struct {
		hex string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			"success-prefix",
			args{
				"0x" + hex.EncodeToString(charlottePublicKey().Bytes()),
			},
			charlottePublicKey().Bytes(),
			false,
		},
		{
			"success-no-prefix",
			args{
				hex.EncodeToString(charlottePublicKey().Bytes()),
			},
			charlottePublicKey().Bytes(),
			false,
		},
		{
			"success-no-fixed-prefix",
			args{
				"bdf6fb97c97c126b492186a4d5b28f34f0671a5aacc974da3bde0be93e45a1c50f89ceff72bd04ac9e25a04a1a6cb010aedaf65f91cec8ebe75901c49b63355d",
			},
			charlottePublicKey().Bytes(),
			false,
		},
		{
			"success-fixed-prefix",
			args{
				"0xbdf6fb97c97c126b492186a4d5b28f34f0671a5aacc974da3bde0be93e45a1c50f89ceff72bd04ac9e25a04a1a6cb010aedaf65f91cec8ebe75901c49b63355d",
			},
			charlottePublicKey().Bytes(),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PublicKeyFromHex(tt.args.hex)
			if (err != nil) != tt.wantErr {
				t.Errorf("PublicKeyFromHex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Bytes(), tt.want) {
				t.Errorf("PublicKeyFromHex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPublicKey_Address(t *testing.T) {
	assert := assert.New(t)
	type fields struct {
		ecdsa ecdsa.PublicKey
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{
			"success",
			fields{
				secp256k1test.ECDSAPublicKeyA(),
			},
			[]byte{0x8f, 0xd3, 0x79, 0x24, 0x68, 0x34, 0xea, 0xc7, 0x4b, 0x84, 0x19, 0xff, 0xda, 0x20, 0x2c, 0xf8, 0x5, 0x1f, 0x7a, 0x3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pk := PublicKey{
				ecdsa: tt.fields.ecdsa,
			}
			if got := pk.Address(); !assert.Equal(tt.want, got) {
				t.Errorf("PublicKey.Address() = %v, want %v", got, tt.want)
			}
		})
	}
}
