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

package nacl

import (
	"crypto/rand"
	"io"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/mailchain/mailchain/crypto/ed25519/ed25519test"
	"github.com/mailchain/mailchain/crypto/secp256k1/secp256k1test"
	"github.com/mailchain/mailchain/encoding/encodingtest"
	"github.com/mailchain/mailchain/internal/keystore"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestNewFileStore(t *testing.T) {
	type args struct {
		path   string
		logger io.Writer
	}
	tests := []struct {
		name string
		args args
		want FileStore
	}{
		{
			"success",
			args{"/test", ioutil.Discard},
			FileStore{fs: afero.NewBasePathFs(afero.NewOsFs(), "/test"), rand: rand.Reader, logger: ioutil.Discard},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFileStore(tt.args.path, tt.args.logger); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFileStore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFileStore_filename(t *testing.T) {
	type fields struct {
		fs   afero.Fs
		rand io.Reader
	}
	type args struct {
		pubKeyBytes []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			"secp256k1-sofia",
			fields{},
			args{
				secp256k1test.SofiaPublicKey.Bytes(),
			},
			"69d908510e355beb1d5bf2df8129e5b6401e1969891e8016a0b2300739bbb00687055e5924a2fd8dd35f069dc14d8147aa11c1f7e2f271573487e1beeb2be9d0.json",
		},
		{
			"ed25519-sofia",
			fields{},
			args{
				ed25519test.SofiaPublicKey.Bytes(),
			},
			"723caa23a5b511af5ad7b7ef6076e414ab7e75a9dc910ea60e417a2b770a5671.json",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := FileStore{
				fs:   tt.fields.fs,
				rand: tt.fields.rand,
			}
			if got := f.filename(tt.args.pubKeyBytes); got != tt.want {
				t.Errorf("FileStore.filename() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFileStore_getEncryptedKey(t *testing.T) {
	type fields struct {
		fs   afero.Fs
		rand io.Reader
	}
	type args struct {
		pubKeyBytes []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *keystore.EncryptedKey
		wantErr bool
	}{
		{
			"success-secp256k1",
			fields{
				func() afero.Fs {
					m := afero.NewMemMapFs()
					afero.WriteFile(m, "./0269d908510e355beb1d5bf2df8129e5b6401e1969891e8016a0b2300739bbb006.json", fileSofiaSECP256k1, 0755)
					afero.WriteFile(m, "./723caa23a5b511af5ad7b7ef6076e414ab7e75a9dc910ea60e417a2b770a5671.json", fileSofiaED25519, 0755)

					return m
				}(),
				nil,
			},
			args{
				encodingtest.MustDecodeHex("0269d908510e355beb1d5bf2df8129e5b6401e1969891e8016a0b2300739bbb006"),
			},
			&encryptedKeySofiaSECP256k1,
			false,
		},
		{
			"success-ed25519",
			fields{
				func() afero.Fs {
					m := afero.NewMemMapFs()
					afero.WriteFile(m, "./0269d908510e355beb1d5bf2df8129e5b6401e1969891e8016a0b2300739bbb006.json", fileSofiaSECP256k1, 0755)
					afero.WriteFile(m, "./723caa23a5b511af5ad7b7ef6076e414ab7e75a9dc910ea60e417a2b770a5671.json", fileSofiaED25519, 0755)

					return m
				}(),
				nil,
			},
			args{
				encodingtest.MustDecodeHex("723caa23a5b511af5ad7b7ef6076e414ab7e75a9dc910ea60e417a2b770a5671"),
			},
			&encryptedKeySofiaED25519,
			false,
		},
		{
			"err-public-mismatch",
			fields{
				func() afero.Fs {
					m := afero.NewMemMapFs()
					afero.WriteFile(m, "./0269d908.json", fileSofiaSECP256k1, 0755)
					afero.WriteFile(m, "./723caa23a5b511af5ad7b7ef6076e414ab7e75a9dc910ea60e417a2b770a5671.json", fileSofiaED25519, 0755)

					return m
				}(),
				nil,
			},
			args{
				encodingtest.MustDecodeHex("0269d908"),
			},
			nil,
			true,
		},
		{
			"err-bad-json",
			fields{
				func() afero.Fs {
					m := afero.NewMemMapFs()
					afero.WriteFile(m, "./0269d908510e355beb1d5bf2df8129e5b6401e1969891e8016a0b2300739bbb006.json", []byte("not-json"), 0755)
					afero.WriteFile(m, "./723caa23a5b511af5ad7b7ef6076e414ab7e75a9dc910ea60e417a2b770a5671.json", fileSofiaED25519, 0755)

					return m
				}(),
				nil,
			},
			args{
				encodingtest.MustDecodeHex("0269d908510e355beb1d5bf2df8129e5b6401e1969891e8016a0b2300739bbb006"),
			},
			nil,
			true,
		},
		{
			"err-not-found",
			fields{
				func() afero.Fs {
					m := afero.NewMemMapFs()
					afero.WriteFile(m, "./0269d908510e355beb1d5bf2df8129e5b6401e1969891e8016a0b2300739bbb006.json", fileSofiaSECP256k1, 0755)
					afero.WriteFile(m, "./723caa23a5b511af5ad7b7ef6076e414ab7e75a9dc910ea60e417a2b770a5671.json", fileSofiaED25519, 0755)

					return m
				}(),
				nil,
			},
			args{
				encodingtest.MustDecodeHex("0269d9df"),
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := FileStore{
				fs:   tt.fields.fs,
				rand: tt.fields.rand,
			}
			got, err := f.getEncryptedKey(tt.args.pubKeyBytes)
			if (err != nil) != tt.wantErr {
				t.Errorf("FileStore.getEncryptedKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("FileStore.getEncryptedKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFileStore_getEncryptedKeys(t *testing.T) {
	type fields struct {
		fs     afero.Fs
		rand   io.Reader
		logger io.Writer
	}
	tests := []struct {
		name    string
		fields  fields
		want    []keystore.EncryptedKey
		wantErr bool
	}{
		{
			"success",
			fields{
				func() afero.Fs {
					m := afero.NewMemMapFs()
					afero.WriteFile(m, "./0269d908510e355beb1d5bf2df8129e5b6401e1969891e8016a0b2300739bbb006.json", fileSofiaSECP256k1, 0755)
					afero.WriteFile(m, "./723caa23a5b511af5ad7b7ef6076e414ab7e75a9dc910ea60e417a2b770a5671.json", fileSofiaED25519, 0755)

					return m
				}(),
				nil,
				ioutil.Discard,
			},
			[]keystore.EncryptedKey{
				encryptedKeySofiaSECP256k1,
				encryptedKeySofiaED25519,
			},
			false,
		},
		{
			"success-invalid-key",
			fields{
				func() afero.Fs {
					m := afero.NewMemMapFs()
					afero.WriteFile(m, "./0269d908510e355beb1d5bf2df8129e5b6401e1969891e8016a0b2300739bbb006.json", fileSofiaSECP256k1, 0755)
					afero.WriteFile(m, "./0269d908.json", fileSofiaSECP256k1, 0755)
					afero.WriteFile(m, "./723caa23a5b511af5ad7b7ef6076e414ab7e75a9dc910ea60e417a2b770a5671.json", fileSofiaED25519, 0755)

					return m
				}(),
				nil,
				ioutil.Discard,
			},
			[]keystore.EncryptedKey{
				encryptedKeySofiaSECP256k1,
				encryptedKeySofiaED25519,
			},
			false,
		},
		{
			"success-invalid-file-name",
			fields{
				func() afero.Fs {
					m := afero.NewMemMapFs()
					afero.WriteFile(m, "./0269d908510e355beb1d5bf2df8129e5b6401e1969891e8016a0b2300739bbb006.json", fileSofiaSECP256k1, 0755)
					afero.WriteFile(m, "./invalid-file-name.json", fileSofiaSECP256k1, 0755)
					afero.WriteFile(m, "./723caa23a5b511af5ad7b7ef6076e414ab7e75a9dc910ea60e417a2b770a5671.json", fileSofiaED25519, 0755)

					return m
				}(),
				nil,
				ioutil.Discard,
			},
			[]keystore.EncryptedKey{
				encryptedKeySofiaSECP256k1,
				encryptedKeySofiaED25519,
			},
			false,
		},
		{
			"success-invalid-extension",
			fields{
				func() afero.Fs {
					m := afero.NewMemMapFs()
					afero.WriteFile(m, "./0269d908510e355beb1d5bf2df8129e5b6401e1969891e8016a0b2300739bbb006.json", fileSofiaSECP256k1, 0755)
					afero.WriteFile(m, "./0269d908510e355beb1d5bf2df8129e5b6401e1969891e8016a0b2300739bbb006.invalid", fileSofiaSECP256k1, 0755)
					afero.WriteFile(m, "./723caa23a5b511af5ad7b7ef6076e414ab7e75a9dc910ea60e417a2b770a5671.json", fileSofiaED25519, 0755)

					return m
				}(),
				nil,
				ioutil.Discard,
			},
			[]keystore.EncryptedKey{
				encryptedKeySofiaSECP256k1,
				encryptedKeySofiaED25519,
			},
			false,
		},
		{
			"success-empty",
			fields{
				afero.NewMemMapFs(),
				nil,
				ioutil.Discard,
			},
			[]keystore.EncryptedKey{},
			false,
		},
		{
			"err-read-dir",
			fields{
				func() afero.Fs {
					m := afero.NewMemMapFs()
					afero.WriteFile(m, "./notdir.txt", []byte{}, 0755)
					return afero.NewBasePathFs(m, "./testdata/GetPublicKeys/notdir.txt")
				}(),
				nil,
				ioutil.Discard,
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := FileStore{
				fs:     tt.fields.fs,
				rand:   tt.fields.rand,
				logger: tt.fields.logger,
			}
			got, err := f.getEncryptedKeys()
			if (err != nil) != tt.wantErr {
				t.Errorf("FileStore.getEncryptedKeys() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("FileStore.getEncryptedKeys() = %v, want %v", got, tt.want)
			}
		})
	}
}
