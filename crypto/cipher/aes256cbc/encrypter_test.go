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

	ethcypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/mailchain/mailchain/crypto/secp256k1"
	"github.com/mailchain/mailchain/crypto/secp256k1/secp256k1test"
	"github.com/mailchain/mailchain/encoding/encodingtest"
	"github.com/stretchr/testify/assert"
)

func Test_encryptCBC(t *testing.T) {
	type args struct {
		data []byte
		iv   []byte
		key  []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			"success-short-text",
			args{
				[]byte("Hi Tim"),
				encodingtest.MustDecodeHex("05050505050505050505050505050505"),
				encodingtest.MustDecodeHex("af0ad81e7d9194721d6c26f6c1f2a2b7fd06e2c99c4f5deefe59fb93936c981e"),
			},
			encodingtest.MustDecodeHex("747ef78a32eb582d325a634e4acffd61"),
			false,
		},
		{
			"success-medium-text",
			args{
				[]byte("Hi Tim, this is a much longer message to make sure there are no problems"),
				encodingtest.MustDecodeHex("05050505050505050505050505050505"),
				encodingtest.MustDecodeHex("af0ad81e7d9194721d6c26f6c1f2a2b7fd06e2c99c4f5deefe59fb93936c981e"),
			},
			encodingtest.MustDecodeHex("2ec66aac453ff543f47830d4b8cbc68d9965bf7c6bb69724fd4de26d41001256dfa6f7f0b3956ce21d4717caf75b0c2ad753852f216df6cfbcda4911619c5fc34798a19f81adff902c1ad906ab0edaec"),
			false,
		},
		{
			"err-iv",
			args{
				[]byte("Hi Tim, this is a much longer message to make sure there are no problems"),
				encodingtest.MustDecodeHex("0505"),
				encodingtest.MustDecodeHex("af0ad81e7d9194721d6c26f6c1f2a2b7fd06e2c99c4f5deefe59fb93936c981e"),
			},
			nil,
			true,
		},
		{
			"err-key",
			args{
				[]byte("Hi Tim, this is a much longer message to make sure there are no problems"),
				encodingtest.MustDecodeHex("05050505050505050505050505050505"),
				encodingtest.MustDecodeHex("af"),
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := encryptCBC(tt.args.data, tt.args.iv, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("encryptCBC() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("encryptCBC() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_encrypt(t *testing.T) {
	type args struct {
		ephemeralPrivateKey *ecies.PrivateKey
		pub                 *ecies.PublicKey
		input               []byte
		iv                  []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *encryptedData
		wantErr bool
	}{
		{
			"success",
			args{
				func() *ecies.PrivateKey {
					tmpEphemeralPrivateKey, err := ethcypto.HexToECDSA("0404040404040404040404040404040404040404040404040404040404040404")
					if err != nil {
						t.Fatal(err)
					}
					return ecies.ImportECDSA(tmpEphemeralPrivateKey)
				}(),
				func() *ecies.PublicKey {
					tp, ok := secp256k1test.AlicePublicKey.(*secp256k1.PublicKey)
					if !ok {
						t.Error("failed to cast")
					}
					pub, err := tp.ECIES()
					if err != nil {
						t.Error(err)
					}
					return pub
				}(),
				[]byte("Hi Tim, this is a much longer message to make sure there are no problems"),
				encodingtest.MustDecodeHex("05050505050505050505050505050505"),
			},
			&encryptedData{
				Ciphertext:                encodingtest.MustDecodeHex("2ec66aac453ff543f47830d4b8cbc68d9965bf7c6bb69724fd4de26d41001256dfa6f7f0b3956ce21d4717caf75b0c2ad753852f216df6cfbcda4911619c5fc34798a19f81adff902c1ad906ab0edaec"),
				EphemeralPublicKey:        encodingtest.MustDecodeHex("04462779ad4aad39514614751a71085f2f10e1c7a593e4e030efb5b8721ce55b0b199c07969f5442000bea455d72ae826a86bfac9089cb18152ed756ebb2a596f5"),
				InitializationVector:      encodingtest.MustDecodeHex("05050505050505050505050505050505"),
				MessageAuthenticationCode: encodingtest.MustDecodeHex("4367ae8a54b65f99e4f2fd315ba65bf85e1138967a7bea451faf80f75cdf3404"),
			},
			false,
		},
		{
			"err-encryptCBC-invalid-iv",
			args{
				func() *ecies.PrivateKey {
					tmpEphemeralPrivateKey, err := ethcypto.HexToECDSA("0404040404040404040404040404040404040404040404040404040404040404")
					if err != nil {
						t.Fatal(err)
					}
					return ecies.ImportECDSA(tmpEphemeralPrivateKey)
				}(),
				func() *ecies.PublicKey {
					tp, ok := secp256k1test.AlicePublicKey.(*secp256k1.PublicKey)
					if !ok {
						t.Error("failed to cast")
					}
					pub, err := tp.ECIES()
					if err != nil {
						t.Error(err)
					}
					return pub
				}(),
				[]byte("Hi Tim, this is a much longer message to make sure there are no problems"),
				encodingtest.MustDecodeHex("0505"),
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := encrypt(tt.args.ephemeralPrivateKey, tt.args.pub, tt.args.input, tt.args.iv)
			if (err != nil) != tt.wantErr {
				t.Errorf("encrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("encrypt() = %v, want %v", got, tt.want)
			}
		})
	}
}
