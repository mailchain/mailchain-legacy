package sr25519

import (
	"github.com/ChainSafe/go-schnorrkel"
	"github.com/mailchain/mailchain/internal/testutil"
)

var ( //nolint

	sofiaSeed       = testutil.MustHexDecodeString("5c6d7adf75bda1180c225d25f3aa8dc174bbfb3cddee11ae9a85982f6faf791a") //nolint: lll
	sofiaPrivateKey = PrivateKey{key: func() *schnorrkel.SecretKey {
		priv, err := keyFromSeed([]byte{0x5c, 0x6d, 0x7a, 0xdf, 0x75, 0xbd, 0xa1, 0x18, 0xc, 0x22, 0x5d, 0x25, 0xf3, 0xaa, 0x8d, 0xc1, 0x74, 0xbb, 0xfb, 0x3c, 0xdd, 0xee, 0x11, 0xae, 0x9a, 0x85, 0x98, 0x2f, 0x6f, 0xaf, 0x79, 0x1a})
		if err != nil {
			panic(err)
		}

		return priv
	}(),
	}

	sofiaPrivateKeyBytes = []byte{0x5c, 0x6d, 0x7a, 0xdf, 0x75, 0xbd, 0xa1, 0x18, 0xc, 0x22, 0x5d, 0x25, 0xf3, 0xaa, 0x8d, 0xc1, 0x74, 0xbb, 0xfb, 0x3c, 0xdd, 0xee, 0x11, 0xae, 0x9a, 0x85, 0x98, 0x2f, 0x6f, 0xaf, 0x79, 0x1a}
	sofiaPublicKey       = PublicKey{key: func() *schnorrkel.PublicKey {
		pubKey, err := schnorrkelPublicKeyFromBytes([]byte{22, 154, 17, 114, 24, 81, 245, 223, 243, 84, 29, 213, 196, 176, 180, 120, 172, 28, 208, 146, 201, 213, 151, 110, 131, 218, 160, 208, 63, 38, 98, 12}) //nolint: lll
		if err != nil {
			panic(err)
		}

		return pubKey
	}()}
	sofiaPublicKeyBytes = []byte{22, 154, 17, 114, 24, 81, 245, 223, 243, 84, 29, 213, 196, 176, 180, 120, 172, 28, 208, 146, 201, 213, 151, 110, 131, 218, 160, 208, 63, 38, 98, 12} //nolint: lll

	charlotteSeed       = testutil.MustHexDecodeString("23b063a581fd8e5e847c4e2b9c494247298791530f5293be369e8bf23a45d2bd") //nolint: lll
	charlottePrivateKey = PrivateKey{key: func() *schnorrkel.SecretKey {
		priv, err := keyFromSeed([]byte{0x23, 0xb0, 0x63, 0xa5, 0x81, 0xfd, 0x8e, 0x5e, 0x84, 0x7c, 0x4e, 0x2b, 0x9c, 0x49, 0x42, 0x47, 0x29, 0x87, 0x91, 0x53, 0xf, 0x52, 0x93, 0xbe, 0x36, 0x9e, 0x8b, 0xf2, 0x3a, 0x45, 0xd2, 0xbd})
		if err != nil {
			panic(err)
		}

		return priv
	}()}
	charlottePrivateKeyBytes = []byte{0x23, 0xb0, 0x63, 0xa5, 0x81, 0xfd, 0x8e, 0x5e, 0x84, 0x7c, 0x4e, 0x2b, 0x9c, 0x49, 0x42, 0x47, 0x29, 0x87, 0x91, 0x53, 0xf, 0x52, 0x93, 0xbe, 0x36, 0x9e, 0x8b, 0xf2, 0x3a, 0x45, 0xd2, 0xbd}
	charlottePublicKey       = PublicKey{key: func() *schnorrkel.PublicKey {
		pubKey, err := schnorrkelPublicKeyFromBytes([]byte{132, 98, 62, 114, 82, 228, 17, 56, 175, 105, 4, 225, 176, 35, 4, 201, 65, 98, 95, 57, 229, 118, 37, 137, 18, 93, 193, 162, 242, 207, 46, 48})
		if err != nil {
			panic(err)
		}
		return pubKey

	}()}
	charlottePublicKeyBytes = []byte{132, 98, 62, 114, 82, 228, 17, 56, 175, 105, 4, 225, 176, 35, 4, 201, 65, 98, 95, 57, 229, 118, 37, 137, 18, 93, 193, 162, 242, 207, 46, 48}
	charlotteSign, _        = charlottePrivateKey.Sign([]byte("message"))
) //nolint: lll
