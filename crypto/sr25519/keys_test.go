package sr25519

import (
	"github.com/ChainSafe/go-schnorrkel"
	"github.com/mailchain/mailchain/internal/testutil"
)

var ( //nolint 
	sofiaSeed       = testutil.MustHexDecodeString("0x5c6d7adf75bda1180c225d25f3aa8dc174bbfb3cddee11ae9a85982f6faf791a") //nolint: lll
	sofiaPrivateKey = PrivateKey{key: func() *schnorrkel.SecretKey {
		priv, err := keyFromSeed([]byte{48, 120, 53, 99, 54, 100, 55, 97, 100, 102, 55, 53, 98, 100, 97, 49, 49, 56, 48, 99, 50, 50, 53, 100, 50, 53, 102, 51, 97, 97, 56})
		if err != nil {
			panic(err)
		}
		return priv
	}(),
	}
	sofiaPrivateKeyBytes = [32]byte{48, 120, 53, 99, 54, 100, 55, 97, 100, 102, 55, 53, 98, 100, 97, 49, 49, 56, 48, 99, 50, 50, 53, 100, 50, 53, 102, 51, 97, 97, 56}                 //nolint: lll
	sofiaPublicKey       = PublicKey{key: [32]byte{48, 120, 49, 54, 57, 97, 49, 49, 55, 50, 49, 56, 53, 49, 102, 53, 100, 102, 102, 51, 53, 52, 49, 100, 100, 53, 99, 52, 98, 48, 98}} //nolint: lll
	sofiaPublicKeyBytes  = [32]byte{48, 120, 49, 54, 57, 97, 49, 49, 55, 50, 49, 56, 53, 49, 102, 53, 100, 102, 102, 51, 53, 52, 49, 100, 100, 53, 99, 52, 98, 48, 98}                 //nolint: lll

	charlotteSeed            = testutil.MustHexDecodeString("0x23b063a581fd8e5e847c4e2b9c494247298791530f5293be369e8bf23a45d2bd")                                                                             //nolint: lll
	charlottePrivateKey      = PrivateKey{&schnorrkel.PrivateKey{[32]byte{48, 120, 50, 51, 98, 48, 54, 51, 97, 53, 56, 49, 102, 100, 56, 101, 53, 101, 56, 52, 55, 99, 52, 101, 50, 98, 57, 99, 52, 57, 52}}} //nolint: lll
	charlottePrivateKeyBytes = [32]byte{48, 120, 50, 51, 98, 48, 54, 51, 97, 53, 56, 49, 102, 100, 56, 101, 53, 101, 56, 52, 55, 99, 52, 101, 50, 98, 57, 99, 52, 57, 52}                                     //nolint: lll
	charlottePublicKey       = PublicKey{&schnorrkel.PublicKey{[32]byte{48, 120, 56, 52, 54, 50, 51, 101, 55, 50, 53, 50, 101, 52, 49, 49, 51, 56, 97, 102, 54, 57, 48, 52, 101, 49, 98, 48, 50, 51, 48}}}    //nolint: lll
	charlottePublicKeyBytes  = [32]byte{48, 120, 56, 52, 54, 50, 51, 101, 55, 50, 53, 50, 101, 52, 49, 49, 51, 56, 97, 102, 54, 57, 48, 52, 101, 49, 98, 48, 50, 51, 48}                                      //nolint: lll
) //
