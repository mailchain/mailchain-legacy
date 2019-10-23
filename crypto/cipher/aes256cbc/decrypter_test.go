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
	"testing"

	"github.com/mailchain/mailchain/crypto/cipher"
	"github.com/mailchain/mailchain/crypto/secp256k1/secp256k1test"
	"github.com/mailchain/mailchain/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestDecryptCBC(t *testing.T) {
	assert := assert.New(t)
	cases := []struct {
		name          string
		iv            []byte
		encryptionKey []byte
		ciphertext    []byte
		expected      []byte
		err           error
	}{
		{"short text",
			testutil.MustHexDecodeString("05050505050505050505050505050505"),
			testutil.MustHexDecodeString("af0ad81e7d9194721d6c26f6c1f2a2b7fd06e2c99c4f5deefe59fb93936c981e"),
			testutil.MustHexDecodeString("747ef78a32eb582d325a634e4acffd61"),
			[]byte("Hi Tim"),
			nil,
		}, {"medium text",
			testutil.MustHexDecodeString("05050505050505050505050505050505"),
			testutil.MustHexDecodeString("af0ad81e7d9194721d6c26f6c1f2a2b7fd06e2c99c4f5deefe59fb93936c981e"),
			testutil.MustHexDecodeString("2ec66aac453ff543f47830d4b8cbc68d9965bf7c6bb69724fd4de26d41001256dfa6f7f0b3956ce21d4717caf75b0c2ad753852f216df6cfbcda4911619c5fc34798a19f81adff902c1ad906ab0edaec"),
			[]byte("Hi Tim, this is a much longer message to make sure there are no problems"),
			nil,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := decryptCBC(tc.encryptionKey, tc.iv, tc.ciphertext)
			assert.EqualValues(string(tc.expected), string(res))
			assert.Equal(tc.err, err)
		})
	}
}

func TestDecrypter(t *testing.T) {
	assert := assert.New(t)
	cases := []struct {
		name          string
		decrypter     Decrypter
		encryptedData cipher.EncryptedContent
		expected      []byte
		err           error
	}{
		{"to-sofia-short-text",
			Decrypter{privateKey: secp256k1test.SofiaPrivateKey},
			testutil.MustHexDecodeString("2efafed52c39d4cd1ef32db24d015e77e002882fde5ee55e0b49d2b84ccd1bbe19ee705b355421e9f7e88edbce5b4b1da6ba08ce4c5adca03ea8f14b45ac09a0c2536b70c0a72cc01a5310b240508ff2cbded2b74094e6d302707b324e43ace545e2"),
			[]byte("Hi Sofia"),
			nil,
		}, {"to-sofia-medium-text",
			Decrypter{privateKey: secp256k1test.SofiaPrivateKey},
			testutil.MustHexDecodeString("2ee80337544404fb07b06cd19515fcd635038621ccd7fb04c7e2771ea1b6ccd3dcacba071bdd92a51ee443c3735a09c32c4d5523950142380c9ac771d95eece221470a3a52db70060ff43dbea5d3891da942c14f515f9f0bc9fa1f7cc25327b8b67668f75a266de53ec000e04375fab30c67adced120090c5dcb7b3bfde491239bb3c6a444aff4610c905f1e8ec82ea0d9f54a6d31e195b4f784e150762be160f52732b99eb503ba42da8eb2baef63adc931"),
			[]byte("Hi Sofia, this is a little bit of a longer message to make sure there are no problems"),
			nil,
		}, {"to-sofia-short-text",
			Decrypter{privateKey: secp256k1test.CharlottePrivateKey},
			testutil.MustHexDecodeString("2e5e33a1a6013a268a7494ffc27a27a5a903c1715fed0ebf62975efea659152054adb11d5e4c43f4894cd5a233da10f33175e4afa699faf96dccfe551e4a11a4d8334f8bce4ec04e3119852b283b55ce208148ce7190b8affe0c6b5b6430bb749576"),
			[]byte("Hi Charlotte"),
			nil,
		}, {"to-sofia-short-text",
			Decrypter{privateKey: secp256k1test.CharlottePrivateKey},
			testutil.MustHexDecodeString("2e7ffbd8a1092e8cc8f6a0a5df4a23ea0b02fba947d56d83a9d6801233ce074ed936811d9d541be757b8b39bdc9d0ff08424b2bd1c5dc26e5789ef4fa5bfb24e074acd8c0a324d3a0598bde06f4d51ac1bc91eed372e93bb23891840882eac1008f56ed77171de248dd7ed0dd33f8ed6b4d9fde39ed690aea3d33f9a4dd09f645fea625603e5ad06c44ced1c2c44d1a2c895b91f5dd8d1a1d540f8f9e88087c3466c359c4e6a8347ee54dbfdd055ede2d29e"),
			[]byte("Hi Charlotte, this is a little bit of a longer message to make sure there are no problems"),
			nil,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := tc.decrypter.Decrypt(tc.encryptedData)
			assert.EqualValues(string(tc.expected), string(res))
			assert.Equal(tc.err, err)
		})
	}
}
