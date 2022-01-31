// Copyright 2022 Mailchain Ltd.
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

func TestZeroX02_URL(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	type fields struct {
		UIBEncryptedLocationHash []byte
		DecryptedHash            []byte
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
				[]byte("DecryptedHash"),
			},
			args{
				func() cipher.Decrypter {
					m := ciphertest.NewMockDecrypter(mockCtrl)

					m.EXPECT().Decrypt(cipher.EncryptedContent([]byte("UIBEncryptedLocationHash"))).Return(encodingtest.MustDecodeHex("010201551220497b22d4e86a3caa9f5baa24435a99ac1154094a0b9302b9bcd9d6544d6efbe9"), nil)
					return m
				}(),
			},
			func() *url.URL {
				u, _ := url.Parse("https://ipfs.io/ipfs/bafkreicjpmrnj2dkhsvj6w5kerbvvgnmcfkassqlsmbltpgz2zke23x35e")
				return u
			}(),
			false,
		},
		{
			"err-cid",
			fields{
				[]byte("UIBEncryptedLocationHash"),
				[]byte("DecryptedHash"),
			},
			args{
				func() cipher.Decrypter {
					m := ciphertest.NewMockDecrypter(mockCtrl)

					m.EXPECT().Decrypt(cipher.EncryptedContent([]byte("UIBEncryptedLocationHash"))).Return(encodingtest.MustDecodeHex("0102f4551220497b22d4e86a3caa9f5baa24435a99ac1154094a0b9302b9bcd9d6544d6efbe9"), nil)
					return m
				}(),
			},
			nil,
			true,
		},
		{
			"err-unknown-loc-code",
			fields{
				[]byte("UIBEncryptedLocationHash"),
				[]byte("DecryptedHash"),
			},
			args{
				func() cipher.Decrypter {
					m := ciphertest.NewMockDecrypter(mockCtrl)

					m.EXPECT().Decrypt(cipher.EncryptedContent([]byte("UIBEncryptedLocationHash"))).Return(encodingtest.MustDecodeHex("010001551220497b22d4e86a3caa9f5baa24435a99ac1154094a0b9302b9bcd9d6544d6efbe9"), nil)
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
				[]byte("DecryptedHash"),
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
				[]byte("DecryptedHash"),
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
			d := &ZeroX02{
				UIBEncryptedLocationHash: tt.fields.UIBEncryptedLocationHash,
				DecryptedHash:            tt.fields.DecryptedHash,
			}
			got, err := d.URL(tt.args.decrypter)
			if (err != nil) != tt.wantErr {
				t.Errorf("ZeroX02.URL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ZeroX02.URL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestZeroX02_ContentsHash(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type fields struct {
		UIBEncryptedLocationHash []byte
		DecryptedHash            []byte
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
				[]byte("DecryptedHash"),
			},
			args{
				nil,
			},
			[]byte("DecryptedHash"),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &ZeroX02{
				UIBEncryptedLocationHash: tt.fields.UIBEncryptedLocationHash,
				DecryptedHash:            tt.fields.DecryptedHash,
			}
			got, err := d.ContentsHash(tt.args.decrypter)
			if (err != nil) != tt.wantErr {
				t.Errorf("ZeroX02.ContentsHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("ZeroX02.ContentsHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestZeroX02_IntegrityHash(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type fields struct {
		UIBEncryptedLocationHash []byte
		DecryptedHash            []byte
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
				[]byte("DecryptedHash"),
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
				[]byte("DecryptedHash"),
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
			d := &ZeroX02{
				UIBEncryptedLocationHash: tt.fields.UIBEncryptedLocationHash,
				DecryptedHash:            tt.fields.DecryptedHash,
			}
			got, err := d.IntegrityHash(tt.args.decrypter)
			if (err != nil) != tt.wantErr {
				t.Errorf("ZeroX02.IntegrityHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("ZeroX02.IntegrityHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestZeroX02_Valid(t *testing.T) {
	type fields struct {
		UIBEncryptedLocationHash []byte
		DecryptedHash            []byte
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
				[]byte("DecryptedHash"),
			},
			false,
		},
		{
			"err-no-decrypted-hash",
			fields{
				[]byte("UIBEncryptedLocationHash"),
				nil,
			},
			true,
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
			d := &ZeroX02{
				UIBEncryptedLocationHash: tt.fields.UIBEncryptedLocationHash,
				DecryptedHash:            tt.fields.DecryptedHash,
			}
			if err := d.Valid(); (err != nil) != tt.wantErr {
				t.Errorf("ZeroX02.Valid() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewZeroX02(t *testing.T) {
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
		want    *ZeroX02
		wantErr bool
	}{
		{
			"success",
			args{
				func() cipher.Encrypter {
					m := ciphertest.NewMockEncrypter(mockCtrl)
					m.EXPECT().Encrypt(gomock.Any()).Return([]byte{0x2a, 0xe1, 0x2, 0xa3, 0xaa, 0xe8, 0xec, 0xaf, 0xd5, 0xe3, 0xb6, 0xbc, 0x8c, 0xef, 0x3b, 0x6a, 0x63, 0x79, 0x27, 0x5b, 0x89, 0xed, 0x78, 0x3c, 0xba, 0x9, 0xe4, 0xa0, 0xbc, 0x43, 0xba, 0x45, 0xd6, 0x3c, 0xc1, 0x39, 0x13, 0xd3, 0x7a, 0x16, 0x15, 0x37, 0x7d, 0x92, 0x3d, 0x47, 0x3a, 0x63, 0xef, 0x7c, 0x7a, 0xea, 0x4a, 0x1, 0xe1, 0x31, 0x41, 0xbc, 0xa7, 0x6b, 0x6b, 0x5b, 0xbc, 0x66, 0xda, 0x10, 0xdc, 0xb5, 0x26, 0x19, 0x4a, 0xe7, 0x91, 0x24, 0x34, 0x76, 0x1a, 0x43, 0xa, 0x6, 0xf9, 0x4e, 0xae, 0x9e, 0x1e, 0x77, 0xcd, 0x77, 0xdf, 0x91, 0x73, 0xa2, 0xf9, 0x42, 0xb7, 0x2d, 0x67, 0xe1, 0xfd, 0xef, 0xd4, 0x36, 0x54, 0xb5, 0xb, 0xfe, 0x3a, 0xd1, 0xbd, 0x80, 0x3, 0xf3, 0xdf}, nil)
					return m
				}(),
				secp256k1test.BobPublicKey,
				&CreateOpts{
					Location:          mli.IPFS,
					DecryptedHash:     encodingtest.MustDecodeHex("2c8432ca28ce929b86a47f2d40413d161f591f8985229060491573d83f82f292f4dc68f918446332837aa57cd5145235cc40702d962cbb53ac27fb2246fb6cba"),
					Resource:          "bafkreicjpmrnj2dkhsvj6w5kerbvvgnmcfkassqlsmbltpgz2zke23x35e",
					EncryptedContents: []byte("test test test"),
				},
			},
			&ZeroX02{
				UIBEncryptedLocationHash: []uint8{0x2a, 0xe1, 0x2, 0xa3, 0xaa, 0xe8, 0xec, 0xaf, 0xd5, 0xe3, 0xb6, 0xbc, 0x8c, 0xef, 0x3b, 0x6a, 0x63, 0x79, 0x27, 0x5b, 0x89, 0xed, 0x78, 0x3c, 0xba, 0x9, 0xe4, 0xa0, 0xbc, 0x43, 0xba, 0x45, 0xd6, 0x3c, 0xc1, 0x39, 0x13, 0xd3, 0x7a, 0x16, 0x15, 0x37, 0x7d, 0x92, 0x3d, 0x47, 0x3a, 0x63, 0xef, 0x7c, 0x7a, 0xea, 0x4a, 0x1, 0xe1, 0x31, 0x41, 0xbc, 0xa7, 0x6b, 0x6b, 0x5b, 0xbc, 0x66, 0xda, 0x10, 0xdc, 0xb5, 0x26, 0x19, 0x4a, 0xe7, 0x91, 0x24, 0x34, 0x76, 0x1a, 0x43, 0xa, 0x6, 0xf9, 0x4e, 0xae, 0x9e, 0x1e, 0x77, 0xcd, 0x77, 0xdf, 0x91, 0x73, 0xa2, 0xf9, 0x42, 0xb7, 0x2d, 0x67, 0xe1, 0xfd, 0xef, 0xd4, 0x36, 0x54, 0xb5, 0xb, 0xfe, 0x3a, 0xd1, 0xbd, 0x80, 0x3, 0xf3, 0xdf},
				DecryptedHash:            []uint8{0x2c, 0x84, 0x32, 0xca, 0x28, 0xce, 0x92, 0x9b, 0x86, 0xa4, 0x7f, 0x2d, 0x40, 0x41, 0x3d, 0x16, 0x1f, 0x59, 0x1f, 0x89, 0x85, 0x22, 0x90, 0x60, 0x49, 0x15, 0x73, 0xd8, 0x3f, 0x82, 0xf2, 0x92, 0xf4, 0xdc, 0x68, 0xf9, 0x18, 0x44, 0x63, 0x32, 0x83, 0x7a, 0xa5, 0x7c, 0xd5, 0x14, 0x52, 0x35, 0xcc, 0x40, 0x70, 0x2d, 0x96, 0x2c, 0xbb, 0x53, 0xac, 0x27, 0xfb, 0x22, 0x46, 0xfb, 0x6c, 0xba},
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
					Location:          mli.Mailchain,
					DecryptedHash:     encodingtest.MustDecodeHex("2c8432ca28ce929b86a47f2d40413d161f591f8985229060491573d83f82f292f4dc68f918446332837aa57cd5145235cc40702d962cbb53ac27fb2246fb6cba"),
					Resource:          "bafkreicjpmrnj2dkhsvj6w5kerbvvgnmcfkassqlsmbltpgz2zke23x35e",
					EncryptedContents: []byte("test test test"),
				},
			},
			nil,
			true,
		},
		{
			"err-contents-dont-match-hash",
			args{
				func() cipher.Encrypter {
					m := ciphertest.NewMockEncrypter(mockCtrl)
					return m
				}(),
				secp256k1test.BobPublicKey,
				&CreateOpts{
					Location:          mli.Mailchain,
					DecryptedHash:     encodingtest.MustDecodeHex("2c8432ca28ce929b86a47f2d40413d161f591f8985229060491573d83f82f292f4dc68f918446332837aa57cd5145235cc40702d962cbb53ac27fb2246fb6cba"),
					Resource:          "bafkreicjpmrnj2dkhsvj6w5kerbvvgnmcfkassqlsmbltpgz2zke23x35e",
					EncryptedContents: []byte("not the hash"),
				},
			},
			nil,
			true,
		},
		{
			"err-decode-cod",
			args{
				func() cipher.Encrypter {
					m := ciphertest.NewMockEncrypter(mockCtrl)
					return m
				}(),
				secp256k1test.BobPublicKey,
				&CreateOpts{
					Location:          mli.Mailchain,
					DecryptedHash:     encodingtest.MustDecodeHex("2c8432ca28ce929b86a47f2d40413d161f591f8985229060491573d83f82f292f4dc68f918446332837aa57cd5145235cc40702d962cbb53ac27fb2246fb6cba"),
					Resource:          "bafkreicjpmrnj2dkhsvj6w5kerbvvgnmcflsmbltpgz2zke23x35e",
					EncryptedContents: []byte("test test test"),
				},
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
					Location:          mli.Mailchain,
					DecryptedHash:     encodingtest.MustDecodeHex("2c8432ca28ce929b86a47f2d40413d161f591f8985229060491573d83f82f292f4dc68f918446332837aa57cd5145235cc40702d962cbb53ac27fb2246fb6cba"),
					Resource:          "",
					EncryptedContents: []byte("test test test"),
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
			got, err := NewZeroX02(tt.args.encrypter, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewZeroX02() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("NewZeroX02() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestZeroX02_DecrypterKind(t *testing.T) {
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
			x := &ZeroX02{
				UIBEncryptedLocationHash: tt.fields.UIBEncryptedLocationHash,
			}
			got, err := x.DecrypterKind()
			if (err != nil) != tt.wantErr {
				t.Errorf("ZeroX02.DecrypterKind() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ZeroX02.DecrypterKind() = %v, want %v", got, tt.want)
			}
		})
	}
}
