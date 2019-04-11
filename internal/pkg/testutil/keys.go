package testutil

import (
	"log"

	"github.com/mailchain/mailchain/internal/pkg/crypto/keys"
	"github.com/mailchain/mailchain/internal/pkg/crypto/keys/secp256k1"
)

var SofiaPrivateKey keys.PrivateKey
var SofiaPublicKey keys.PublicKey
var CharlottePrivateKey keys.PrivateKey
var CharlottePublicKey keys.PublicKey

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
