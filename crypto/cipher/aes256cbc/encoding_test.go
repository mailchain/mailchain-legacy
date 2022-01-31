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

package aes256cbc

import (
	"bytes"
	"testing"

	"github.com/mailchain/mailchain/encoding/encodingtest"
	"github.com/stretchr/testify/assert"
)

func TestEncodeDecode(t *testing.T) {
	cases := []struct {
		name        string
		encodedData []byte
		err         error
	}{
		{"tc1",
			encodingtest.MustDecodeHex("2eb3d72325f94ed8b9e1b7f28e2fb26492029dce5444ad23a68a76dd1821b9b2b3a9c6e53d464420e2363a80df94cc7b05f58412f3436593821021308c64d4d18482d224e79b9cb2b14b177214f3b023ebe6a6537a3781ed4927228bd7d80d1d6f07"),
			nil,
		},
		{"tc2",
			encodingtest.MustDecodeHex("2ef8307114bb283da496056a8502376cdf0287a2cd646044a0f9639aa3b50aa26b170f21fbedd20e079ab890d3a9c880dea458b3398eccbfeaaa08b350c6226e984a7e70a04f8a97c07f0f5a8e9a36394cf19110ac2e87fcbe9c73faf41183d23a27"),
			nil,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			decoded, err := bytesDecode(tc.encodedData)
			assert.NoError(t, err)

			encoded, err := bytesEncode(decoded)
			assert.NoError(t, err)
			assert.EqualValues(t, tc.encodedData, encoded)
			assert.True(t, bytes.Equal(tc.encodedData, encoded))
			assert.Equal(t, tc.err, err)
		})
	}
}

func TestDecodeEncode(t *testing.T) {
	cases := []struct {
		name          string
		encryptedData *encryptedData
		err           error
	}{
		{"tc1",
			&encryptedData{
				Ciphertext:                encodingtest.MustDecodeHex("a6537a3781ed4927228bd7d80d1d6f07"),
				EphemeralPublicKey:        encodingtest.MustDecodeHex("049dce5444ad23a68a76dd1821b9b2b3a9c6e53d464420e2363a80df94cc7b05f5c0896985fc8156846a42d1b922f253e1e2537b9279cafe44bce66552cbc58c04"),
				InitializationVector:      encodingtest.MustDecodeHex("b3d72325f94ed8b9e1b7f28e2fb26492"),
				MessageAuthenticationCode: encodingtest.MustDecodeHex("8412f3436593821021308c64d4d18482d224e79b9cb2b14b177214f3b023ebe6"),
			},
			nil,
		},
		{"tc2",
			&encryptedData{
				Ciphertext:                encodingtest.MustDecodeHex("9110ac2e87fcbe9c73faf41183d23a27"),
				EphemeralPublicKey:        encodingtest.MustDecodeHex("0487a2cd646044a0f9639aa3b50aa26b170f21fbedd20e079ab890d3a9c880dea4cbdaab93155fa43441dca3e7e94dc2ff67882ec4908e82b0496821cffb4d7cc8"),
				InitializationVector:      encodingtest.MustDecodeHex("f8307114bb283da496056a8502376cdf"),
				MessageAuthenticationCode: encodingtest.MustDecodeHex("58b3398eccbfeaaa08b350c6226e984a7e70a04f8a97c07f0f5a8e9a36394cf1"),
			},
			nil,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			encoded, err := bytesEncode(tc.encryptedData)
			assert.NoError(t, err)
			actual, err := bytesDecode(encoded)
			assert.NoError(t, err)
			assert.EqualValues(t, len(tc.encryptedData.Ciphertext), cap(actual.Ciphertext))
			assert.EqualValues(t, len(tc.encryptedData.InitializationVector), cap(actual.InitializationVector))
			assert.EqualValues(t, len(tc.encryptedData.MessageAuthenticationCode), cap(actual.MessageAuthenticationCode))
			assert.EqualValues(t, len(tc.encryptedData.EphemeralPublicKey), cap(actual.EphemeralPublicKey))
			assert.EqualValues(t, tc.encryptedData, actual)
			assert.Equal(t, tc.err, err)
		})
	}
}

