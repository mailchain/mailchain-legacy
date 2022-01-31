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

package secp256k1

import (
	"log"

	"github.com/mailchain/mailchain/encoding/encodingtest"
)

var (
	alicePrivateKey = func() PrivateKey {
		k, err := PrivateKeyFromBytes(alicePrivateKeyBytes)
		if err != nil {
			log.Fatal(err)
		}
		return *k
	}() //nolint: lll

	alicePrivateKeyBytes = encodingtest.MustDecodeHex("01901E63389EF02EAA7C5782E08B40D98FAEF835F28BD144EECF5614A415943F")
	alicePublicKey       = func() PublicKey {
		k, err := PublicKeyFromBytes(alicePublicKeyBytes)
		if err != nil {
			log.Fatal(err)
		}

		return *k.(*PublicKey)
	}()
	alicePublicKeyBytes = []byte{0x2, 0x69, 0xd9, 0x8, 0x51, 0xe, 0x35, 0x5b, 0xeb, 0x1d, 0x5b, 0xf2, 0xdf, 0x81, 0x29, 0xe5, 0xb6, 0x40, 0x1e, 0x19, 0x69, 0x89, 0x1e, 0x80, 0x16, 0xa0, 0xb2, 0x30, 0x7, 0x39, 0xbb, 0xb0, 0x6}
)

var (
	bobPrivateKey = func() PrivateKey {
		k, err := PrivateKeyFromBytes(bobPrivateKeyBytes)
		if err != nil {
			log.Fatal(err)
		}
		return *k
	}() //nolint: lll

	bobPrivateKeyBytes = encodingtest.MustDecodeHex("DF4BA9F6106AD2846472F759476535E55C5805D8337DF5A11C3B139F438B98B3")
	bobPublicKey       = func() PublicKey {
		k, err := PublicKeyFromBytes(bobPublicKeyBytes)
		if err != nil {
			log.Fatal(err)
		}
		return *k.(*PublicKey)
	}()
	bobPublicKeyBytes = []byte{0x3, 0xbd, 0xf6, 0xfb, 0x97, 0xc9, 0x7c, 0x12, 0x6b, 0x49, 0x21, 0x86, 0xa4, 0xd5, 0xb2, 0x8f, 0x34, 0xf0, 0x67, 0x1a, 0x5a, 0xac, 0xc9, 0x74, 0xda, 0x3b, 0xde, 0xb, 0xe9, 0x3e, 0x45, 0xa1, 0xc5}
)
