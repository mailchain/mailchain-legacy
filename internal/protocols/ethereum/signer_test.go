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
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/ed25519/ed25519test"
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
	type args struct {
		pk crypto.PrivateKey
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"success-secp256k1-charlotte",
			args{
				secp256k1test.CharlottePrivateKey,
			},
			false,
		},
		{
			"err-ed25519-charlotte",
			args{
				ed25519test.CharlottePrivateKey,
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
			err := validatePrivateKeyType(tt.args.pk)
			if (err != nil) != tt.wantErr {
				t.Errorf("validatePrivateKeyType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
