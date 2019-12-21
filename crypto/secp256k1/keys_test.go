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
	"log"

	"github.com/mailchain/mailchain/encoding/encodingtest"
)

var (
	sofiaPrivateKey = func() PrivateKey {
		k, err := PrivateKeyFromBytes(sofiaPrivateKeyBytes)
		if err != nil {
			log.Fatal(err)
		}
		return *k
	}() //nolint: lll

	sofiaPrivateKeyBytes = encodingtest.MustDecodeHex("01901E63389EF02EAA7C5782E08B40D98FAEF835F28BD144EECF5614A415943F")
	sofiaPublicKey       = func() PublicKey {
		k, err := PublicKeyFromBytes(sofiaPublicKeyBytes)
		if err != nil {
			log.Fatal(err)
		}

		return *k.(*PublicKey)
	}()
	sofiaPublicKeyBytes = encodingtest.MustDecodeHexZeroX("0x69d908510e355beb1d5bf2df8129e5b6401e1969891e8016a0b2300739bbb00687055e5924a2fd8dd35f069dc14d8147aa11c1f7e2f271573487e1beeb2be9d0") //nolint: lll
)

var (
	charlottePrivateKey = func() PrivateKey {
		k, err := PrivateKeyFromBytes(charlottePrivateKeyBytes)
		if err != nil {
			log.Fatal(err)
		}
		return *k
	}() //nolint: lll

	charlottePrivateKeyBytes = encodingtest.MustDecodeHex("DF4BA9F6106AD2846472F759476535E55C5805D8337DF5A11C3B139F438B98B3")
	charlottePublicKey       = func() PublicKey {
		k, err := PublicKeyFromBytes(charlottePublicKeyBytes)
		if err != nil {
			log.Fatal(err)
		}
		return *k.(*PublicKey)
	}()
	charlottePublicKeyBytes = encodingtest.MustDecodeHexZeroX("0xbdf6fb97c97c126b492186a4d5b28f34f0671a5aacc974da3bde0be93e45a1c50f89ceff72bd04ac9e25a04a1a6cb010aedaf65f91cec8ebe75901c49b63355d") //nolint: lll
)
