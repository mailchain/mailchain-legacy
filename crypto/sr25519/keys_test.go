package sr25519

import (
	"github.com/ChainSafe/go-schnorrkel"
	"github.com/mailchain/mailchain/internal/testutil"
)

var ( //nolint
	sofiaSeed       = testutil.MustHexDecodeString("5c6d7adf75bda1180c225d25f3aa8dc174bbfb3cddee11ae9a85982f6faf791a") //nolint: lll
	sofiaPrivateKey = PrivateKey{key: func() *schnorrkel.SecretKey {
		priv, err := keyFromSeed([]byte{92, 109, 122, 223, 117, 189, 161, 24, 12, 34, 93, 37, 243, 170, 141, 193, 116, 187, 251, 60, 221, 238, 17, 174, 154, 133, 152, 47, 111, 175, 121, 26}) //nolint: lll
		if err != nil {
			panic(err)
		}

		return priv
	}(),
	}
	sofiaPrivateKeyBytes = []byte{92, 109, 122, 223, 117, 189, 161, 24, 12, 34, 93, 37, 243, 170, 141, 193, 116, 187, 251, 60, 221, 238, 17, 174, 154, 133, 152, 47, 111, 175, 121, 26} //nolint: lll
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
		priv, err := keyFromSeed([]byte{35, 176, 99, 165, 129, 253, 142, 94, 132, 124, 78, 43, 156, 73, 66, 71, 41, 135, 145, 83, 15, 82, 147, 190, 54, 158, 139, 242, 58, 69, 210, 189})
		if err != nil {
			panic(err)
		}

		return priv
	}()}
	charlottePrivateKeyBytes = []byte{35, 176, 99, 165, 129, 253, 142, 94, 132, 124, 78, 43, 156, 73, 66, 71, 41, 135, 145, 83, 15, 82, 147, 190, 54, 158, 139, 242, 58, 69, 210, 189}
	charlottePublicKey       = PublicKey{key: func() *schnorrkel.PublicKey {
		pubKey, err := schnorrkelPublicKeyFromBytes([]byte{132, 98, 62, 114, 82, 228, 17, 56, 175, 105, 4, 225, 176, 35, 4, 201, 65, 98, 95, 57, 229, 118, 37, 137, 18, 93, 193, 162, 242, 207, 46, 48})
		if err != nil {
			panic(err)
		}
		return pubKey

	}()}
	charlottePublicKeyBytes = []byte{132, 98, 62, 114, 82, 228, 17, 56, 175, 105, 4, 225, 176, 35, 4, 201, 65, 98, 95, 57, 229, 118, 37, 137, 18, 93, 193, 162, 242, 207, 46, 48}
) //nolint: lll
