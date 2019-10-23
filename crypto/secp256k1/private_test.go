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
		{"A", fields{ecdsaPrivateKeyA()}, []byte{0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa}},
		{"B", fields{ecdsaPrivateKeyB()}, []byte{0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pk := PrivateKeyFromECDSA(tt.fields.ecdsa)
			if got := pk.Bytes(); !assert.Equal(got, tt.want) {
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
		{"A", fields{ecdsaPrivateKeyA()}, []byte{0x2, 0x6a, 0x4, 0xab, 0x98, 0xd9, 0xe4, 0x77, 0x4a, 0xd8, 0x6, 0xe3, 0x2, 0xdd, 0xde, 0xb6, 0x3b, 0xea, 0x16, 0xb5, 0xcb, 0x5f, 0x22, 0x3e, 0xe7, 0x74, 0x78, 0xe8, 0x61, 0xbb, 0x58, 0x3e, 0xb3}},
		{"B", fields{ecdsaPrivateKeyB()}, []byte{0x2, 0x68, 0x68, 0x7, 0x37, 0xc7, 0x6d, 0xab, 0xb8, 0x1, 0xcb, 0x22, 0x4, 0xf5, 0x7d, 0xbe, 0x4e, 0x45, 0x79, 0xe4, 0xf7, 0x10, 0xcd, 0x67, 0xdc, 0x1b, 0x42, 0x27, 0x59, 0x2c, 0x81, 0xe9, 0xb5}},
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
		{"A", args{ecdsaPrivateKeyA()}, []byte{0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa}},
		{"B", args{ecdsaPrivateKeyB()}, []byte{0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PrivateKeyFromECDSA(tt.args.pk); !assert.Equal(got.Bytes(), tt.want) {
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
				ecdsaPrivateKeyA(),
			},
			ecdsaPrivateKeyA(),
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
				ecdsaPrivateKeyA(),
			},
			func() *ecdsa.PrivateKey {
				k := ecdsaPrivateKeyA()
				return &k
			}(),
			false,
		},
		{
			"err",
			fields{
				func() ecdsa.PrivateKey {
					k := ecdsaPrivateKeyA()
					k.Y = nil
					return k
				}(),
			},
			func() *ecdsa.PrivateKey {
				k := ecdsaPrivateKeyA()
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
