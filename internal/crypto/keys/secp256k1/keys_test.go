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

package secp256k1_test

import (
	"crypto/ecdsa"
	"encoding/hex"
	"log"

	"github.com/ethereum/go-ethereum/crypto"
)

var privateKeyA *ecdsa.PrivateKey
var publicKeyA ecdsa.PublicKey
var privateKeyB *ecdsa.PrivateKey
var publicKeyB ecdsa.PublicKey

func init() {
	var err error
	pkAHex, _ := hex.DecodeString("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	privateKeyA, err = crypto.ToECDSA(pkAHex)
	if err != nil {
		log.Fatal(err)
	}
	publicKeyA = privateKeyA.PublicKey

	pkBHex, _ := hex.DecodeString("BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB")
	privateKeyB, err = crypto.ToECDSA(pkBHex)
	if err != nil {
		log.Fatal(err)
	}

	publicKeyB = privateKeyB.PublicKey
}
