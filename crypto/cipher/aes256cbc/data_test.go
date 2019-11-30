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

package aes256cbc

import (
	"testing"

	"github.com/mailchain/mailchain/internal/encoding/encodingtest"
)

func Test_encryptedData_verify(t *testing.T) {
	type fields struct {
		InitializationVector      []byte
		EphemeralPublicKey        []byte
		Ciphertext                []byte
		MessageAuthenticationCode []byte
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			"success",
			fields{
				encodingtest.MustDecodeHex("2c8432ca28ce929b86a47f2d40413d16"),
				encodingtest.MustDecodeHex("2c8432ca28ce929b86a47f2d40413d161f591f8985229060491573d83f82f292f4dc68f918446332837aa57cd5145235cc40702d962cbb53ac27fb2246fb6cbadc"),
				encodingtest.MustDecodeHex("2c8432ca"),
				encodingtest.MustDecodeHex("2c8432ca28ce929b86a47f2d40413d161f591f8985229060491573d83f82f292"),
			},
			false,
		},
		{
			"err-iv",
			fields{
				encodingtest.MustDecodeHex(""),
				encodingtest.MustDecodeHex("2c8432ca28ce929b86a47f2d40413d161f591f8985229060491573d83f82f292f4dc68f918446332837aa57cd5145235cc40702d962cbb53ac27fb2246fb6cbadc"),
				encodingtest.MustDecodeHex("2c8432ca"),
				encodingtest.MustDecodeHex("2c8432ca28ce929b86a47f2d40413d161f591f8985229060491573d83f82f292"),
			},
			true,
		},
		{
			"err-ethemeral-pk",
			fields{
				encodingtest.MustDecodeHex("2c8432ca28ce929b86a47f2d40413d16"),
				encodingtest.MustDecodeHex(""),
				encodingtest.MustDecodeHex("2c8432ca"),
				encodingtest.MustDecodeHex("2c8432ca28ce929b86a47f2d40413d161f591f8985229060491573d83f82f292"),
			},
			true,
		},
		{
			"err-cipher-text",
			fields{
				encodingtest.MustDecodeHex("2c8432ca28ce929b86a47f2d40413d16"),
				encodingtest.MustDecodeHex("2c8432ca28ce929b86a47f2d40413d161f591f8985229060491573d83f82f292f4dc68f918446332837aa57cd5145235cc40702d962cbb53ac27fb2246fb6cbadc"),
				encodingtest.MustDecodeHex(""),
				encodingtest.MustDecodeHex("2c8432ca28ce929b86a47f2d40413d161f591f8985229060491573d83f82f292"),
			},
			true,
		},
		{
			"err-mac",
			fields{
				encodingtest.MustDecodeHex("2c8432ca28ce929b86a47f2d40413d16"),
				encodingtest.MustDecodeHex("2c8432ca28ce929b86a47f2d40413d161f591f8985229060491573d83f82f292f4dc68f918446332837aa57cd5145235cc40702d962cbb53ac27fb2246fb6cbadc"),
				encodingtest.MustDecodeHex("2c8432ca"),
				encodingtest.MustDecodeHex(""),
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &encryptedData{
				InitializationVector:      tt.fields.InitializationVector,
				EphemeralPublicKey:        tt.fields.EphemeralPublicKey,
				Ciphertext:                tt.fields.Ciphertext,
				MessageAuthenticationCode: tt.fields.MessageAuthenticationCode,
			}
			if err := e.verify(); (err != nil) != tt.wantErr {
				t.Errorf("encryptedData.verify() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
