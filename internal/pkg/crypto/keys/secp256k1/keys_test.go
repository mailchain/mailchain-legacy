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
