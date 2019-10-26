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

	"github.com/mailchain/mailchain/crypto"
)

func sofiaPrivateKey() crypto.PrivateKey {
	k, err := PrivateKeyFromHex("01901E63389EF02EAA7C5782E08B40D98FAEF835F28BD144EECF5614A415943F")
	if err != nil {
		log.Fatal(err)
	}
	return k
}

func sofiaPublicKey() crypto.PublicKey {
	return sofiaPrivateKey().PublicKey()
}

func charlottePrivateKey() crypto.PrivateKey {
	k, err := PrivateKeyFromHex("DF4BA9F6106AD2846472F759476535E55C5805D8337DF5A11C3B139F438B98B3")
	if err != nil {
		log.Fatal(err)
	}
	return k
}

func charlottePublicKey() crypto.PublicKey {
	return charlottePrivateKey().PublicKey()
}
