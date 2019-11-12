package sr25519test

import (
	"log"

	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/internal/testutil"
	"github.com/mailchain/mailchain/crypto/sr25519"
)

// SofiaPrivateKey ed25519 key for testing purposes. Key is compromised do not use on mainnet's.
var SofiaPrivateKey crypto.PrivateKey //nolint: gochecknoglobals
// SofiaPublicKey ed25519 key for testing purposes. Key is compromised do not use on mainnet's.
var SofiaPublicKey crypto.PublicKey //nolint: gochecknoglobals
// CharlottePrivateKey ed25519 key for testing purposes. Key is compromised do not use on mainnet's.
var CharlottePrivateKey crypto.PrivateKey //nolint: gochecknoglobals
// CharlottePublicKey ed25519 key for testing purposes. Key is compromised do not use on mainnet's.
var CharlottePublicKey crypto.PublicKey //nolint: gochecknoglobals

func int() {
	var err error
	msc, err := NewMiniSecretKeyFromRaw([32]byte{})
	if err != nil {
		t.Fatal(err)
	}

	sc := msc.ExpandEd25519()
	expected, err := crypto.NewMiniSecretKey(testutil.MustHexDecodeString("0d9b4a3c10721991c6b806f0f343535dc2b46c74bece50a0a0d6b9f0070d3157"))
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(sc.key[:], expected[:32]) {
		t.Errorf("Fail to expand key: got %x expected %x", sc.key, expected[:32])
	}

	if !bytes.Equal(sc.nonce[:], expected[32:64]) {
		t.Errorf("Fail to expand nonce: got %x expected %x", sc.nonce, expected[32:64])
	}

	pub := msc.Public().Compress()
	if !bytes.Equal(pub[:], expected[64:]) {
		t.Errorf("Fail to expand nonce: got %x expected %x", sc.nonce, expected[32:64])
	}

	SofiaPrivateKey, err = crypto.PrivateKeyFromBytes(testutil.MustHexDecodeString("0d9b4a3c10721991c6b806f0f343535dc2b46c74bece50a0a0d6b9f0070d3157"))
	if err != nil {
		log.Fatal(err)
	}

	SofiaPublicKey = SofiaPrivateKey.PublicKey()

	CharlottePrivateKey, err = crypto.PrivateKeyFromBytes(testutil.MustHexDecodeString("39d4c97d6a7f9e3306a2b5aae604ee67ec8b1387fffb39128fc055656cff05bb"))
	if err != nil {
		log.Fatal(err)
	}

	CharlottePublicKey = CharlottePrivateKey.PublicKey()

	
}