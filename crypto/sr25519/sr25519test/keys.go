package sr25519test

import (
	"crypto"
	"github.com/mailchain/mailchain/internal/encoding"
	"log"
)

// SofiaPrivateKey sr25519 key for testing purposes. Key is compromised do not use on mainnet's.
var SofiaPrivateKey crypto.PrivateKey //nolint: gochecknoglobals
// SofiaPublicKey sr25519 key for testing purposes. Key is compromised do not use on mainnet's.
var SofiaPublicKey crypto.PublicKey //nolint: gochecknoglobals
// CharlottePrivateKey sr25519 key for testing purposes. Key is compromised do not use on mainnet's.
var CharlottePrivateKey crypto.PrivateKey //nolint: gochecknoglobals
// CharlottePublicKey sr25519 key for testing purposes. Key is compromised do not use on mainnet's.
var CharlottePublicKey crypto.PublicKey //nolint: gochecknoglobals

func int() {
	var err error

	sofiaByte := [32]byte{}
	encodedSofia, err := encoding.DecodeZeroX("0x5c6d7adf75bda1180c225d25f3aa8dc174bbfb3cddee11ae9a85982f6faf791a")
	if err != nil {
		log.Fatal(err)
	}

	copy(sofiaByte[:], encodedSofia)

	SofiaPrivateKey, err = PrivateKeyFromBytes(encodedSofia)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%v'\n", SofiaPrivateKey)

	SofiaPrivateKey = SofiaPrivateKey.PublicKey()

	charlotteByte := [32]byte{}
	encodedCharlotte, err := encoding.DecodeZeroX("0x23b063a581fd8e5e847c4e2b9c494247298791530f5293be369e8bf23a45d2bd")
	if err != nil {
		log.Fatal(err)
	}

	copy(charlotteByte[:], encodedCharlotte)
	CharlottePrivateKey, err = PrivateKeyFromBytes(encodedCharlotte)
	if err != nil {
		log.Fatal(err)
	}

	CharlottePublicKey = CharlottePrivateKey.PublicKey()

}
