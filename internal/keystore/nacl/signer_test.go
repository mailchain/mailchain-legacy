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

	"github.com/mailchain/mailchain/internal/keystore/kdf/multi"
	"github.com/mailchain/mailchain/internal/keystore/kdf/scrypt"
	"github.com/mailchain/mailchain/internal/mailbox/signer"
	"github.com/mailchain/mailchain/internal/protocols/ethereum"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestFileStore_GetSigner(t *testing.T) {
	assert := assert.New(t)
	type fields struct {
		fs     afero.Fs
		rand   io.Reader
		logger io.Writer
	}
	type args struct {
		address          []byte
		protocol         string
		network          string
		deriveKeyOptions multi.OptionsBuilders
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    signer.Signer
		wantErr bool
	}{
		{
			"success-sofia-secp256k1",
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
			args{
				[]byte{0xd5, 0xab, 0x4c, 0xe3, 0x60, 0x5c, 0xd5, 0x90, 0xdb, 0x60, 0x9b, 0x6b, 0x5c, 0x89, 0x1, 0xfd, 0xb2, 0xef, 0x7f, 0xe6},
				"ethereum",
				"mainnet",
				multi.OptionsBuilders{
					Scrypt: []scrypt.DeriveOptionsBuilder{scrypt.WithPassphrase("test")},
				},
			},
			&ethereum.Signer{},
			false,
		},
		{
			"err-read-private-key",
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
			args{
				[]byte{0xd5, 0xab, 0x4c, 0xe3, 0x60, 0x5c, 0xd5, 0x90, 0xdb, 0x60, 0x9b, 0x6b, 0x5c, 0x89, 0x1, 0xfd, 0xb2, 0xef, 0x7f, 0xe6},
				"ethereum",
				"mainnet",
				multi.OptionsBuilders{
					Scrypt: []scrypt.DeriveOptionsBuilder{scrypt.WithPassphrase("invalid")},
				},
			},
			nil,
			true,
		},
		{
			"err-key-not-found",
			fields{
				func() afero.Fs {
					m := afero.NewMemMapFs()
					return m
				}(),
				nil,
				ioutil.Discard,
			},
			args{
				[]byte{0xd5, 0xab, 0x4c, 0xe3, 0x60, 0x5c, 0xd5, 0x90, 0xdb, 0x60, 0x9b, 0x6b, 0x5c, 0x89, 0x1, 0xfd, 0xb2, 0xef, 0x7f, 0xe6},
				"ethereum",
				"mainnet",
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
			f := FileStore{
				fs:     tt.fields.fs,
				rand:   tt.fields.rand,
				logger: tt.fields.logger,
			}
			got, err := f.GetSigner(tt.args.address, tt.args.protocol, tt.args.network, tt.args.deriveKeyOptions)
			if (err != nil) != tt.wantErr {
				t.Errorf("FileStore.GetSigner() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.IsType(got, tt.want) {
				t.Errorf("FileStore.GetSigner() = %v, want %v", got, tt.want)
			}
		})
	}
}
