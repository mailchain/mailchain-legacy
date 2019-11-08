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
	"log"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
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
		{
			"err-could-not-decode",
			args{
				"0xbdf6fb97c97c126b492",
			},
			nil,
			true,
		},
		{
			"err-invalid-length",
			args{
				"0xbdf6fb97c97c126b492186a4d5b28f34f0671a5aacc974da3bde0be93e45a1c50f89ceff72bd04ac9e25a04a1a6cb010aedaf65f91cec8ebe75901c4",
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PublicKeyFromHex(tt.args.hex)
			if (err != nil) != tt.wantErr {
				t.Errorf("PublicKeyFromHex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			var gotBytes []byte
			if got != nil {
				gotBytes = got.Bytes()
			}
			if !reflect.DeepEqual(gotBytes, tt.want) {
				t.Errorf("PublicKeyFromHex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPublicKeyFromBytes(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		keyBytes []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			"success-65",
			args{
				func() []byte {
					pub := make([]byte, 65)
					pub[0] = byte(4)
					copy(pub[1:], ecdsaPrivateKeySofia().X.Bytes())
					copy(pub[33:], ecdsaPrivateKeySofia().Y.Bytes())
					return pub
				}(),
			},
			[]byte{0x2, 0x69, 0xd9, 0x8, 0x51, 0xe, 0x35, 0x5b, 0xeb, 0x1d, 0x5b, 0xf2, 0xdf, 0x81, 0x29, 0xe5, 0xb6, 0x40, 0x1e, 0x19, 0x69, 0x89, 0x1e, 0x80, 0x16, 0xa0, 0xb2, 0x30, 0x7, 0x39, 0xbb, 0xb0, 0x6},
			false,
		},
		{
			"err-65",
			args{
				[]byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
			},
			nil,
			true,
		},
		{
			"success-64",
			args{
				func() []byte {
					pub := make([]byte, 64)
					copy(pub, ecdsaPrivateKeySofia().X.Bytes())
					copy(pub[32:], ecdsaPrivateKeySofia().Y.Bytes())
					return pub
				}(),
			},
			[]byte{0x2, 0x69, 0xd9, 0x8, 0x51, 0xe, 0x35, 0x5b, 0xeb, 0x1d, 0x5b, 0xf2, 0xdf, 0x81, 0x29, 0xe5, 0xb6, 0x40, 0x1e, 0x19, 0x69, 0x89, 0x1e, 0x80, 0x16, 0xa0, 0xb2, 0x30, 0x7, 0x39, 0xbb, 0xb0, 0x6},
			false,
		},
		{
			"err-64",
			args{
				[]byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
			},
			nil,
			true,
		},
		{
			"success-33",
			args{
				charlottePublicKey().Bytes(),
			},
			[]byte{0x3, 0xbd, 0xf6, 0xfb, 0x97, 0xc9, 0x7c, 0x12, 0x6b, 0x49, 0x21, 0x86, 0xa4, 0xd5, 0xb2, 0x8f, 0x34, 0xf0, 0x67, 0x1a, 0x5a, 0xac, 0xc9, 0x74, 0xda, 0x3b, 0xde, 0xb, 0xe9, 0x3e, 0x45, 0xa1, 0xc5},
			false,
		},
		{
			"err-63",
			args{
				[]byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
			},
			nil,
			true,
		},
		{
			"err",
			args{
				[]byte{0x3, 0xbd, 0xf6, 0xfb, 0x97},
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PublicKeyFromBytes(tt.args.keyBytes)
			if (err != nil) != tt.wantErr {
				t.Errorf("PublicKeyFromBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			var gotBytes []byte
			if got != nil {
				gotBytes = got.Bytes()
			}
			if !assert.Equal(tt.want, gotBytes) {
				t.Errorf("PublicKeyFromBytes() = %v, want %v", gotBytes, tt.want)
			}
		})
	}
}

func TestPublicKey_Bytes(t *testing.T) {
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
			"charlotte",
			fields{
				func() ecdsa.PublicKey {
					b, _ := hex.DecodeString("01901E63389EF02EAA7C5782E08B40D98FAEF835F28BD144EECF5614A415943F")
					key, err := crypto.ToECDSA(b)
					if err != nil {
						log.Fatal(err)
					}
					return key.PublicKey
				}(),
			},
			[]byte{0x2, 0x69, 0xd9, 0x8, 0x51, 0xe, 0x35, 0x5b, 0xeb, 0x1d, 0x5b, 0xf2, 0xdf, 0x81, 0x29, 0xe5, 0xb6, 0x40, 0x1e, 0x19, 0x69, 0x89, 0x1e, 0x80, 0x16, 0xa0, 0xb2, 0x30, 0x7, 0x39, 0xbb, 0xb0, 0x6},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pk := PublicKey{
				ecdsa: tt.fields.ecdsa,
			}
			if got := pk.Bytes(); !assert.Equal(tt.want, got) {
				t.Errorf("PublicKey.Bytes() = %v, want %v", got, tt.want)
			}
		})
	}
}
