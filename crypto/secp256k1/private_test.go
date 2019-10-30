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

package secp256k1

import (
	"crypto/ecdsa"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrivateKey_Bytes(t *testing.T) {
	assert := assert.New(t)
	type fields struct {
		ecdsa ecdsa.PrivateKey
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{"A", fields{ecdsaPrivateKeySofia()}, []byte{0x1, 0x90, 0x1e, 0x63, 0x38, 0x9e, 0xf0, 0x2e, 0xaa, 0x7c, 0x57, 0x82, 0xe0, 0x8b, 0x40, 0xd9, 0x8f, 0xae, 0xf8, 0x35, 0xf2, 0x8b, 0xd1, 0x44, 0xee, 0xcf, 0x56, 0x14, 0xa4, 0x15, 0x94, 0x3f}},
		{"B", fields{ecdsaPrivateKeyCharlotte()}, []byte{0xdf, 0x4b, 0xa9, 0xf6, 0x10, 0x6a, 0xd2, 0x84, 0x64, 0x72, 0xf7, 0x59, 0x47, 0x65, 0x35, 0xe5, 0x5c, 0x58, 0x5, 0xd8, 0x33, 0x7d, 0xf5, 0xa1, 0x1c, 0x3b, 0x13, 0x9f, 0x43, 0x8b, 0x98, 0xb3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pk := PrivateKeyFromECDSA(tt.fields.ecdsa)
			if got := pk.Bytes(); !assert.Equal(tt.want, got) {
				t.Errorf("PrivateKey.Bytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrivateKey_PublicKey(t *testing.T) {
	assert := assert.New(t)
	type fields struct {
		ecdsa ecdsa.PrivateKey
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{"A", fields{ecdsaPrivateKeySofia()}, []byte{0x2, 0x69, 0xd9, 0x8, 0x51, 0xe, 0x35, 0x5b, 0xeb, 0x1d, 0x5b, 0xf2, 0xdf, 0x81, 0x29, 0xe5, 0xb6, 0x40, 0x1e, 0x19, 0x69, 0x89, 0x1e, 0x80, 0x16, 0xa0, 0xb2, 0x30, 0x7, 0x39, 0xbb, 0xb0, 0x6}},
		{"B", fields{ecdsaPrivateKeyCharlotte()}, []byte{0x3, 0xbd, 0xf6, 0xfb, 0x97, 0xc9, 0x7c, 0x12, 0x6b, 0x49, 0x21, 0x86, 0xa4, 0xd5, 0xb2, 0x8f, 0x34, 0xf0, 0x67, 0x1a, 0x5a, 0xac, 0xc9, 0x74, 0xda, 0x3b, 0xde, 0xb, 0xe9, 0x3e, 0x45, 0xa1, 0xc5}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pk := PrivateKeyFromECDSA(tt.fields.ecdsa)
			if got := pk.PublicKey(); !assert.Equal(tt.want, got.Bytes()) {
				t.Errorf("PrivateKey.PublicKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrivateKeyFromECDSA(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		pk ecdsa.PrivateKey
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{"A", args{ecdsaPrivateKeySofia()}, []byte{0x1, 0x90, 0x1e, 0x63, 0x38, 0x9e, 0xf0, 0x2e, 0xaa, 0x7c, 0x57, 0x82, 0xe0, 0x8b, 0x40, 0xd9, 0x8f, 0xae, 0xf8, 0x35, 0xf2, 0x8b, 0xd1, 0x44, 0xee, 0xcf, 0x56, 0x14, 0xa4, 0x15, 0x94, 0x3f}},
		{"B", args{ecdsaPrivateKeyCharlotte()}, []byte{0xdf, 0x4b, 0xa9, 0xf6, 0x10, 0x6a, 0xd2, 0x84, 0x64, 0x72, 0xf7, 0x59, 0x47, 0x65, 0x35, 0xe5, 0x5c, 0x58, 0x5, 0xd8, 0x33, 0x7d, 0xf5, 0xa1, 0x1c, 0x3b, 0x13, 0x9f, 0x43, 0x8b, 0x98, 0xb3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PrivateKeyFromECDSA(tt.args.pk); !assert.Equal(tt.want, got.Bytes()) {
				t.Errorf("PrivateKeyFromECDSA() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrivateKeyFromHex(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		hexkey string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			"success",
			args{
				"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
			},
			[]byte{0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa},
			false,
		},
		{
			"err-decode",
			args{
				"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
			},
			[]byte{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PrivateKeyFromHex(tt.args.hexkey)
			if (err != nil) != tt.wantErr {
				t.Errorf("PrivateKeyFromHex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			gotBytes := []byte{}
			if got != nil {
				gotBytes = got.Bytes()
			}
			if !assert.Equal(tt.want, gotBytes) {
				t.Errorf("PrivateKeyFromHex() = %v, want %v", gotBytes, tt.want)
			}
		})
	}
}

func TestPrivateKeyFromBytes(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		pk []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			"success",
			args{
				[]byte{0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa},
			},
			[]byte{0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa},
			false,
		},
		{
			"err-to-ecdsa",
			args{
				[]byte{0xaa},
			},
			[]byte{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PrivateKeyFromBytes(tt.args.pk)
			if (err != nil) != tt.wantErr {
				t.Errorf("PrivateKeyFromBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			gotBytes := []byte{}
			if got != nil {
				gotBytes = got.Bytes()
			}
			if !assert.Equal(tt.want, gotBytes) {
				t.Errorf("PrivateKeyFromBytes() = %v, want %v", gotBytes, tt.want)
			}
		})
	}
}

func TestPrivateKey_ECIES(t *testing.T) {
	assert := assert.New(t)
	type fields struct {
		ecdsa ecdsa.PrivateKey
	}
	tests := []struct {
		name   string
		fields fields
		want   ecdsa.PrivateKey
	}{
		{
			"success",
			fields{
				ecdsaPrivateKeySofia(),
			},
			ecdsaPrivateKeySofia(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pk := PrivateKey{
				ecdsa: tt.fields.ecdsa,
			}
			got := pk.ECIES().ExportECDSA()
			if !assert.Equal(got, &tt.want) {
				t.Errorf("PrivateKey.ECIES() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrivateKey_ECDSA(t *testing.T) {
	type fields struct {
		ecdsa ecdsa.PrivateKey
	}
	tests := []struct {
		name    string
		fields  fields
		want    *ecdsa.PrivateKey
		wantErr bool
	}{
		{
			"success",
			fields{
				ecdsaPrivateKeySofia(),
			},
			func() *ecdsa.PrivateKey {
				k := ecdsaPrivateKeySofia()
				return &k
			}(),
			false,
		},
		{
			"err",
			fields{
				func() ecdsa.PrivateKey {
					k := ecdsaPrivateKeySofia()
					k.Y = nil
					return k
				}(),
			},
			func() *ecdsa.PrivateKey {
				k := ecdsaPrivateKeySofia()
				return &k
			}(),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pk := PrivateKey{
				ecdsa: tt.fields.ecdsa,
			}
			got, err := pk.ECDSA()
			if (err != nil) != tt.wantErr {
				t.Errorf("PrivateKey.ECDSA() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PrivateKey.ECDSA() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrivateKey_Kind(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			"success",
			"secp256k1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pk := PrivateKey{}
			if got := pk.Kind(); got != tt.want {
				t.Errorf("PrivateKey.Kind() = %v, want %v", got, tt.want)
			}
		})
	}
}
