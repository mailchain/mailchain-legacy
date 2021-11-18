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
	"testing"

	"github.com/mailchain/mailchain/crypto/ed25519/ed25519test"
	"github.com/mailchain/mailchain/crypto/secp256k1/secp256k1test"
	"github.com/mailchain/mailchain/internal/keystore"
	"github.com/mailchain/mailchain/internal/protocols/substrate"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestFileStore_getPublicKeys(t *testing.T) {
	type fields struct {
		fs   afero.Fs
		rand io.Reader
	}
	type args struct {
		protocol string
		network  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    [][]byte
		wantErr bool
	}{
		{
			"success",
			fields{
				func() afero.Fs {
					m := afero.NewMemMapFs()
					afero.WriteFile(m, "./ethereum/mainnet/0269d908510e355beb1d5bf2df8129e5b6401e1969891e8016a0b2300739bbb006.json", fileAliceSECP256k1, 0755)
					afero.WriteFile(m, "./ethereum/mainnet/723caa23a5b511af5ad7b7ef6076e414ab7e75a9dc910ea60e417a2b770a5671.json", fileAliceED25519, 0755)

					return m
				}(),
				nil,
			},
			args{
				"ethereum",
				"mainnet",
			},
			[][]byte{
				secp256k1test.AlicePublicKey.Bytes(),
				ed25519test.AlicePublicKey.Bytes(),
			},
			false,
		},
		{
			"success-invalid-curve-type",
			fields{
				func() afero.Fs {
					m := afero.NewMemMapFs()
					afero.WriteFile(m, "./ethereum/mainnet/0269d908510e355beb1d5bf2df8129e5b6401e1969891e8016a0b2300739bbb006.json", fileInvalidCurve, 0755)
					afero.WriteFile(m, "./ethereum/mainnet/723caa23a5b511af5ad7b7ef6076e414ab7e75a9dc910ea60e417a2b770a5671.json", fileAliceED25519, 0755)

					return m
				}(),
				nil,
			},
			args{
				"ethereum",
				"mainnet",
			},
			[][]byte{
				ed25519test.AlicePublicKey.Bytes(),
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := FileStore{
				fs:   tt.fields.fs,
				rand: tt.fields.rand,
			}
			got, err := f.getPublicKeys(tt.args.protocol, tt.args.network)
			if (err != nil) != tt.wantErr {
				t.Errorf("FileStore.GetPublicKeys() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			gotBytes := [][]byte{}
			for _, x := range got {
				gotBytes = append(gotBytes, x.Bytes())
			}
			if !assert.Equal(t, tt.want, gotBytes) {
				t.Errorf("FileStore.GetPublicKeys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFileStore_GetAddresses(t *testing.T) {
	type fields struct {
		fs   afero.Fs
		rand io.Reader
	}
	type args struct {
		protocol string
		network  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]map[string][][]byte
		wantErr bool
	}{
		{
			"empty-query-protocol-network",
			fields{
				func() afero.Fs {
					m := afero.NewMemMapFs()
					return m
				}(),
				nil,
			},
			args{
				"ethereum",
				"mainnet",
			},
			map[string]map[string][][]uint8{
				"ethereum": {
					"mainnet": [][]uint8{},
				},
			},
			false,
		},
		{
			"empty-query-protocol",
			fields{
				func() afero.Fs {
					m := afero.NewMemMapFs()
					return m
				}(),
				nil,
			},
			args{
				"ethereum",
				"",
			},
			map[string]map[string][][]uint8{
				"ethereum": {
					"goerli":  [][]uint8{},
					"kovan":   [][]uint8{},
					"mainnet": [][]uint8{},
					"rinkeby": [][]uint8{},
					"ropsten": [][]uint8{},
				},
			},
			false,
		},
		{
			"success-ethereum",
			fields{
				func() afero.Fs {
					m := afero.NewMemMapFs()
					afero.WriteFile(m, "./ethereum/mainnet/0269d908510e355beb1d5bf2df8129e5b6401e1969891e8016a0b2300739bbb006.json", fileAliceSECP256k1, 0755)
					afero.WriteFile(m, "./ethereum/goerli/0269d908510e355beb1d5bf2df8129e5b6401e1969891e8016a0b2300739bbb006.json", fileAliceSECP256k1, 0755)
					afero.WriteFile(m, "./substrate/edgeware-beresheet/723caa23a5b511af5ad7b7ef6076e414ab7e75a9dc910ea60e417a2b770a5671.json", fileAliceED25519, 0755)
					afero.WriteFile(m, "./substrate/edgeware-beresheet/2e322f8740c60172111ac8eadcdda2512f90d06d0e503ef189979a159bece1e8.json", fileBobED25519, 0755)
					afero.WriteFile(m, "./substrate/edgeware-mainnet/723caa23a5b511af5ad7b7ef6076e414ab7e75a9dc910ea60e417a2b770a5671.json", fileAliceED25519, 0755)
					afero.WriteFile(m, "./substrate/edgeware-mainnet/2e322f8740c60172111ac8eadcdda2512f90d06d0e503ef189979a159bece1e8.json", fileBobED25519, 0755)

					return m
				}(),
				nil,
			},
			args{
				"ethereum",
				"mainnet",
			},
			map[string]map[string][][]uint8{
				"ethereum": {
					"mainnet": [][]uint8{{0xd5, 0xab, 0x4c, 0xe3, 0x60, 0x5c, 0xd5, 0x90, 0xdb, 0x60, 0x9b, 0x6b, 0x5c, 0x89, 0x1, 0xfd, 0xb2, 0xef, 0x7f, 0xe6}},
				},
			},
			false,
		},
		{
			"success-edgeware-beresheet",
			fields{
				func() afero.Fs {
					m := afero.NewMemMapFs()
					afero.WriteFile(m, "./ethereum/mainnet/0269d908510e355beb1d5bf2df8129e5b6401e1969891e8016a0b2300739bbb006.json", fileAliceSECP256k1, 0755)
					afero.WriteFile(m, "./ethereum/goerli/0269d908510e355beb1d5bf2df8129e5b6401e1969891e8016a0b2300739bbb006.json", fileAliceSECP256k1, 0755)
					afero.WriteFile(m, "./substrate/edgeware-beresheet/723caa23a5b511af5ad7b7ef6076e414ab7e75a9dc910ea60e417a2b770a5671.json", fileAliceED25519, 0755)
					afero.WriteFile(m, "./substrate/edgeware-beresheet/2e322f8740c60172111ac8eadcdda2512f90d06d0e503ef189979a159bece1e8.json", fileBobED25519, 0755)
					afero.WriteFile(m, "./substrate/edgeware-mainnet/723caa23a5b511af5ad7b7ef6076e414ab7e75a9dc910ea60e417a2b770a5671.json", fileAliceED25519, 0755)
					afero.WriteFile(m, "./substrate/edgeware-mainnet/2e322f8740c60172111ac8eadcdda2512f90d06d0e503ef189979a159bece1e8.json", fileBobED25519, 0755)

					return m
				}(),
				nil,
			},
			args{
				"substrate",
				substrate.EdgewareBeresheet,
			},
			map[string]map[string][][]uint8{
				"substrate": {
					"edgeware-beresheet": [][]uint8{
						{0x2a, 0x2e, 0x32, 0x2f, 0x87, 0x40, 0xc6, 0x1, 0x72, 0x11, 0x1a, 0xc8, 0xea, 0xdc, 0xdd, 0xa2, 0x51, 0x2f, 0x90, 0xd0, 0x6d, 0xe, 0x50, 0x3e, 0xf1, 0x89, 0x97, 0x9a, 0x15, 0x9b, 0xec, 0xe1, 0xe8, 0x6d, 0x48},
						{0x2a, 0x72, 0x3c, 0xaa, 0x23, 0xa5, 0xb5, 0x11, 0xaf, 0x5a, 0xd7, 0xb7, 0xef, 0x60, 0x76, 0xe4, 0x14, 0xab, 0x7e, 0x75, 0xa9, 0xdc, 0x91, 0xe, 0xa6, 0xe, 0x41, 0x7a, 0x2b, 0x77, 0xa, 0x56, 0x71, 0x63, 0x83},
					},
				},
			},
			false,
		},
		{
			"success-ethereum",
			fields{
				func() afero.Fs {
					m := afero.NewMemMapFs()
					afero.WriteFile(m, "./ethereum/mainnet/0269d908510e355beb1d5bf2df8129e5b6401e1969891e8016a0b2300739bbb006.json", fileAliceSECP256k1, 0755)
					afero.WriteFile(m, "./ethereum/goerli/0269d908510e355beb1d5bf2df8129e5b6401e1969891e8016a0b2300739bbb006.json", fileAliceSECP256k1, 0755)
					afero.WriteFile(m, "./substrate/edgeware-beresheet/723caa23a5b511af5ad7b7ef6076e414ab7e75a9dc910ea60e417a2b770a5671.json", fileAliceED25519, 0755)
					afero.WriteFile(m, "./substrate/edgeware-beresheet/2e322f8740c60172111ac8eadcdda2512f90d06d0e503ef189979a159bece1e8.json", fileBobED25519, 0755)
					afero.WriteFile(m, "./substrate/edgeware-mainnet/723caa23a5b511af5ad7b7ef6076e414ab7e75a9dc910ea60e417a2b770a5671.json", fileAliceED25519, 0755)
					afero.WriteFile(m, "./substrate/edgeware-mainnet/2e322f8740c60172111ac8eadcdda2512f90d06d0e503ef189979a159bece1e8.json", fileBobED25519, 0755)

					return m
				}(),
				nil,
			},
			args{
				"ethereum",
				"",
			},
			map[string]map[string][][]uint8{
				"ethereum": {
					"goerli":  [][]uint8{{0xd5, 0xab, 0x4c, 0xe3, 0x60, 0x5c, 0xd5, 0x90, 0xdb, 0x60, 0x9b, 0x6b, 0x5c, 0x89, 0x1, 0xfd, 0xb2, 0xef, 0x7f, 0xe6}},
					"kovan":   [][]uint8{},
					"mainnet": [][]uint8{{0xd5, 0xab, 0x4c, 0xe3, 0x60, 0x5c, 0xd5, 0x90, 0xdb, 0x60, 0x9b, 0x6b, 0x5c, 0x89, 0x1, 0xfd, 0xb2, 0xef, 0x7f, 0xe6}},
					"rinkeby": [][]uint8{},
					"ropsten": [][]uint8{}},
			},
			false,
		},
		{
			"success-all",
			fields{
				func() afero.Fs {
					m := afero.NewMemMapFs()
					afero.WriteFile(m, "./ethereum/mainnet/0269d908510e355beb1d5bf2df8129e5b6401e1969891e8016a0b2300739bbb006.json", fileAliceSECP256k1, 0755)
					afero.WriteFile(m, "./ethereum/goerli/0269d908510e355beb1d5bf2df8129e5b6401e1969891e8016a0b2300739bbb006.json", fileAliceSECP256k1, 0755)
					afero.WriteFile(m, "./substrate/edgeware-beresheet/723caa23a5b511af5ad7b7ef6076e414ab7e75a9dc910ea60e417a2b770a5671.json", fileAliceED25519, 0755)
					afero.WriteFile(m, "./substrate/edgeware-beresheet/2e322f8740c60172111ac8eadcdda2512f90d06d0e503ef189979a159bece1e8.json", fileBobED25519, 0755)
					afero.WriteFile(m, "./substrate/edgeware-mainnet/723caa23a5b511af5ad7b7ef6076e414ab7e75a9dc910ea60e417a2b770a5671.json", fileAliceED25519, 0755)
					afero.WriteFile(m, "./substrate/edgeware-mainnet/2e322f8740c60172111ac8eadcdda2512f90d06d0e503ef189979a159bece1e8.json", fileBobED25519, 0755)

					return m
				}(),
				nil,
			},
			args{
				"",
				"",
			},
			map[string]map[string][][]uint8{
				"algorand": {
					"betanet": [][]uint8{}, "mainnet": [][]uint8{}, "testnet": [][]uint8{},
				},
				"ethereum": {
					"goerli":  [][]uint8{{0xd5, 0xab, 0x4c, 0xe3, 0x60, 0x5c, 0xd5, 0x90, 0xdb, 0x60, 0x9b, 0x6b, 0x5c, 0x89, 0x1, 0xfd, 0xb2, 0xef, 0x7f, 0xe6}},
					"kovan":   [][]uint8{},
					"mainnet": [][]uint8{{0xd5, 0xab, 0x4c, 0xe3, 0x60, 0x5c, 0xd5, 0x90, 0xdb, 0x60, 0x9b, 0x6b, 0x5c, 0x89, 0x1, 0xfd, 0xb2, 0xef, 0x7f, 0xe6}},
					"rinkeby": [][]uint8{},
					"ropsten": [][]uint8{}},
				"substrate": {
					"edgeware-beresheet": [][]uint8{
						{0x2a, 0x2e, 0x32, 0x2f, 0x87, 0x40, 0xc6, 0x1, 0x72, 0x11, 0x1a, 0xc8, 0xea, 0xdc, 0xdd, 0xa2, 0x51, 0x2f, 0x90, 0xd0, 0x6d, 0xe, 0x50, 0x3e, 0xf1, 0x89, 0x97, 0x9a, 0x15, 0x9b, 0xec, 0xe1, 0xe8, 0x6d, 0x48},
						{0x2a, 0x72, 0x3c, 0xaa, 0x23, 0xa5, 0xb5, 0x11, 0xaf, 0x5a, 0xd7, 0xb7, 0xef, 0x60, 0x76, 0xe4, 0x14, 0xab, 0x7e, 0x75, 0xa9, 0xdc, 0x91, 0xe, 0xa6, 0xe, 0x41, 0x7a, 0x2b, 0x77, 0xa, 0x56, 0x71, 0x63, 0x83},
					},
					"edgeware-local": [][]uint8{},
					"edgeware-mainnet": [][]uint8{
						{0x7, 0x2e, 0x32, 0x2f, 0x87, 0x40, 0xc6, 0x1, 0x72, 0x11, 0x1a, 0xc8, 0xea, 0xdc, 0xdd, 0xa2, 0x51, 0x2f, 0x90, 0xd0, 0x6d, 0xe, 0x50, 0x3e, 0xf1, 0x89, 0x97, 0x9a, 0x15, 0x9b, 0xec, 0xe1, 0xe8, 0x9b, 0x76},
						{0x7, 0x72, 0x3c, 0xaa, 0x23, 0xa5, 0xb5, 0x11, 0xaf, 0x5a, 0xd7, 0xb7, 0xef, 0x60, 0x76, 0xe4, 0x14, 0xab, 0x7e, 0x75, 0xa9, 0xdc, 0x91, 0xe, 0xa6, 0xe, 0x41, 0x7a, 0x2b, 0x77, 0xa, 0x56, 0x71, 0xda, 0xb},
					},
				},
			},
			false,
		},
		{
			"err-network-only",
			fields{
				func() afero.Fs {
					m := afero.NewMemMapFs()
					afero.WriteFile(m, "./ethereum/mainnet/0269d908510e355beb1d5bf2df8129e5b6401e1969891e8016a0b2300739bbb006.json", fileAliceSECP256k1, 0755)
					afero.WriteFile(m, "./ethereum/goerli/0269d908510e355beb1d5bf2df8129e5b6401e1969891e8016a0b2300739bbb006.json", fileAliceSECP256k1, 0755)
					afero.WriteFile(m, "./substrate/edgeware-beresheet/0269d908510e355beb1d5bf2df8129e5b6401e1969891e8016a0b2300739bbb006.json", fileAliceSECP256k1, 0755)
					afero.WriteFile(m, "./substrate/edgeware-beresheet/723caa23a5b511af5ad7b7ef6076e414ab7e75a9dc910ea60e417a2b770a5671.json", fileBobED25519, 0755)
					afero.WriteFile(m, "./substrate/edgeware-mainnet/1269d908510e355beb1d5bf2df8129e5b6401e1969891e8016a0b2300739bbb006.json", fileAliceSECP256k1, 0755)
					afero.WriteFile(m, "./substrate/edgeware-mainnet/823caa23a5b511af5ad7b7ef6076e414ab7e75a9dc910ea60e417a2b770a5671.json", fileBobED25519, 0755)

					return m
				}(),
				nil,
			},
			args{
				"",
				"mainnet",
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
			got, err := f.GetAddresses(tt.args.protocol, tt.args.network)
			if (err != nil) != tt.wantErr {
				t.Errorf("FileStore.GetAddresses() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("FileStore.GetAddresses() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFileStore_getEncryptedKeyByAddress(t *testing.T) {
	type fields struct {
		fs   afero.Fs
		rand io.Reader
	}
	type args struct {
		searchAddress []byte
		protocol      string
		network       string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *keystore.EncryptedKey
		wantErr bool
	}{
		{
			"success-alice-secp256k1",
			fields{
				func() afero.Fs {
					m := afero.NewMemMapFs()
					afero.WriteFile(m, "./ethereum/mainnet/0269d908510e355beb1d5bf2df8129e5b6401e1969891e8016a0b2300739bbb006.json", fileAliceSECP256k1, 0755)
					afero.WriteFile(m, "./ethereum/mainnet/723caa23a5b511af5ad7b7ef6076e414ab7e75a9dc910ea60e417a2b770a5671.json", fileAliceED25519, 0755)

					return m
				}(),
				nil,
			},
			args{
				[]byte{0xd5, 0xab, 0x4c, 0xe3, 0x60, 0x5c, 0xd5, 0x90, 0xdb, 0x60, 0x9b, 0x6b, 0x5c, 0x89, 0x1, 0xfd, 0xb2, 0xef, 0x7f, 0xe6},
				"ethereum",
				"mainnet",
			},
			&encryptedKeyAliceSECP256k1,
			false,
		},
		{
			"success-alice-ed255419",
			fields{
				func() afero.Fs {
					m := afero.NewMemMapFs()
					afero.WriteFile(m, "./substrate/edgeware-mainnet/0269d908510e355beb1d5bf2df8129e5b6401e1969891e8016a0b2300739bbb006.json", fileAliceSECP256k1, 0755)
					afero.WriteFile(m, "./substrate/edgeware-mainnet/723caa23a5b511af5ad7b7ef6076e414ab7e75a9dc910ea60e417a2b770a5671.json", fileAliceED25519, 0755)

					return m
				}(),
				nil,
			},
			args{
				[]byte{0x7, 0x72, 0x3c, 0xaa, 0x23, 0xa5, 0xb5, 0x11, 0xaf, 0x5a, 0xd7, 0xb7, 0xef, 0x60, 0x76, 0xe4, 0x14, 0xab, 0x7e, 0x75, 0xa9, 0xdc, 0x91, 0xe, 0xa6, 0xe, 0x41, 0x7a, 0x2b, 0x77, 0xa, 0x56, 0x71, 0xda, 0x0b},
				"substrate",
				substrate.EdgewareMainnet,
			},
			&encryptedKeyAliceED25519,
			false,
		},
		{
			"success-bob-ed25519",
			fields{
				func() afero.Fs {
					m := afero.NewMemMapFs()
					// afero.WriteFile(m, "./0269d908510e355beb1d5bf2df8129e5b6401e1969891e8016a0b2300739bbb006.json", fileInvalidCurve, 0755)
					// afero.WriteFile(m, "./723caa23a5b511af5ad7b7ef6076e414ab7e75a9dc910ea60e417a2b770a5671.json", fileAliceED25519, 0755)
					afero.WriteFile(m, "./substrate/edgeware-mainnet/2e322f8740c60172111ac8eadcdda2512f90d06d0e503ef189979a159bece1e8.json", fileBobED25519, 0755)

					return m
				}(),
				nil,
			},
			args{
				[]byte{0x7, 0x2e, 0x32, 0x2f, 0x87, 0x40, 0xc6, 0x1, 0x72, 0x11, 0x1a, 0xc8, 0xea, 0xdc, 0xdd, 0xa2, 0x51, 0x2f, 0x90, 0xd0, 0x6d, 0xe, 0x50, 0x3e, 0xf1, 0x89, 0x97, 0x9a, 0x15, 0x9b, 0xec, 0xe1, 0xe8, 0x9b, 0x76},
				"substrate",
				substrate.EdgewareMainnet,
			},
			&encryptedKeyBobED25519,
			false,
		},
		{
			"err-empty-dir",
			fields{
				func() afero.Fs {
					m := afero.NewMemMapFs()

					return m
				}(),
				nil,
			},
			args{
				[]byte{0xd5, 0xab, 0x4c, 0xe3, 0x60, 0x5c, 0xd5, 0x90, 0xdb, 0x60, 0x9b, 0x6b, 0x5c, 0x89, 0x1, 0xfd, 0xb2, 0xef, 0x7f, 0xe6},
				"ethereum",
				"mainnet",
			},
			nil,
			true,
		},
		{
			"err-invalid-curve",
			fields{
				func() afero.Fs {
					m := afero.NewMemMapFs()
					afero.WriteFile(m, "./ethereum/mainnet/0269d908510e355beb1d5bf2df8129e5b6401e1969891e8016a0b2300739bbb006.json", fileInvalidCurve, 0755)
					afero.WriteFile(m, "./ethereum/mainnet/723caa23a5b511af5ad7b7ef6076e414ab7e75a9dc910ea60e417a2b770a5671.json", fileAliceED25519, 0755)

					return m
				}(),
				nil,
			},
			args{
				[]byte{0xd5, 0xab, 0x4c, 0xe3, 0x60, 0x5c, 0xd5, 0x90, 0xdb, 0x60, 0x9b, 0x6b, 0x5c, 0x89, 0x1, 0xfd, 0xb2, 0xef, 0x7f, 0xe6},
				"ethereum",
				"mainnet",
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
			got, err := f.getEncryptedKeyByAddress(tt.args.searchAddress, tt.args.protocol, tt.args.network)
			if (err != nil) != tt.wantErr {
				t.Errorf("FileStore.getEncryptedKeyByAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("FileStore.getEncryptedKeyByAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFileStore_HasAddress(t *testing.T) {
	type fields struct {
		fs   afero.Fs
		rand io.Reader
	}
	type args struct {
		searchAddress []byte
		protocol      string
		network       string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			"true-alice-secp256k1",
			fields{
				func() afero.Fs {
					m := afero.NewMemMapFs()
					afero.WriteFile(m, "./ethereum/mainnet/0269d908510e355beb1d5bf2df8129e5b6401e1969891e8016a0b2300739bbb006.json", fileAliceSECP256k1, 0755)
					afero.WriteFile(m, "./ethereum/mainnet/723caa23a5b511af5ad7b7ef6076e414ab7e75a9dc910ea60e417a2b770a5671.json", fileAliceED25519, 0755)

					return m
				}(),
				nil,
			},
			args{
				[]byte{0xd5, 0xab, 0x4c, 0xe3, 0x60, 0x5c, 0xd5, 0x90, 0xdb, 0x60, 0x9b, 0x6b, 0x5c, 0x89, 0x1, 0xfd, 0xb2, 0xef, 0x7f, 0xe6},
				"ethereum",
				"mainnet",
			},
			true,
		},
		{
			"true-bob-ed25519",
			fields{
				func() afero.Fs {
					m := afero.NewMemMapFs()
					afero.WriteFile(m, "./substrate/edgeware-mainnet/0269d908510e355beb1d5bf2df8129e5b6401e1969891e8016a0b2300739bbb006.json", fileInvalidCurve, 0755)
					afero.WriteFile(m, "./substrate/edgeware-mainnet/723caa23a5b511af5ad7b7ef6076e414ab7e75a9dc910ea60e417a2b770a5671.json", fileAliceED25519, 0755)
					afero.WriteFile(m, "./substrate/edgeware-mainnet/2e322f8740c60172111ac8eadcdda2512f90d06d0e503ef189979a159bece1e8.json", fileBobED25519, 0755)

					return m
				}(),
				nil,
			},
			args{
				[]byte{0x7, 0x2e, 0x32, 0x2f, 0x87, 0x40, 0xc6, 0x1, 0x72, 0x11, 0x1a, 0xc8, 0xea, 0xdc, 0xdd, 0xa2, 0x51, 0x2f, 0x90, 0xd0, 0x6d, 0xe, 0x50, 0x3e, 0xf1, 0x89, 0x97, 0x9a, 0x15, 0x9b, 0xec, 0xe1, 0xe8, 0x9b, 0x76},
				"substrate",
				substrate.EdgewareMainnet,
			},
			true,
		},
		{
			"false-alice-secp256k1-invalid-key",
			fields{
				func() afero.Fs {
					m := afero.NewMemMapFs()
					afero.WriteFile(m, "./ethereum/mainnet/0269d908510e355beb1d5bf2df8129e5b6401e1969891e8016a0b2300739bbb006.json", fileAliceSECP256k1, 0755)
					afero.WriteFile(m, "./ethereum/mainnet/723caa23a5b511af5ad7b7ef6076e414ab7e75a9dc910ea60e417a2b770a5671.json", fileAliceED25519, 0755)

					return m
				}(),
				nil,
			},
			args{
				[]byte{0xFF, 0xFF, 0x4c, 0xe3, 0x60, 0x5c, 0xd5, 0x90, 0xdb, 0x60, 0x9b, 0x6b, 0x5c, 0x89, 0x1, 0xfd, 0xb2, 0xef, 0x7f, 0xe6},
				"ethereum",
				"mainnet",
			},
			false,
		},
		{
			"false-empty-dir",
			fields{
				func() afero.Fs {
					m := afero.NewMemMapFs()

					return m
				}(),
				nil,
			},
			args{
				[]byte{0xd5, 0xab, 0x4c, 0xe3, 0x60, 0x5c, 0xd5, 0x90, 0xdb, 0x60, 0x9b, 0x6b, 0x5c, 0x89, 0x1, 0xfd, 0xb2, 0xef, 0x7f, 0xe6},
				"ethereum",
				"mainnet",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := FileStore{
				fs:   tt.fields.fs,
				rand: tt.fields.rand,
			}
			if got := f.HasAddress(tt.args.searchAddress, tt.args.protocol, tt.args.network); got != tt.want {
				t.Errorf("FileStore.HasAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFileStore_getProtocolNetworkAddresses(t *testing.T) {
	type fields struct {
		fs   afero.Fs
		rand io.Reader
	}
	type args struct {
		protocol string
		network  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    [][]byte
		wantErr bool
	}{
		{
			"multiple-addresses",
			fields{
				func() afero.Fs {
					m := afero.NewMemMapFs()
					afero.WriteFile(m, "./ethereum/mainnet/0269d908510e355beb1d5bf2df8129e5b6401e1969891e8016a0b2300739bbb006.json", fileAliceSECP256k1, 0755)
					afero.WriteFile(m, "./ethereum/goerli/0269d908510e355beb1d5bf2df8129e5b6401e1969891e8016a0b2300739bbb006.json", fileAliceSECP256k1, 0755)
					afero.WriteFile(m, "./substrate/edgeware-beresheet/723caa23a5b511af5ad7b7ef6076e414ab7e75a9dc910ea60e417a2b770a5671.json", fileAliceED25519, 0755)
					afero.WriteFile(m, "./substrate/edgeware-beresheet/2e322f8740c60172111ac8eadcdda2512f90d06d0e503ef189979a159bece1e8.json", fileBobED25519, 0755)
					afero.WriteFile(m, "./substrate/edgeware-mainnet/723caa23a5b511af5ad7b7ef6076e414ab7e75a9dc910ea60e417a2b770a5671.json", fileAliceED25519, 0755)
					afero.WriteFile(m, "./substrate/edgeware-mainnet/2e322f8740c60172111ac8eadcdda2512f90d06d0e503ef189979a159bece1e8.json", fileBobED25519, 0755)

					return m
				}(),
				nil,
			},
			args{
				"substrate",
				"edgeware-beresheet",
			},
			[][]uint8{
				{0x2a, 0x2e, 0x32, 0x2f, 0x87, 0x40, 0xc6, 0x1, 0x72, 0x11, 0x1a, 0xc8, 0xea, 0xdc, 0xdd, 0xa2, 0x51, 0x2f, 0x90, 0xd0, 0x6d, 0xe, 0x50, 0x3e, 0xf1, 0x89, 0x97, 0x9a, 0x15, 0x9b, 0xec, 0xe1, 0xe8, 0x6d, 0x48},
				{0x2a, 0x72, 0x3c, 0xaa, 0x23, 0xa5, 0xb5, 0x11, 0xaf, 0x5a, 0xd7, 0xb7, 0xef, 0x60, 0x76, 0xe4, 0x14, 0xab, 0x7e, 0x75, 0xa9, 0xdc, 0x91, 0xe, 0xa6, 0xe, 0x41, 0x7a, 0x2b, 0x77, 0xa, 0x56, 0x71, 0x63, 0x83},
			},
			false,
		},
		{
			"single-address",
			fields{
				func() afero.Fs {
					m := afero.NewMemMapFs()
					afero.WriteFile(m, "./ethereum/mainnet/0269d908510e355beb1d5bf2df8129e5b6401e1969891e8016a0b2300739bbb006.json", fileAliceSECP256k1, 0755)
					afero.WriteFile(m, "./ethereum/goerli/0269d908510e355beb1d5bf2df8129e5b6401e1969891e8016a0b2300739bbb006.json", fileAliceSECP256k1, 0755)
					afero.WriteFile(m, "./substrate/edgeware-beresheet/723caa23a5b511af5ad7b7ef6076e414ab7e75a9dc910ea60e417a2b770a5671.json", fileAliceED25519, 0755)
					afero.WriteFile(m, "./substrate/edgeware-beresheet/2e322f8740c60172111ac8eadcdda2512f90d06d0e503ef189979a159bece1e8.json", fileBobED25519, 0755)
					afero.WriteFile(m, "./substrate/edgeware-mainnet/723caa23a5b511af5ad7b7ef6076e414ab7e75a9dc910ea60e417a2b770a5671.json", fileAliceED25519, 0755)
					afero.WriteFile(m, "./substrate/edgeware-mainnet/2e322f8740c60172111ac8eadcdda2512f90d06d0e503ef189979a159bece1e8.json", fileBobED25519, 0755)

					return m
				}(),
				nil,
			},
			args{
				"ethereum",
				"mainnet",
			},
			[][]uint8{{0xd5, 0xab, 0x4c, 0xe3, 0x60, 0x5c, 0xd5, 0x90, 0xdb, 0x60, 0x9b, 0x6b, 0x5c, 0x89, 0x1, 0xfd, 0xb2, 0xef, 0x7f, 0xe6}},
			false,
		},
		{
			"no-addresses",
			fields{
				func() afero.Fs {
					m := afero.NewMemMapFs()
					afero.WriteFile(m, "./ethereum/mainnet/0269d908510e355beb1d5bf2df8129e5b6401e1969891e8016a0b2300739bbb006.json", fileAliceSECP256k1, 0755)
					afero.WriteFile(m, "./ethereum/goerli/0269d908510e355beb1d5bf2df8129e5b6401e1969891e8016a0b2300739bbb006.json", fileAliceSECP256k1, 0755)
					afero.WriteFile(m, "./substrate/edgeware-beresheet/723caa23a5b511af5ad7b7ef6076e414ab7e75a9dc910ea60e417a2b770a5671.json", fileAliceED25519, 0755)
					afero.WriteFile(m, "./substrate/edgeware-beresheet/2e322f8740c60172111ac8eadcdda2512f90d06d0e503ef189979a159bece1e8.json", fileBobED25519, 0755)
					afero.WriteFile(m, "./substrate/edgeware-mainnet/723caa23a5b511af5ad7b7ef6076e414ab7e75a9dc910ea60e417a2b770a5671.json", fileAliceED25519, 0755)
					afero.WriteFile(m, "./substrate/edgeware-mainnet/2e322f8740c60172111ac8eadcdda2512f90d06d0e503ef189979a159bece1e8.json", fileBobED25519, 0755)

					return m
				}(),
				nil,
			},
			args{
				"ethereum",
				"rinkeby",
			},
			[][]byte{},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := FileStore{
				fs:   tt.fields.fs,
				rand: tt.fields.rand,
			}
			got, err := f.getProtocolNetworkAddresses(tt.args.protocol, tt.args.network)
			if (err != nil) != tt.wantErr {
				t.Errorf("FileStore.getProtocolNetworkAddresses() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("FileStore.getProtocolNetworkAddresses() = %v, want %v", got, tt.want)
			}
		})
	}
}
