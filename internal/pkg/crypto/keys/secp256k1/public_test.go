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

package secp256k1_test

import (
	"encoding/hex"
	"reflect"
	"testing"

	"github.com/mailchain/mailchain/internal/pkg/crypto/keys/secp256k1"
	"github.com/mailchain/mailchain/internal/pkg/testutil"
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
				"0x" + hex.EncodeToString(testutil.CharlottePublicKey.Bytes()),
			},
			testutil.CharlottePublicKey.Bytes(),
			false,
		},
		{
			"success-no-prefix",
			args{
				hex.EncodeToString(testutil.CharlottePublicKey.Bytes()),
			},
			testutil.CharlottePublicKey.Bytes(),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := secp256k1.PublicKeyFromHex(tt.args.hex)
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
