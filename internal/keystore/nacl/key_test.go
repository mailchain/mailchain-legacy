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
	"io"
	"io/ioutil"
	"testing"

	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/ed25519/ed25519test"
	"github.com/mailchain/mailchain/crypto/secp256k1/secp256k1test"
	"github.com/mailchain/mailchain/internal/keystore"
	"github.com/mailchain/mailchain/internal/keystore/kdf/multi"
	"github.com/mailchain/mailchain/internal/keystore/kdf/scrypt"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func Test_deriveKey(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		ek               *keystore.EncryptedKey
		deriveKeyOptions multi.OptionsBuilders
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			"success-kdf-secp256k1",
			args{
				&encryptedKeySofiaSECP256k1,
				multi.OptionsBuilders{
					Scrypt: []scrypt.DeriveOptionsBuilder{scrypt.WithPassphrase("test")},
				},
			},
			[]byte{0xce, 0xc0, 0x96, 0x7d, 0xb4, 0x55, 0xbb, 0x94, 0x0, 0xd0, 0xb7, 0x88, 0xaf, 0xdb, 0xad, 0x43, 0xde, 0xe6, 0x83, 0x1a, 0xca, 0x9b, 0xfc, 0x1a, 0x87, 0xe6, 0x4c, 0x31, 0x9a, 0xcd, 0x30, 0xde},
			false,
		},
		{
			"success-kdf-ed25519",
			args{
				&encryptedKeySofiaED25519,
				multi.OptionsBuilders{
					Scrypt: []scrypt.DeriveOptionsBuilder{scrypt.WithPassphrase("sofia-ed25519")},
				},
			},
			[]byte{0x7c, 0xb3, 0xf1, 0xbf, 0xe8, 0x19, 0xdb, 0x82, 0x55, 0xb2, 0x19, 0xd4, 0x1e, 0xa8, 0x7f, 0xb7, 0x13, 0x67, 0x20, 0x45, 0x7f, 0x6a, 0xcf, 0x4c, 0xb8, 0xde, 0x52, 0x91, 0xf3, 0x2e, 0xd0, 0xb6},
			false,
		},
		{
			"err-nil-script-params",
			args{
				func() *keystore.EncryptedKey {
					m := encryptedKeySofiaED25519
					m.ScryptParams = nil
					return &m
				}(),
				multi.OptionsBuilders{
					Scrypt: []scrypt.DeriveOptionsBuilder{scrypt.WithPassphrase("test")},
				},
			},
			nil,
			true,
		},
		{
			"err-kdf-invalid",
			args{
				func() *keystore.EncryptedKey {
					m := encryptedKeySofiaED25519
					m.KDF = "invalid"
					return &m
				}(),
				multi.OptionsBuilders{
					Scrypt: []scrypt.DeriveOptionsBuilder{scrypt.WithPassphrase("test")},
				},
			},
			nil,
			true,
		},
		{
			"err-nil-encrypted-key",
			args{
				nil,
				multi.OptionsBuilders{
					Scrypt: []scrypt.DeriveOptionsBuilder{scrypt.WithPassphrase("test")},
				},
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := deriveKey(tt.args.ek, tt.args.deriveKeyOptions)
			if (err != nil) != tt.wantErr {
				t.Errorf("deriveKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(tt.want, got) {
				t.Errorf("deriveKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFileStore_getPrivateKey(t *testing.T) {
	assert := assert.New(t)
	type fields struct {
		fs     afero.Fs
		rand   io.Reader
		logger io.Writer
	}
	type args struct {
		encryptedKey     *keystore.EncryptedKey
		deriveKeyOptions multi.OptionsBuilders
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    crypto.PrivateKey
		wantErr bool
	}{
		{
			"success-sofia-secp256k1",
			fields{
				nil,
				nil,
				ioutil.Discard,
			},
			args{
				&encryptedKeySofiaSECP256k1,
				multi.OptionsBuilders{
					Scrypt: []scrypt.DeriveOptionsBuilder{scrypt.WithPassphrase("test")},
				},
			},
			secp256k1test.SofiaPrivateKey,
			false,
		},
		{
			"success-sofia-ed25519",
			fields{
				nil,
				nil,
				ioutil.Discard,
			},
			args{
				&encryptedKeySofiaED25519,
				multi.OptionsBuilders{
					Scrypt: []scrypt.DeriveOptionsBuilder{scrypt.WithPassphrase("sofia-ed25519")},
				},
			},
			ed25519test.SofiaPrivateKey,
			false,
		},
		//{
		//	"success-charlotte-sr25519",
		//	fields{
		//		nil,
		//		nil,
		//		ioutil.Discard,
		//	},
		//	args{
		//		&encryptedKeyCharlotteSR25519,
		//		multi.OptionsBuilders{
		//			Scrypt: []scrypt.DeriveOptionsBuilder{scrypt.WithPassphrase("sofia-ed25519")},
		//		},
		//	},
		//	sr25519test.SofiaPrivateKey,
		//	false,
		//},
		{
			"err-private-key-bytes",
			fields{
				nil,
				nil,
				ioutil.Discard,
			},
			args{
				func() *keystore.EncryptedKey {
					m := encryptedKeySofiaED25519
					m.CurveType = "invalid"
					return &m
				}(),
				multi.OptionsBuilders{
					Scrypt: []scrypt.DeriveOptionsBuilder{scrypt.WithPassphrase("sofia-ed25519")},
				},
			},
			nil,
			true,
		},
		{
			"err-wrong-key",
			fields{
				nil,
				nil,
				ioutil.Discard,
			},
			args{
				&encryptedKeyCharlotteED25519,
				multi.OptionsBuilders{
					Scrypt: []scrypt.DeriveOptionsBuilder{scrypt.WithPassphrase("sofia-ed25519")},
				},
			},
			nil,
			true,
		},
		{
			"err-derive-key",
			fields{
				nil,
				nil,
				ioutil.Discard,
			},
			args{
				func() *keystore.EncryptedKey {
					m := encryptedKeySofiaED25519
					m.KDF = "invalid"
					return &m
				}(),
				multi.OptionsBuilders{
					Scrypt: []scrypt.DeriveOptionsBuilder{scrypt.WithPassphrase("sofia-ed25519")},
				},
			},
			nil,
			true,
		},
		{
			"err-nil-encrypted-key",
			fields{
				nil,
				nil,
				ioutil.Discard,
			},
			args{
				nil,
				multi.OptionsBuilders{
					Scrypt: []scrypt.DeriveOptionsBuilder{scrypt.WithPassphrase("sofia-ed25519")},
				},
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
			got, err := f.getPrivateKey(tt.args.encryptedKey, tt.args.deriveKeyOptions)
			if (err != nil) != tt.wantErr {
				t.Errorf("FileStore.getPrivateKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(tt.want, got) {
				t.Errorf("FileStore.getPrivateKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
