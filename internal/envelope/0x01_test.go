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
	"net/url"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/cipher"
	"github.com/mailchain/mailchain/crypto/cipher/ciphertest"
	"github.com/mailchain/mailchain/crypto/secp256k1/secp256k1test"
	"github.com/mailchain/mailchain/encoding/encodingtest"
	"github.com/mailchain/mailchain/internal/mli"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestZeroX01_URL(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	type fields struct {
		UIBEncryptedLocationHash []byte
		EncryptedHash            []byte
	}
	type args struct {
		decrypter cipher.Decrypter
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *url.URL
		wantErr bool
	}{
		{
			"success",
			fields{
				[]byte("UIBEncryptedLocationHash"),
				[]byte("EncryptedHash"),
			},
			args{
				func() cipher.Decrypter {
					m := ciphertest.NewMockDecrypter(mockCtrl)

					m.EXPECT().Decrypt(cipher.EncryptedContent([]byte("UIBEncryptedLocationHash"))).Return([]byte{0x1, 0x01, 0x76, 0x61, 0x6c, 0x75, 0x65}, nil)
					return m
				}(),
			},
			func() *url.URL {
				u, _ := url.Parse("https://mcx.mx/76616c7565")
				return u
			}(),
			false,
		},
		{
			"err-unknown-loc-code",
			fields{
				[]byte("UIBEncryptedLocationHash"),
				[]byte("EncryptedHash"),
			},
			args{
				func() cipher.Decrypter {
					m := ciphertest.NewMockDecrypter(mockCtrl)

					m.EXPECT().Decrypt(cipher.EncryptedContent([]byte("UIBEncryptedLocationHash"))).Return([]byte{0x1, 0x00, 0x76, 0x61, 0x6c, 0x75, 0x65}, nil)
					return m
				}(),
			},
			nil,
			true,
		},
		{
			"err-invalid-uint64-bytes",
			fields{
				[]byte("UIBEncryptedLocationHash"),
				[]byte("EncryptedHash"),
			},
			args{
				func() cipher.Decrypter {
					m := ciphertest.NewMockDecrypter(mockCtrl)

					m.EXPECT().Decrypt(cipher.EncryptedContent([]byte("UIBEncryptedLocationHash"))).Return([]byte{0x1}, nil)
					return m
				}(),
			},
			nil,
			true,
		},
		{
			"err-decrypt",
			fields{
				[]byte("UIBEncryptedLocationHash"),
				[]byte("EncryptedHash"),
			},
			args{
				func() cipher.Decrypter {
					m := ciphertest.NewMockDecrypter(mockCtrl)

					m.EXPECT().Decrypt(cipher.EncryptedContent([]byte("UIBEncryptedLocationHash"))).Return(nil, errors.Errorf("failed"))
					return m
				}(),
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &ZeroX01{
				UIBEncryptedLocationHash: tt.fields.UIBEncryptedLocationHash,
				EncryptedHash:            tt.fields.EncryptedHash,
			}
			got, err := d.URL(tt.args.decrypter)
			if (err != nil) != tt.wantErr {
				t.Errorf("ZeroX01.URL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ZeroX01.URL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestZeroX01_ContentsHash(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type fields struct {
		UIBEncryptedLocationHash []byte
		EncryptedHash            []byte
	}
	type args struct {
		decrypter cipher.Decrypter
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
			fields{
				[]byte("UIBEncryptedLocationHash"),
				[]byte("EncryptedHash"),
			},
			args{
				func() cipher.Decrypter {
					m := ciphertest.NewMockDecrypter(mockCtrl)

					m.EXPECT().Decrypt(cipher.EncryptedContent([]byte("UIBEncryptedLocationHash"))).Return([]byte{0x1, 0x01, 0x76, 0x61, 0x6c, 0x75, 0x65}, nil)
					return m
				}(),
			},
			[]byte{0x76, 0x61, 0x6c, 0x75, 0x65},
			false,
		},
		{
			"err-decrypt",
			fields{
				[]byte("UIBEncryptedLocationHash"),
				[]byte("EncryptedHash"),
			},
			args{
				func() cipher.Decrypter {
					m := ciphertest.NewMockDecrypter(mockCtrl)

					m.EXPECT().Decrypt(cipher.EncryptedContent([]byte("UIBEncryptedLocationHash"))).Return(nil, errors.Errorf("failed"))
					return m
				}(),
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &ZeroX01{
				UIBEncryptedLocationHash: tt.fields.UIBEncryptedLocationHash,
				EncryptedHash:            tt.fields.EncryptedHash,
			}
			got, err := d.ContentsHash(tt.args.decrypter)
			if (err != nil) != tt.wantErr {
				t.Errorf("ZeroX01.ContentsHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("ZeroX01.ContentsHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestZeroX01_IntegrityHash(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type fields struct {
		UIBEncryptedLocationHash []byte
		EncryptedHash            []byte
	}
	type args struct {
		decrypter cipher.Decrypter
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
			fields{
				[]byte("UIBEncryptedLocationHash"),
				[]byte("EncryptedHash"),
			},
			args{
				func() cipher.Decrypter {
					m := ciphertest.NewMockDecrypter(mockCtrl)
					return m
				}(),
			},
			[]byte("EncryptedHash"),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &ZeroX01{
				UIBEncryptedLocationHash: tt.fields.UIBEncryptedLocationHash,
				EncryptedHash:            tt.fields.EncryptedHash,
			}
			got, err := d.IntegrityHash(tt.args.decrypter)
			if (err != nil) != tt.wantErr {
				t.Errorf("ZeroX01.IntegrityHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ZeroX01.IntegrityHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestZeroX01_Valid(t *testing.T) {
	type fields struct {
		UIBEncryptedLocationHash []byte
		EncryptedHash            []byte
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			"success",
			fields{
				[]byte("UIBEncryptedLocationHash"),
				nil,
			},
			false,
		},
		{
			"err-no-EncryptedLocationHash",
			fields{
				nil,
				nil,
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &ZeroX01{
				UIBEncryptedLocationHash: tt.fields.UIBEncryptedLocationHash,
				EncryptedHash:            tt.fields.EncryptedHash,
			}
			if err := d.Valid(); (err != nil) != tt.wantErr {
				t.Errorf("ZeroX01.Valid() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewZeroX01(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type args struct {
		encrypter cipher.Encrypter
		pubkey    crypto.PublicKey
		opts      *CreateOpts
	}
	tests := []struct {
		name    string
		args    args
		want    *ZeroX01
		wantErr bool
	}{
		{
			"success",
			args{
				func() cipher.Encrypter {
					m := ciphertest.NewMockEncrypter(mockCtrl)
					m.EXPECT().Encrypt(gomock.Any()).Return([]byte("encrypted"), nil)
					return m
				}(),
				secp256k1test.BobPublicKey,
				&CreateOpts{
					Location:      mli.Mailchain,
					DecryptedHash: encodingtest.MustDecodeHex("2c8432ca28ce929b86a47f2d40413d161f591f8985229060491573d83f82f292f4dc68f918446332837aa57cd5145235cc40702d962cbb53ac27fb2246fb6cba"),
					Resource:      "2c8432ca28ce929b86a47f2d40413d161f591f8985229060491573d83f82f292f4dc68f918446332837aa57cd5145235cc40702d962cbb53ac27fb2246fb6cba",
				},
			},
			&ZeroX01{
				UIBEncryptedLocationHash: []uint8{0x65, 0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x65, 0x64},
				EncryptedHash:            []uint8(nil),
			},
			false,
		},
		{
			"success-encrypted-hash",
			args{
				func() cipher.Encrypter {
					m := ciphertest.NewMockEncrypter(mockCtrl)
					m.EXPECT().Encrypt(gomock.Any()).Return([]byte("encrypted"), nil)
					return m
				}(),
				secp256k1test.BobPublicKey,
				&CreateOpts{
					EncryptedHash: encodingtest.MustDecodeHex("220455078214"),
					Location:      mli.Mailchain,
					DecryptedHash: encodingtest.MustDecodeHex("2c8432ca28ce929b86a47f2d40413d161f591f8985229060491573d83f82f292f4dc68f918446332837aa57cd5145235cc40702d962cbb53ac27fb2246fb6cba"),
					Resource:      "2c8432ca28ce929b86a47f2d40413d161f591f8985229060491573d83f82f292f4dc68f918446332837aa57cd5145235cc40702d962cbb53ac27fb2246fb6cba",
				},
			},
			&ZeroX01{
				UIBEncryptedLocationHash: []uint8{0x65, 0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x65, 0x64},
				EncryptedHash:            []uint8{0x22, 0x4, 0x55, 0x7, 0x82, 0x14},
			},
			false,
		},
		{
			"err-encrypt",
			args{
				func() cipher.Encrypter {
					m := ciphertest.NewMockEncrypter(mockCtrl)
					m.EXPECT().Encrypt(gomock.Any()).Return(nil, errors.Errorf("failed"))
					return m
				}(),
				secp256k1test.BobPublicKey,
				&CreateOpts{
					Location:      mli.Mailchain,
					Resource:      "2c8432ca28ce929b86a47f2d40413d161f591f8985229060491573d83f82f292f4dc68f918446332837aa57cd5145235cc40702d962cbb53ac27fb2246fb6cba",
					DecryptedHash: encodingtest.MustDecodeHex("2c8432ca28ce929b86a47f2d40413d161f591f8985229060491573d83f82f292f4dc68f918446332837aa57cd5145235cc40702d962cbb53ac27fb2246fb6cba"),
				},
			},
			nil,
			true,
		},
		{
			"err-missing-decrypted-hash",
			args{
				func() cipher.Encrypter {
					m := ciphertest.NewMockEncrypter(mockCtrl)
					return m
				}(),
				secp256k1test.BobPublicKey,
				&CreateOpts{
					Location: mli.Mailchain,
				},
			},
			nil,
			true,
		},
		{
			"err-missing-location",
			args{
				func() cipher.Encrypter {
					m := ciphertest.NewMockEncrypter(mockCtrl)
					return m
				}(),
				secp256k1test.BobPublicKey,
				&CreateOpts{},
			},
			nil,
			true,
		},
		{
			"err-missing-resource",
			args{
				func() cipher.Encrypter {
					m := ciphertest.NewMockEncrypter(mockCtrl)
					return m
				}(),
				secp256k1test.BobPublicKey,
				&CreateOpts{
					Location:      mli.Mailchain,
					DecryptedHash: encodingtest.MustDecodeHex("2c8432ca28ce929b86a47f2d40413d161f591f8985229060491573d83f82f292f4dc68f918446332837aa57cd5145235cc40702d962cbb53ac27fb2246fb6cba"),
				},
			},
			nil,
			true,
		},
		{
			"err-invalid-resource",
			args{
				func() cipher.Encrypter {
					m := ciphertest.NewMockEncrypter(mockCtrl)
					return m
				}(),
				secp256k1test.BobPublicKey,
				&CreateOpts{
					Location:      mli.Mailchain,
					Resource:      "invalid",
					DecryptedHash: encodingtest.MustDecodeHex("2c8432ca28ce929b86a47f2d40413d161f591f8985229060491573d83f82f292f4dc68f918446332837aa57cd5145235cc40702d962cbb53ac27fb2246fb6cba"),
				},
			},
			nil,
			true,
		},
		{
			"err-resource-not-match-hash",
			args{
				func() cipher.Encrypter {
					m := ciphertest.NewMockEncrypter(mockCtrl)
					return m
				}(),
				secp256k1test.BobPublicKey,
				&CreateOpts{
					Location:      mli.Mailchain,
					Resource:      "2c8432ca",
					DecryptedHash: encodingtest.MustDecodeHex("2c8432ca28ce929b86a47f2d40413d161f591f8985229060491573d83f82f292f4dc68f918446332837aa57cd5145235cc40702d962cbb53ac27fb2246fb6cba"),
				},
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewZeroX01(tt.args.encrypter, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewZeroX01() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("NewZeroX01() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestZeroX01_DecrypterKind(t *testing.T) {
	type fields struct {
		UIBEncryptedLocationHash []byte
	}
	tests := []struct {
		name    string
		fields  fields
		want    byte
		wantErr bool
	}{
		{
			"success",
			fields{
				encodingtest.MustDecodeHexZeroX("0x2ee10c59024c836d7ca12470b5ac74673002127ddedadbc6fc4375a8c086b650060ede199f603a158bc7884a903eadf97a2dd0fbe69ac81c216830f94e56b847d924b51a7d8227c80714219e6821a51bc7cba922f291a47bdffe29e7c3f67ad908ff377bfcc0b603007ead4bfd87ff0acc272528ca03d6381e6d0e1e2c5dfd24d521"),
			},
			cipher.AES256CBC,
			false,
		},
		{
			"err-empty",
			fields{
				[]byte{},
			},
			0x0,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &ZeroX01{
				UIBEncryptedLocationHash: tt.fields.UIBEncryptedLocationHash,
			}
			got, err := x.DecrypterKind()
			if (err != nil) != tt.wantErr {
				t.Errorf("ZeroX01.DecrypterKind() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ZeroX01.DecrypterKind() = %v, want %v", got, tt.want)
			}
		})
	}
}
