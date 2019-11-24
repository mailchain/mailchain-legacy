package sr25519

import (
	"github.com/ChainSafe/go-schnorrkel"
	"github.com/mailchain/mailchain/internal/testutil"
)

var ( //nolint
	sofiaSeed       = testutil.MustHexDecodeString("0x5c6d7adf75bda1180c225d25f3aa8dc174bbfb3cddee11ae9a85982f6faf791a") //nolint: lll
	sofiaPrivateKey = PrivateKey{key: func() *schnorrkel.SecretKey {
		priv, err := keyFromSeed(testutil.MustHexDecodeString("0x5c6d7adf75bda1180c225d25f3aa8dc174bbfb3cddee11ae9a85982f6faf791a"))
		if err != nil {
			panic(err)
		}

		return priv
	}(),
	}
	sofiaPrivateKeyBytes = testutil.MustHexDecodeString("0x5c6d7adf75bda1180c225d25f3aa8dc174bbfb3cddee11ae9a85982f6faf791a") //nolint: lll
	sofiaPublicKey       = PublicKey{key: func() *schnorrkel.PublicKey {
		pubKey, err := schnorrkelPublicKeyFromBytes(testutil.MustHexDecodeString("0x169a11721851f5dff3541dd5c4b0b478ac1cd092c9d5976e83daa0d03f26620c")) //nolint: lll
		if err != nil {
			panic(err)
		}

		return pubKey

	}()}
	sofiaPublicKeyBytes = testutil.MustHexDecodeString("0x169a11721851f5dff3541dd5c4b0b478ac1cd092c9d5976e83daa0d03f26620c") //nolint: lll

	charlotteSeed       = testutil.MustHexDecodeString("0x23b063a581fd8e5e847c4e2b9c494247298791530f5293be369e8bf23a45d2bd") //nolint: lll
	charlottePrivateKey = PrivateKey{key: func() *schnorrkel.SecretKey {
		priv, err := keyFromSeed(testutil.MustHexDecodeString("0x23b063a581fd8e5e847c4e2b9c494247298791530f5293be369e8bf23a45d2bd")) //nolint: lll
		if err != nil {
			panic(err)
		}

		return priv
	}(),
	}
	charlottePrivateKeyBytes = testutil.MustHexDecodeString("0x23b063a581fd8e5e847c4e2b9c494247298791530f5293be369e8bf23a45d2bd") //nolint: lll
	charlottePublicKey       = PublicKey{key: func() *schnorrkel.PublicKey {
		pubKey, err := schnorrkelPublicKeyFromBytes(testutil.MustHexDecodeString("0x84623e7252e41138af6904e1b02304c941625f39e5762589125dc1a2f2cf2e30")) //nolint: lll
		if err != nil {
			panic(err)
		}
		return pubKey

	}()}
	charlottePublicKeyBytes = testutil.MustHexDecodeString("0x84623e7252e41138af6904e1b02304c941625f39e5762589125dc1a2f2cf2e30") //nolint: lll
) //nolint: lll
