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

package testutil

import (
	"log"

	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/secp256k1"
)

var SofiaPrivateKey crypto.PrivateKey     //nolint: gochecknoglobals
var SofiaPublicKey crypto.PublicKey       //nolint: gochecknoglobals
var CharlottePrivateKey crypto.PrivateKey //nolint: gochecknoglobals
var CharlottePublicKey crypto.PublicKey   //nolint: gochecknoglobals

//nolint: gochecknoinits
func init() {
	var err error
	SofiaPrivateKey, err = secp256k1.PrivateKeyFromHex("01901E63389EF02EAA7C5782E08B40D98FAEF835F28BD144EECF5614A415943F")
	if err != nil {
		log.Fatal(err)
	}
	SofiaPublicKey = SofiaPrivateKey.PublicKey()

	CharlottePrivateKey, err = secp256k1.PrivateKeyFromHex("DF4BA9F6106AD2846472F759476535E55C5805D8337DF5A11C3B139F438B98B3")
	if err != nil {
		log.Fatal(err)
	}

	CharlottePublicKey = CharlottePrivateKey.PublicKey()
}
