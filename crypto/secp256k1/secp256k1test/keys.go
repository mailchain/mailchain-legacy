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

package secp256k1test

import (
	"crypto/ecdsa"
	"encoding/hex"
	"log"

	"github.com/ethereum/go-ethereum/crypto"
)

func ECDSAPrivateKeyA() ecdsa.PrivateKey {
	b, _ := hex.DecodeString("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	key, err := crypto.ToECDSA(b)
	if err != nil {
		log.Fatal(err)
	}
	return *key
}

func ECDSAPublicKeyA() ecdsa.PublicKey {
	return ECDSAPrivateKeyA().PublicKey
}

func ECDSAPrivateKeyB() ecdsa.PrivateKey {
	b, _ := hex.DecodeString("BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB")
	key, err := crypto.ToECDSA(b)
	if err != nil {
		log.Fatal(err)
	}
	return *key
}

func ECDSAPublicKeyB() ecdsa.PublicKey {
	return ECDSAPrivateKeyB().PublicKey
}
