package secp256k1test

import (
	"log"

	"github.com/ethereum/go-ethereum/common/hexutil"
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

	encryptedSofiaPrivateKey, err := hexutil.Decode("01901E63389EF02EAA7C5782E08B40D98FAEF835F28BD144EECF5614A415943F")
	if err != nil {
		log.Fatal(err)
	}
	SofiaPrivateKey, err = secp256k1.PrivateKeyFromBytes(encryptedSofiaPrivateKey)
	if err != nil {
		log.Fatal(err)
	}
	SofiaPublicKey = SofiaPrivateKey.PublicKey()

	CharlottePrivateKey, err = secp256k1.PrivateKeyFromBytes(encryptedSofiaPrivateKey)
	if err != nil {
		log.Fatal(err)
	}

	CharlottePublicKey = CharlottePrivateKey.PublicKey()
}
