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
	"github.com/mailchain/mailchain/crypto/cipher"
	"github.com/mailchain/mailchain/crypto/cipher/ciphertest"
	"github.com/mailchain/mailchain/encoding/encodingtest"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestZeroX50_URL(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type fields struct {
		Version       int32
		EncryptedURL  []byte
		DecryptedHash []byte
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
				0,
				[]byte("EncryptedURL"),
				[]byte("DecryptedHash"),
			},
			args{
				func() cipher.Decrypter {
					m := ciphertest.NewMockDecrypter(mockCtrl)

					m.EXPECT().Decrypt(cipher.EncryptedContent([]byte("EncryptedURL"))).Return([]byte("https://domain.com/resource"), nil)
					return m
				}(),
			},
			func() *url.URL {
				u, _ := url.Parse("https://domain.com/resource")
				return u
			}(),
			false,
		},
		{
			"err-invalid-url",
			fields{
				0,
				[]byte("EncryptedURL"),
				[]byte("DecryptedHash"),
			},
			args{
				func() cipher.Decrypter {
					m := ciphertest.NewMockDecrypter(mockCtrl)

					m.EXPECT().Decrypt(cipher.EncryptedContent([]byte("EncryptedURL"))).Return([]byte("://resource"), nil)
					return m
				}(),
			},
			nil,
			true,
		},
		{
			"err-decrypt",
			fields{
				0,
				[]byte("EncryptedURL"),
				[]byte("DecryptedHash"),
			},
			args{
				func() cipher.Decrypter {
					m := ciphertest.NewMockDecrypter(mockCtrl)

					m.EXPECT().Decrypt(cipher.EncryptedContent([]byte("EncryptedURL"))).Return(nil, errors.Errorf("failed"))
					return m
				}(),
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &ZeroX50{
				Version:       tt.fields.Version,
				EncryptedURL:  tt.fields.EncryptedURL,
				DecryptedHash: tt.fields.DecryptedHash,
			}
			got, err := d.URL(tt.args.decrypter)
			if (err != nil) != tt.wantErr {
				t.Errorf("ZeroX50.URL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ZeroX50.URL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestZeroX50_ContentsHash(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type fields struct {
		Version       int32
		EncryptedURL  []byte
		DecryptedHash []byte
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
				0,
				[]byte("EncryptedURL"),
				[]byte("DecryptedHash"),
			},
			args{
				func() cipher.Decrypter {
					m := ciphertest.NewMockDecrypter(mockCtrl)

					return m
				}(),
			},
			[]byte("DecryptedHash"),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &ZeroX50{
				Version:       tt.fields.Version,
				EncryptedURL:  tt.fields.EncryptedURL,
				DecryptedHash: tt.fields.DecryptedHash,
			}
			got, err := d.ContentsHash(tt.args.decrypter)
			if (err != nil) != tt.wantErr {
				t.Errorf("ZeroX50.ContentsHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ZeroX50.ContentsHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestZeroX50_IntegrityHash(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type fields struct {
		Version       int32
		EncryptedURL  []byte
		DecryptedHash []byte
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
				0,
				[]byte("EncryptedURL"),
				[]byte("DecryptedHash"),
			},
			args{
				func() cipher.Decrypter {
					m := ciphertest.NewMockDecrypter(mockCtrl)

					m.EXPECT().Decrypt(cipher.EncryptedContent([]byte("EncryptedURL"))).Return(cipher.PlainContent("https://domain.com/resource-220455078214"), nil)
					return m
				}(),
			},
			encodingtest.MustDecodeHex("220455078214"),
			false,
		},
		{
			"err-invalid-hash",
			fields{
				0,
				[]byte("EncryptedURL"),
				[]byte("DecryptedHash"),
			},
			args{
				func() cipher.Decrypter {
					m := ciphertest.NewMockDecrypter(mockCtrl)

					m.EXPECT().Decrypt(cipher.EncryptedContent([]byte("EncryptedURL"))).Return(cipher.PlainContent("https://domain.com/resource-h22055078214"), nil)
					return m
				}(),
			},
			[]byte{},
			true,
		},
		{
			"err-parts",
			fields{
				0,
				[]byte("EncryptedURL"),
				[]byte("DecryptedHash"),
			},
			args{
				func() cipher.Decrypter {
					m := ciphertest.NewMockDecrypter(mockCtrl)

					m.EXPECT().Decrypt(cipher.EncryptedContent([]byte("EncryptedURL"))).Return(cipher.PlainContent("https://domain.com/resource"), nil)
					return m
				}(),
			},
			nil,
			true,
		},
		{
			"err-decrypt",
			fields{
				0,
				[]byte("EncryptedURL"),
				[]byte("DecryptedHash"),
			},
			args{
				func() cipher.Decrypter {
					m := ciphertest.NewMockDecrypter(mockCtrl)

					m.EXPECT().Decrypt(cipher.EncryptedContent([]byte("EncryptedURL"))).Return(nil, errors.Errorf("failed"))
					return m
				}(),
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &ZeroX50{
				Version:       tt.fields.Version,
				EncryptedURL:  tt.fields.EncryptedURL,
				DecryptedHash: tt.fields.DecryptedHash,
			}
			got, err := d.IntegrityHash(tt.args.decrypter)
			if (err != nil) != tt.wantErr {
				t.Errorf("ZeroX50.IntegrityHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("ZeroX50.IntegrityHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestZeroX50_Valid(t *testing.T) {
	type fields struct {
		Version       int32
		EncryptedURL  []byte
		DecryptedHash []byte
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			"success",
			fields{
				0,
				[]byte("EncryptedURL"),
				[]byte("DecryptedHash"),
			},
			false,
		},
		{
			"err-decryptedHash",
			fields{
				0,
				[]byte("EncryptedURL"),
				[]byte{},
			},
			true,
		},
		{
			"err-EncryptedURL",
			fields{
				0,
				[]byte{},
				[]byte{},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &ZeroX50{
				Version:       tt.fields.Version,
				EncryptedURL:  tt.fields.EncryptedURL,
				DecryptedHash: tt.fields.DecryptedHash,
			}
			if err := d.Valid(); (err != nil) != tt.wantErr {
				t.Errorf("ZeroX50.Valid() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestZeroX50_DecrypterKind(t *testing.T) {
	type fields struct {
		EncryptedURL []byte
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
			x := &ZeroX50{
				EncryptedURL: tt.fields.EncryptedURL,
			}
			got, err := x.DecrypterKind()
			if (err != nil) != tt.wantErr {
				t.Errorf("ZeroX50.DecrypterKind() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ZeroX50.DecrypterKind() = %v, want %v", got, tt.want)
			}
		})
	}
}
