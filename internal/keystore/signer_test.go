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

package keystore

import (
	"testing"

	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/secp256k1/secp256k1test"
	"github.com/mailchain/mailchain/internal/mailbox/signer"
	"github.com/mailchain/mailchain/internal/protocols"
	"github.com/mailchain/mailchain/internal/protocols/ethereum"
	"github.com/stretchr/testify/assert"
)

func TestSigner(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		chain string
		pk    crypto.PrivateKey
	}
	tests := []struct {
		name    string
		args    args
		want    signer.Signer
		wantErr bool
	}{
		{
			"ethereum",
			args{
				protocols.Ethereum,
				secp256k1test.CharlottePrivateKey,
			},
			func() signer.Signer {
				m, _ := ethereum.NewSigner(secp256k1test.CharlottePrivateKey)
				return m
			}(),
			false,
		},
		{
			"err",
			args{
				"invalid",
				secp256k1test.CharlottePrivateKey,
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Signer(tt.args.chain, tt.args.pk)
			if (err != nil) != tt.wantErr {
				t.Errorf("Signer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(tt.want, got) {
				t.Errorf("Signer() = %v, want %v", got, tt.want)
			}
		})
	}
}
