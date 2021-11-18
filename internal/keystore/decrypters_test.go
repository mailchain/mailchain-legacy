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

	"github.com/golang/mock/gomock"
	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/cipher"
	"github.com/mailchain/mailchain/crypto/cipher/aes256cbc"
	"github.com/mailchain/mailchain/crypto/cipher/nacl"
	"github.com/mailchain/mailchain/crypto/cipher/noop"
	"github.com/mailchain/mailchain/crypto/cryptotest"
	"github.com/mailchain/mailchain/crypto/ed25519/ed25519test"
	"github.com/mailchain/mailchain/crypto/secp256k1/secp256k1test"
	"github.com/stretchr/testify/assert"
)

func TestDecrypter(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	type args struct {
		cipherType byte
		pk         crypto.PrivateKey
	}
	tests := []struct {
		name    string
		args    args
		want    cipher.Decrypter
		wantErr bool
	}{
		{
			"aes256cbc",
			args{
				cipher.AES256CBC,
				secp256k1test.BobPrivateKey,
			},
			func() cipher.Decrypter {
				m, _ := aes256cbc.NewDecrypter(secp256k1test.BobPrivateKey)
				return m
			}(),
			false,
		},
		{
			"nacl",
			args{
				cipher.NACLECDH,
				ed25519test.BobPrivateKey,
			},
			func() cipher.Decrypter {
				m, _ := nacl.NewDecrypter(ed25519test.BobPrivateKey)
				return m
			}(),
			false,
		},
		{
			"noop",
			args{
				cipher.NoOperation,
				ed25519test.BobPrivateKey,
			},
			func() cipher.Decrypter {
				m := noop.NewDecrypter()
				return m
			}(),
			false,
		},
		{
			"err-invalid-key-type",
			args{
				0xFF,
				secp256k1test.BobPrivateKey,
			},
			nil,
			true,
		},
		{
			"err-invalid-nacl-key-type",
			args{
				cipher.NACLECDH,
				cryptotest.NewMockPrivateKey(mockCtrl),
			},
			(*nacl.Decrypter)(nil),
			true,
		},
		{
			"err-invalid-aes256cbc-key-type",
			args{
				cipher.AES256CBC,
				ed25519test.BobPrivateKey,
			},
			(*aes256cbc.Decrypter)(nil),
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Decrypter(tt.args.cipherType, tt.args.pk)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decrypter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("Decrypter() = %v, want %v", got, tt.want)
			}
		})
	}
}
