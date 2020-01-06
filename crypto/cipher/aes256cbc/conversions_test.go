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

package aes256cbc

import (
	"testing"

	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/ed25519/ed25519test"
	"github.com/mailchain/mailchain/crypto/secp256k1"
	"github.com/mailchain/mailchain/crypto/sr25519/sr25519test"
	"github.com/mailchain/mailchain/encoding"
	"github.com/mailchain/mailchain/crypto/secp256k1/secp256k1test"
)

func Test_asPrivateECIES(t *testing.T) {
	type args struct {
		pk crypto.PrivateKey
	}
	tests := []struct {
		name    string
		args    args
		wantNil bool
		wantErr bool
	}{
		{
			"success-secp256k1-sofia-val",
			args{
				func() secp256k1.PrivateKey {
					t := secp256k1test.SofiaPrivateKey.(*secp256k1.PrivateKey)
					return *t
				}(),
			},
			false,
			false,
		},
		{
			"success-secp256k1-sofia-pointer",
			args{
				func() *secp256k1.PrivateKey {
					return secp256k1test.SofiaPrivateKey.(*secp256k1.PrivateKey)
				}(),
			},
			false,
			false,
		},
		{
			"err-unsupported",
			args{
				ed25519test.SofiaPrivateKey,
			},
			true,
			true,
		},
		{
			"err-unsupported-sr25519",
			args{
				sr25519test.SofiaPrivateKey,
			},
			true,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := asPrivateECIES(tt.args.pk)
			if (err != nil) != tt.wantErr {
				t.Errorf("asPrivateECIES() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if (got == nil) != tt.wantNil {
				t.Errorf("asPrivateECIES() = %v, wantNil %v", got, tt.wantNil)
			}
		})
	}
}

func Test_asPublicECIES(t *testing.T) {
	type args struct {
		pk crypto.PublicKey
	}
	tests := []struct {
		name    string
		args    args
		wantNil bool
		wantErr bool
	}{
		{
			"success-secp256k1-sofia-pointer",
			args{
				func() crypto.PublicKey {
					return secp256k1test.SofiaPublicKey.(*secp256k1.PublicKey)
				}(),
			},
			false,
			false,
		},
		{
			"success-secp256k1-sofia-val",
			args{
				func() crypto.PublicKey {
					pk := secp256k1test.SofiaPublicKey.(*secp256k1.PublicKey)
					return *pk
				}(),
			},
			false,
			false,
		},
		{
			"err-invalid",
			args{
				ed25519test.SofiaPublicKey,
			},
			true,
			true,
		},
		{
			"err-invalid-sr25519",
			args{
				sr25519test.CharlottePublicKey,
			},
			true,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := asPublicECIES(tt.args.pk)
			if (err != nil) != tt.wantErr {
				t.Errorf("asPublicECIES() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got == nil) != tt.wantNil {
				t.Errorf("asPublicECIES() = %v, wantNil %v", got, tt.wantNil)
			}
		})
	}
}
