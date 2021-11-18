package sr25519test

import (
	"log"

	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/sr25519"
	"github.com/mailchain/mailchain/encoding/encodingtest"
)

// AlicePrivateKey sr25519 key for testing purposes. Key is compromised do not use on mainnet's.
var AlicePrivateKey crypto.PrivateKey //nolint: gochecknoglobals test key
// AlicePublicKey sr25519 key for testing purposes. Key is compromised do not use on mainnet's.
var AlicePublicKey crypto.PublicKey //nolint: gochecknoglobals test key
// BobPrivateKey sr25519 key for testing purposes. Key is compromised do not use on mainnet's.
var BobPrivateKey crypto.PrivateKey //nolint: gochecknoglobals test key
// BobPublicKey sr25519 key for testing purposes. Key is compromised do not use on mainnet's.
var BobPublicKey crypto.PublicKey //nolint: gochecknoglobals test key
// EvePrivateKey sr25519 key for testing purposes. Key is compromised do not use on mainnet's.
var EvePrivateKey crypto.PrivateKey //nolint: gochecknoglobals test key
// EvePublicKey sr25519 key for testing purposes. Key is compromised do not use on mainnet's.
var EvePublicKey crypto.PublicKey //nolint: gochecknoglobals test key

//nolint: gochecknoinits test key
func init() {
	var err error
	AlicePrivateKey, err = sr25519.PrivateKeyFromBytes(encodingtest.MustDecodeHex("5c6d7adf75bda1180c225d25f3aa8dc174bbfb3cddee11ae9a85982f6faf791a")) //nolint: lll test key
	if err != nil {
		log.Fatal(err)
	}

	AlicePublicKey = AlicePrivateKey.PublicKey()

	BobPrivateKey, err = sr25519.PrivateKeyFromBytes(encodingtest.MustDecodeHex("23b063a581fd8e5e847c4e2b9c494247298791530f5293be369e8bf23a45d2bd")) //nolint: lll test key
	if err != nil {
		log.Fatal(err)
	}

	BobPublicKey = BobPrivateKey.PublicKey()

	EvePrivateKey, err = sr25519.PrivateKeyFromBytes(encodingtest.MustDecodeHex("000102030405060708090a0b0c0d0e0f000102030405060708090a0b0c0d0e0f")) //nolint: lll test key
	if err != nil {
		log.Fatal(err)
	}

	EvePublicKey = EvePrivateKey.PublicKey()
}
