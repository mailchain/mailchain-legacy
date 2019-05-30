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

package multikey

import (
	"testing"

	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestPrivateKeyFromHex(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		hex     string
		keyType string
	}
	tests := []struct {
		name    string
		args    args
		want    crypto.PrivateKey
		wantErr bool
	}{
		{
			"ethereum",
			args{
				"01901E63389EF02EAA7C5782E08B40D98FAEF835F28BD144EECF5614A415943F",
				"secp256k1",
			},
			testutil.SofiaPrivateKey,
			false,
		},
		{
			"err",
			args{
				"01901E63389EF02EAA7C5782E08B40D98FAEF835F28BD144EECF5614A415943F",
				"unknown",
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PrivateKeyFromHex(tt.args.hex, tt.args.keyType)
			if (err != nil) != tt.wantErr {
				t.Errorf("PrivateKeyFromHex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(tt.want, got) {
				t.Errorf("PrivateKeyFromHex() = %v, want %v", got, tt.want)
			}
		})
	}
}
