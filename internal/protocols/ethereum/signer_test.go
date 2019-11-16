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
	"crypto/ecdsa"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/ed25519/ed25519test"
	"github.com/mailchain/mailchain/crypto/secp256k1"
	"github.com/mailchain/mailchain/crypto/secp256k1/secp256k1test"
	"github.com/mailchain/mailchain/internal/mailbox/signer"
	"github.com/stretchr/testify/assert"
)

func TestNewSigner(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		privateKey crypto.PrivateKey
	}
	tests := []struct {
		name    string
		args    args
		want    *Signer
		wantErr bool
	}{
		{
			"success",
			args{
				secp256k1test.CharlottePrivateKey,
			},
			&Signer{
				secp256k1test.CharlottePrivateKey,
			},
			false,
		},
		{
			"invalid-key",
			args{
				ed25519test.CharlottePrivateKey,
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewSigner(tt.args.privateKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSigner() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !assert.Equal(tt.want, got) {
				t.Errorf("NewSigner() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSigner_Sign(t *testing.T) {
	type fields struct {
		privateKey crypto.PrivateKey
	}
	type args struct {
		opts signer.Options
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantNil bool
		wantErr bool
	}{
		{
			"err-nil-opts",
			fields{
				nil,
			},
			args{
				nil,
			},
			true,
			true,
		},
		{
			"err-invalid-SignerOptions",
			fields{
				secp256k1test.CharlottePrivateKey,
			},
			args{
				func() interface{} {
					type NotSignerOptions struct {
						Tx *types.Transaction
					}
					return NotSignerOptions{Tx: &types.Transaction{}}
				}(),
			},
			true,
			true,
		},
		{
			"success-SignerOptions",
			fields{
				secp256k1test.CharlottePrivateKey,
			},
			args{
				SignerOptions{
					Tx:      &types.Transaction{},
					ChainID: big.NewInt(1000),
				},
			},
			false,
			false,
		},
		{
			"success-SignerOptions-chainid-nil",
			fields{
				secp256k1test.CharlottePrivateKey,
			},
			args{
				SignerOptions{
					Tx: &types.Transaction{},
				},
			},
			false,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Signer{
				privateKey: tt.fields.privateKey,
			}
			gotSignedTransaction, err := e.Sign(tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("Signer.Sign() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (gotSignedTransaction == nil) != tt.wantNil {
				t.Errorf("Signer.Sign() = %v, wantErr %v", gotSignedTransaction, tt.wantNil)
				return
			}
		})
	}
}

func Test_validatePrivateKeyType(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		pk crypto.PrivateKey
	}
	tests := []struct {
		name    string
		args    args
		want    *ecdsa.PrivateKey
		wantErr bool
	}{
		{
			"success-secp256k1-charlotte",
			args{
				secp256k1test.CharlottePrivateKey,
			},
			func() *ecdsa.PrivateKey {
				m, _ := ethcrypto.ToECDSA([]byte{0xdf, 0x4b, 0xa9, 0xf6, 0x10, 0x6a, 0xd2, 0x84, 0x64, 0x72, 0xf7, 0x59, 0x47, 0x65, 0x35, 0xe5, 0x5c, 0x58, 0x05, 0xd8, 0x33, 0x7d, 0xf5, 0xa1, 0x1c, 0x3b, 0x13, 0x9f, 0x43, 0x8b, 0x98, 0xb3})
				return m
			}(),
			false,
		},
		{
			"success-secp256k1-charlotte",
			args{
				*secp256k1test.CharlottePrivateKey.(*secp256k1.PrivateKey),
			},
			func() *ecdsa.PrivateKey {
				m, _ := ethcrypto.ToECDSA([]byte{0xdf, 0x4b, 0xa9, 0xf6, 0x10, 0x6a, 0xd2, 0x84, 0x64, 0x72, 0xf7, 0x59, 0x47, 0x65, 0x35, 0xe5, 0x5c, 0x58, 0x05, 0xd8, 0x33, 0x7d, 0xf5, 0xa1, 0x1c, 0x3b, 0x13, 0x9f, 0x43, 0x8b, 0x98, 0xb3})
				return m
			}(),
			false,
		},
		{
			"err-ed25519-charlotte",
			args{
				ed25519test.CharlottePrivateKey,
			},
			nil,
			true,
		},
		{
			"err-nil",
			args{
				nil,
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validatePrivateKeyType(tt.args.pk)
			if (err != nil) != tt.wantErr {
				t.Errorf("validatePrivateKeyType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !assert.Equal(tt.want, got) {
				t.Errorf("validatePrivateKeyType() = %v, want %v", got, tt.want)
			}
		})
	}
}