func Test_bytesDecode(t *testing.T) {
	type args struct {
		raw []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *encryptedData
		wantErr bool
	}{
		{
			"success-testcase-1",
			args{
				encodingtest.MustDecodeHex("2eb3d72325f94ed8b9e1b7f28e2fb26492029dce5444ad23a68a76dd1821b9b2b3a9c6e53d464420e2363a80df94cc7b05f58412f3436593821021308c64d4d18482d224e79b9cb2b14b177214f3b023ebe6a6537a3781ed4927228bd7d80d1d6f07"),
			},
			&encryptedData{
				Ciphertext:                encodingtest.MustDecodeHex("a6537a3781ed4927228bd7d80d1d6f07"),
				EphemeralPublicKey:        encodingtest.MustDecodeHex("049dce5444ad23a68a76dd1821b9b2b3a9c6e53d464420e2363a80df94cc7b05f5c0896985fc8156846a42d1b922f253e1e2537b9279cafe44bce66552cbc58c04"),
				InitializationVector:      encodingtest.MustDecodeHex("b3d72325f94ed8b9e1b7f28e2fb26492"),
				MessageAuthenticationCode: encodingtest.MustDecodeHex("8412f3436593821021308c64d4d18482d224e79b9cb2b14b177214f3b023ebe6"),
			},
			false,
		},
		{
			"success-testcase-2",
			args{
				encodingtest.MustDecodeHex("2ef8307114bb283da496056a8502376cdf0287a2cd646044a0f9639aa3b50aa26b170f21fbedd20e079ab890d3a9c880dea458b3398eccbfeaaa08b350c6226e984a7e70a04f8a97c07f0f5a8e9a36394cf19110ac2e87fcbe9c73faf41183d23a27"),
			},
			&encryptedData{
				Ciphertext:                encodingtest.MustDecodeHex("9110ac2e87fcbe9c73faf41183d23a27"),
				EphemeralPublicKey:        encodingtest.MustDecodeHex("0487a2cd646044a0f9639aa3b50aa26b170f21fbedd20e079ab890d3a9c880dea4cbdaab93155fa43441dca3e7e94dc2ff67882ec4908e82b0496821cffb4d7cc8"),
				InitializationVector:      encodingtest.MustDecodeHex("f8307114bb283da496056a8502376cdf"),
				MessageAuthenticationCode: encodingtest.MustDecodeHex("58b3398eccbfeaaa08b350c6226e984a7e70a04f8a97c07f0f5a8e9a36394cf1"),
			},
			false,
		},
		{
			"err-empty",
			args{
				[]byte{},
			},
			nil,
			true,
		},
		{
			"err-too-short",
			args{
				[]byte{0x01, 0x02, 0x03},
			},
			nil,
			true,
		},
		{
			"err-invalid-prefix",
			args{
				encodingtest.MustDecodeHex("fff8307114bb283da496056a8502376cdf0287a2cd646044a0f9639aa3b50aa26b170f21fbedd20e079ab890d3a9c880dea458b3398eccbfeaaa08b350c6226e984a7e70a04f8a97c07f0f5a8e9a36394cf19110ac2e87fcbe9c73faf41183d23a27"),
			},
			nil,
			true,
		},
		{
			"err-no-cipher-text",
			args{
				encodingtest.MustDecodeHex("2ef8307114bb283da496056a8502376cdf0287a2cd646044a0f9639aa3b50aa26b170f21fbedd20e079ab890d3a9c880dea458b3398eccbfeaaa08b350c6226e984a7e70a04f8a97c07f0f5a8e9a36394cf1"),
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := bytesDecode(tt.args.raw)
			if (err != nil) != tt.wantErr {
				t.Errorf("bytesDecode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("bytesDecode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_bytesEncode(t *testing.T) {
	type args struct {
		data *encryptedData
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			"success-testcase-1",
			args{
				&encryptedData{
					Ciphertext:                encodingtest.MustDecodeHex("a6537a3781ed4927228bd7d80d1d6f07"),
					EphemeralPublicKey:        encodingtest.MustDecodeHex("049dce5444ad23a68a76dd1821b9b2b3a9c6e53d464420e2363a80df94cc7b05f5c0896985fc8156846a42d1b922f253e1e2537b9279cafe44bce66552cbc58c04"),
					InitializationVector:      encodingtest.MustDecodeHex("b3d72325f94ed8b9e1b7f28e2fb26492"),
					MessageAuthenticationCode: encodingtest.MustDecodeHex("8412f3436593821021308c64d4d18482d224e79b9cb2b14b177214f3b023ebe6"),
				},
			},
			encodingtest.MustDecodeHex("2eb3d72325f94ed8b9e1b7f28e2fb26492029dce5444ad23a68a76dd1821b9b2b3a9c6e53d464420e2363a80df94cc7b05f58412f3436593821021308c64d4d18482d224e79b9cb2b14b177214f3b023ebe6a6537a3781ed4927228bd7d80d1d6f07"),
			false,
		},
		{
			"success-testcase-1",
			args{
				&encryptedData{
					Ciphertext:                encodingtest.MustDecodeHex("9110ac2e87fcbe9c73faf41183d23a27"),
					EphemeralPublicKey:        encodingtest.MustDecodeHex("0487a2cd646044a0f9639aa3b50aa26b170f21fbedd20e079ab890d3a9c880dea4cbdaab93155fa43441dca3e7e94dc2ff67882ec4908e82b0496821cffb4d7cc8"),
					InitializationVector:      encodingtest.MustDecodeHex("f8307114bb283da496056a8502376cdf"),
					MessageAuthenticationCode: encodingtest.MustDecodeHex("58b3398eccbfeaaa08b350c6226e984a7e70a04f8a97c07f0f5a8e9a36394cf1"),
				},
			},
			encodingtest.MustDecodeHex("2ef8307114bb283da496056a8502376cdf0287a2cd646044a0f9639aa3b50aa26b170f21fbedd20e079ab890d3a9c880dea458b3398eccbfeaaa08b350c6226e984a7e70a04f8a97c07f0f5a8e9a36394cf19110ac2e87fcbe9c73faf41183d23a27"),
			false,
		},
		{
			"err-encrypted-data-invalid",
			args{
				&encryptedData{
					Ciphertext:                nil,
					EphemeralPublicKey:        nil,
					InitializationVector:      nil,
					MessageAuthenticationCode: nil,
				},
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := bytesEncode(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("bytesEncode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("bytesEncode() = %v, want %v", got, tt.want)
			}
		})
	}
}
