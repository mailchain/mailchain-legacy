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
	"bytes"
	"flag"
	"io"
	"testing"
	"testing/iotest"

	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/ed25519/ed25519test"
	"github.com/mailchain/mailchain/crypto/secp256k1/secp256k1test"
	"github.com/mailchain/mailchain/crypto/sr25519/sr25519test"
	"github.com/mailchain/mailchain/internal/keystore/kdf/multi"
	"github.com/mailchain/mailchain/internal/keystore/kdf/scrypt"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

var update = flag.Bool("update", false, "update .golden files")

func Test_writeTemporaryKeyFile(t *testing.T) {
	type args struct {
		fs      afero.Fs
		file    string
		content []byte
	}
	tests := []struct {
		name         string
		args         args
		wantFileName string
		wantErr      bool
	}{
		{
			"success",
			args{
				afero.NewMemMapFs(),
				"file.json.tmp",
				[]byte("contents"),
			},
			".file.json.",
			false,
		},
		{
			"err-mkdir-failed",
			args{
				afero.NewReadOnlyFs(afero.NewMemMapFs()),
				"file.json.tmp",
				[]byte("contents"),
			},
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := writeTemporaryKeyFile(tt.args.fs, tt.args.file, tt.args.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("writeTemporaryKeyFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Contains(t, got, tt.wantFileName) {
				t.Errorf("writeTemporaryKeyFile() = %v, want %v", got, tt.wantFileName)
			}
		})
	}
}

func TestFileStore_Store(t *testing.T) {
	type fields struct {
		fs   afero.Fs
		rand io.Reader
	}
	type args struct {
		protocol         string
		network          string
		private          crypto.PrivateKey
		deriveKeyOptions multi.OptionsBuilders
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      crypto.PublicKey
		wantErr   bool
		wantFiles [][]byte
	}{
		{
			"success-alice-secp256k1",
			fields{
				afero.NewMemMapFs(),
				bytes.NewReader([]byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")),
			},
			args{
				"ethereum",
				"mainnet",
				secp256k1test.AlicePrivateKey,
				multi.OptionsBuilders{
					Scrypt: []scrypt.DeriveOptionsBuilder{
						scrypt.DefaultDeriveOptions(),
						scrypt.WithPassphrase("passphrase"),
					},
				},
			},
			secp256k1test.AlicePublicKey,
			false,
			nil,
		},
		{
			"success-alice-ed25519",
			fields{
				afero.NewMemMapFs(),
				bytes.NewReader([]byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")),
			},
			args{
				"algorand",
				"mainnet",
				ed25519test.AlicePrivateKey,
				multi.OptionsBuilders{
					Scrypt: []scrypt.DeriveOptionsBuilder{
						scrypt.DefaultDeriveOptions(),
						scrypt.WithPassphrase("passphrase"),
					},
				},
			},
			ed25519test.AlicePublicKey,
			false,
			nil,
		},
		{
			"success-bob-sr25519",
			fields{
				afero.NewMemMapFs(),
				bytes.NewReader([]byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")),
			},
			args{
				"substrate",
				"mainnet",
				sr25519test.BobPrivateKey,
				multi.OptionsBuilders{
					Scrypt: []scrypt.DeriveOptionsBuilder{
						scrypt.DefaultDeriveOptions(),
						scrypt.WithPassphrase("passphrase"),
					},
				},
			},
			sr25519test.BobPublicKey,
			false,
			nil,
		},
		{
			"err-seal",
			fields{
				afero.NewMemMapFs(),
				iotest.DataErrReader(bytes.NewReader(nil)),
			},
			args{
				"algorand",
				"mainnet",
				ed25519test.AlicePrivateKey,
				multi.OptionsBuilders{
					Scrypt: []scrypt.DeriveOptionsBuilder{
						scrypt.DefaultDeriveOptions(),
						scrypt.WithPassphrase("passphrase"),
					},
				},
			},
			nil,
			true,
			nil,
		},
		{
			"err-write-fail",
			fields{
				afero.NewReadOnlyFs(afero.NewMemMapFs()),
				bytes.NewReader([]byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")),
			},
			args{
				"algorand",
				"mainnet",
				ed25519test.AlicePrivateKey,
				multi.OptionsBuilders{
					Scrypt: []scrypt.DeriveOptionsBuilder{
						scrypt.DefaultDeriveOptions(),
						scrypt.WithPassphrase("passphrase"),
					},
				},
			},
			nil,
			true,
			nil,
		},
		{
			"err-storage-key",
			fields{
				afero.NewReadOnlyFs(afero.NewMemMapFs()),
				bytes.NewReader([]byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")),
			},
			args{
				"algorand",
				"mainnet",
				ed25519test.AlicePrivateKey,
				multi.OptionsBuilders{
					Scrypt: []scrypt.DeriveOptionsBuilder{},
				},
			},
			nil,
			true,
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := FileStore{
				fs:   tt.fields.fs,
				rand: tt.fields.rand,
			}
			got, err := f.Store(tt.args.protocol, tt.args.network, tt.args.private, tt.args.deriveKeyOptions)
			if (err != nil) != tt.wantErr {
				t.Errorf("FileStore.Store() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("FileStore.Store() = %v, want %v", got, tt.want)
			}
		})
	}
}
