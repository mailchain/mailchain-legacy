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

package aes256cbc

import (
	"encoding/hex"
	"log"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/mailchain/mailchain/internal/pkg/crypto/keys/secp256k1"
	"github.com/mailchain/mailchain/internal/pkg/testutil"
	"github.com/stretchr/testify/assert"
)

func TestEncryptCBC(t *testing.T) {
	assert := assert.New(t)
	cases := []struct {
		name          string
		iv            []byte
		encryptionKey []byte
		data          []byte
		expected      []byte
		err           error
	}{
		{"short text",
			testutil.MustHexDecodeString("05050505050505050505050505050505"),
			testutil.MustHexDecodeString("af0ad81e7d9194721d6c26f6c1f2a2b7fd06e2c99c4f5deefe59fb93936c981e"),
			[]byte("Hi Tim"),
			testutil.MustHexDecodeString("747ef78a32eb582d325a634e4acffd61"),
			nil,
		}, {"medium text",
			testutil.MustHexDecodeString("05050505050505050505050505050505"),
			testutil.MustHexDecodeString("af0ad81e7d9194721d6c26f6c1f2a2b7fd06e2c99c4f5deefe59fb93936c981e"),
			[]byte("Hi Tim, this is a much longer message to make sure there are no problems"),
			testutil.MustHexDecodeString("2ec66aac453ff543f47830d4b8cbc68d9965bf7c6bb69724fd4de26d41001256dfa6f7f0b3956ce21d4717caf75b0c2ad753852f216df6cfbcda4911619c5fc34798a19f81adff902c1ad906ab0edaec"),
			nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := encryptCBC(tc.data, tc.iv, tc.encryptionKey)
			assert.EqualValues(tc.expected, res)
			assert.Equal(tc.err, err)
		})
	}
}

func TestInternalEncrypt(t *testing.T) {
	assert := assert.New(t)
	iv := testutil.MustHexDecodeString("05050505050505050505050505050505")
	tmpEphemeralPrivateKey, err := crypto.HexToECDSA("0404040404040404040404040404040404040404040404040404040404040404")
	if err != nil {
		log.Fatal(err)
	}
	ephemeralPrivateKey := ecies.ImportECDSA(tmpEphemeralPrivateKey)

	pub, err := secp256k1.PublicKeyToECIES(testutil.SofiaPublicKey)
	if err != nil {
		log.Fatal(err)
	}

	actual, err := encrypt(ephemeralPrivateKey, pub, []byte("Hi Tim, this is a much longer message to make sure there are no problems"), iv)
	if err != nil {
		log.Fatal(err)
	}
	assert.NoError(err)
	assert.NotNil(actual)
	assert.Equal("2ec66aac453ff543f47830d4b8cbc68d9965bf7c6bb69724fd4de26d41001256dfa6f7f0b3956ce21d4717caf75b0c2ad753852f216df6cfbcda4911619c5fc34798a19f81adff902c1ad906ab0edaec", hex.EncodeToString(actual.Ciphertext))
	assert.Equal("04462779ad4aad39514614751a71085f2f10e1c7a593e4e030efb5b8721ce55b0b199c07969f5442000bea455d72ae826a86bfac9089cb18152ed756ebb2a596f5", hex.EncodeToString(actual.EphemeralPublicKey))
	assert.Equal("05050505050505050505050505050505", hex.EncodeToString(actual.InitializationVector))
	assert.Equal("4367ae8a54b65f99e4f2fd315ba65bf85e1138967a7bea451faf80f75cdf3404", hex.EncodeToString(actual.MessageAuthenticationCode))
}
