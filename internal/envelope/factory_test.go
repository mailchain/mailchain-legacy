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

package envelope

import (
	"encoding/hex"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/cipher"
	"github.com/mailchain/mailchain/crypto/cipher/ciphertest"
	"github.com/mailchain/mailchain/crypto/secp256k1/secp256k1test"
	"github.com/mailchain/mailchain/internal/mli"
	"github.com/stretchr/testify/assert"
)

func TestNewEnvelope(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	opts := func(envelopeKind byte) []CreateOptionsBuilder {
		locOpt, err := WithMessageLocationIdentifier(mli.Mailchain)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		encodedMsg, _ := hex.DecodeString("220455078214")
		return []CreateOptionsBuilder{
			WithKind(envelopeKind),
			WithURL("https://domain.com/220455078214"),
			WithResource("220455078214"),
			WithDecryptedHash(encodedMsg),
			locOpt,
		}

	}
	type args struct {
		encrypter cipher.Encrypter
		pubkey    crypto.PublicKey
		o         []CreateOptionsBuilder
	}
	tests := []struct {
		name    string
		args    args
		want    Data
		wantErr bool
	}{
		{
			"invalid",
			args{
				func() cipher.Encrypter {
					m := ciphertest.NewMockEncrypter(mockCtrl)
					return m
				}(),
				secp256k1test.BobPublicKey,
				opts(0x00),
			},
			nil,
			true,
		},
		{
			"0x01",
			args{
				func() cipher.Encrypter {
					m := ciphertest.NewMockEncrypter(mockCtrl)
					m.EXPECT().Encrypt(gomock.Any()).Return([]byte("encrypted"), nil)
					return m
				}(),
				secp256k1test.BobPublicKey,
				opts(Kind0x01),
			},
			&ZeroX01{},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewEnvelope(tt.args.encrypter, tt.args.o)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewEnvelope() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.IsType(t, tt.want, got) {
				t.Errorf("NewEnvelope() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseEvelope(t *testing.T) {
	type args struct {
		envelope string
	}
	tests := []struct {
		name    string
		args    args
		want    byte
		wantErr bool
	}{
		{
			"0x01",
			args{
				KindString0x01,
			},
			Kind0x01,
			false,
		},
		{
			"0x02",
			args{
				KindString0x02,
			},
			Kind0x02,
			false,
		},
		{
			"0x50",
			args{
				KindString0x50,
			},
			Kind0x50,
			false,
		},
		{
			"empty-envelope",
			args{
				"",
			},
			0x0,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseEnvelope(tt.args.envelope)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseEnvelope() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != got {
				t.Errorf("ParseEnvelope() = %v, want %v", got, tt.want)
			}
		})
	}
}
