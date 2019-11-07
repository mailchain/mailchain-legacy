package secp256k1test

import (
	"encoding/hex"
	"log"

	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/secp256k1"
)

// SofiaPrivateKey secp256k1 key for testing purposes. Key is compromised do not use on mainnet's.
var SofiaPrivateKey crypto.PrivateKey //nolint: gochecknoglobals
// SofiaPublicKey secp256k1 key for testing purposes. Key is compromised do not use on mainnet's.
var SofiaPublicKey crypto.PublicKey //nolint: gochecknoglobals
// CharlottePrivateKey secp256k1 key for testing purposes. Key is compromised do not use on mainnet's.
var CharlottePrivateKey crypto.PrivateKey //nolint: gochecknoglobals
// CharlottePublicKey secp256k1 key for testing purposes. Key is compromised do not use on mainnet's.
var CharlottePublicKey crypto.PublicKey //nolint: gochecknoglobals

//nolint: gochecknoinits
func init() {
	var err error
	encryptedSofiaPrivateKey, err := hex.DecodeString("01901E63389EF02EAA7C5782E08B40D98FAEF835F28BD144EECF5614A415943F")

	if err != nil {
		log.Fatal(err)
	}

	SofiaPrivateKey, err = secp256k1.PrivateKeyFromBytes(encryptedSofiaPrivateKey)
	if err != nil {
		log.Fatal(err)
	}

	SofiaPublicKey = SofiaPrivateKey.PublicKey()

	encryptedCharlottePrivateKey, err := hex.DecodeString("DF4BA9F6106AD2846472F759476535E55C5805D8337DF5A11C3B139F438B98B3")
	if err != nil {
		log.Fatal(err)
	}

	CharlottePrivateKey, err = secp256k1.PrivateKeyFromBytes(encryptedCharlottePrivateKey)

	if err != nil {
		log.Fatal(err)
	}

	CharlottePublicKey = CharlottePrivateKey.PublicKey()
}
