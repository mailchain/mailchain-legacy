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

	"github.com/mailchain/mailchain/internal/chains/ethereum"
	"github.com/mailchain/mailchain/internal/crypto/keys"
	"github.com/mailchain/mailchain/internal/encoding"
	"github.com/mailchain/mailchain/internal/mailbox"
	"github.com/mailchain/mailchain/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestSigner(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		chain string
		pk    keys.PrivateKey
	}
	tests := []struct {
		name    string
		args    args
		want    mailbox.Signer
		wantErr bool
	}{
		{
			"ethereum",
			args{
				encoding.Ethereum,
				testutil.CharlottePrivateKey,
			},
			ethereum.NewSigner(testutil.CharlottePrivateKey),
			false,
		},
		{
			"err",
			args{
				"invalid",
				testutil.CharlottePrivateKey,
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
