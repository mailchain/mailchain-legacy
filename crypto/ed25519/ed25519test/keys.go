package ed25519test

import (
	"log"
	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/ed25519"
	"github.com/mailchain/mailchain/internal/testutil"
)

var SofiaPrivateKey crypto.PrivateKey
var SofiaPublicKey crypto.PublicKey
var CharlottePrivateKey crypto.PrivateKey
var CharlottePublicKey crypto.PublicKey

func init() {
	var err error
	SofiaPrivateKey, err = ed25519.PrivateKeyFromSeed(testutil.MustHexDecodeString("0d9b4a3c10721991c6b806f0f343535dc2b46c74bece50a0a0d6b9f0070d3157"))
	if err != nil {
		log.Fatal(err)
	}

	SofiaPublicKey = SofiaPrivateKey.PublicKey()

	CharlottePrivateKey, err = ed25519.PrivateKeyFromSeed(testutil.MustHexDecodeString("39d4c97d6a7f9e3306a2b5aae604ee67ec8b1387fffb39128fc055656cff05bb"))
	if err != nil {
		log.Fatal(err)
	}

	CharlottePublicKey = CharlottePrivateKey.PublicKey()
}
