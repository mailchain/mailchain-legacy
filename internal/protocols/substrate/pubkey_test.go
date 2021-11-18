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
	"context"
	"testing"

	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/sr25519/sr25519test"
	"github.com/mailchain/mailchain/encoding/encodingtest"
	"github.com/mr-tron/base58"
	"github.com/stretchr/testify/assert"
)

func TestPublicKeyFinder_PublicKeyFromAddress(t *testing.T) {
	type args struct {
		ctx      context.Context
		protocol string
		network  string
		address  []byte
	}
	tests := []struct {
		name    string
		args    args
		want    crypto.PublicKey
		wantErr bool
	}{
		{
			"success",
			args{
				context.Background(),
				"substrate",
				"edgeware-mainnet",
				encodingtest.MustDecodeBase58("kWCKFhkg5m5XmGz2pbcUYhJK7RAf6jLj1wLM5J5junonVTH"),
			},
			sr25519test.BobPublicKey,
			false,
		},
		{
			"err-invalid-length",
			args{
				context.Background(),
				"substrate",
				"beresheet",
				func() []byte {
					num, err := base58.Decode("5CLmNK8f16nagFeF2h3iNeeChaxPiAsJu7piNYJgdPpmaRzPD")
					if err != nil {
						t.Errorf("got error %s\n", err)
						t.FailNow()
					}
					return num
				}(),
			},
			nil,
			true,
		},
		{
			"err-protocol",
			args{
				context.Background(),
				"invalid",
				"beresheet",
				func() []byte {
					num, err := base58.Decode("5CLmNK8f16nagFeF2h3iNeeChaxPiAsJu7piNYJgdPpmaRzPD")
					if err != nil {
						t.Errorf("got error %s\n", err)
						t.FailNow()
					}
					return num
				}(),
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkf := &PublicKeyFinder{}
			got, err := pkf.PublicKeyFromAddress(tt.args.ctx, tt.args.protocol, tt.args.network, tt.args.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("PublicKeyFinder.PublicKeyFromAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil {
				if !assert.Equal(t, tt.want, got) {
					t.Errorf("PublicKeyFinder.PublicKeyFromAddress() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestNewPublicKeyFinder(t *testing.T) {
	tests := []struct {
		name string
		want *PublicKeyFinder
	}{
		{
			"success",
			&PublicKeyFinder{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPublicKeyFinder(); !assert.Equal(t, tt.want, got) {
				t.Errorf("NewPublicKeyFinder() = %v, want %v", got, tt.want)
			}
		})
	}
}
