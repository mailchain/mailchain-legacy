package secp256k1test

import (
	"log"

	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/secp256k1"
	"github.com/mailchain/mailchain/encoding/encodingtest"
)

// AlicePrivateKey secp256k1 key for testing purposes. Key is compromised do not use on mainnet's.
var AlicePrivateKey crypto.PrivateKey //nolint: gochecknoglobals
// AlicePublicKey secp256k1 key for testing purposes. Key is compromised do not use on mainnet's.
var AlicePublicKey crypto.PublicKey //nolint: gochecknoglobals
// BobPrivateKey secp256k1 key for testing purposes. Key is compromised do not use on mainnet's.
var BobPrivateKey crypto.PrivateKey //nolint: gochecknoglobals
// BobPublicKey secp256k1 key for testing purposes. Key is compromised do not use on mainnet's.
var BobPublicKey crypto.PublicKey //nolint: gochecknoglobals

//nolint: gochecknoinits
func init() {
	var err error

	AlicePrivateKey, err = secp256k1.PrivateKeyFromBytes(encodingtest.MustDecodeHex("01901E63389EF02EAA7C5782E08B40D98FAEF835F28BD144EECF5614A415943F"))
	if err != nil {
		log.Fatal(err)
	}

	AlicePublicKey = AlicePrivateKey.PublicKey()

	BobPrivateKey, err = secp256k1.PrivateKeyFromBytes(encodingtest.MustDecodeHex("DF4BA9F6106AD2846472F759476535E55C5805D8337DF5A11C3B139F438B98B3"))
	if err != nil {
		log.Fatal(err)
	}

	BobPublicKey = BobPrivateKey.PublicKey()
}
