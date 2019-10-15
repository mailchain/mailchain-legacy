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

	"github.com/mailchain/mailchain/internal/testutil"
	"github.com/mr-tron/base58"
	"github.com/stretchr/testify/assert"
)

func TestPublicKeyFinder_PublicKeyFromAddress(t *testing.T) {
	assert := assert.New(t)
	type fields struct {
		supportedNetworks []string
	}
	type args struct {
		ctx      context.Context
		protocol string
		network  string
		address  []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		{
			"success",
			fields{},
			args{
				context.Background(),
				"substrate",
				"testnet",
				func() []byte {
					num, err := base58.Decode("5CLmNK8f16nagFeF2h3iNeeChaxPiAsJu7piNYJgdPpmaRzP")
					if err != nil {
						t.Errorf("got error %s\n", err)
						t.FailNow()
					}
					return num
				}(),
			},
			testutil.MustHexDecodeString("0c3fbef5c06307444e8078036c217b2907f2459e906ff0f1a670986743f2494f"),
			false,
		},
		{
			"err-invalid-length",
			fields{},
			args{
				context.Background(),
				"substrate",
				"testnet",
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
			fields{},
			args{
				context.Background(),
				"invalid",
				"testnet",
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
			pkf := &PublicKeyFinder{
				supportedNetworks: tt.fields.supportedNetworks,
			}
			got, err := pkf.PublicKeyFromAddress(tt.args.ctx, tt.args.protocol, tt.args.network, tt.args.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("PublicKeyFinder.PublicKeyFromAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(tt.want, got) {
				t.Errorf("PublicKeyFinder.PublicKeyFromAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}
