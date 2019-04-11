package secp256k1

import "crypto/ecdsa"

// PrivateKey based on the secp256k1 curve
type PrivateKey struct {
	ecdsa ecdsa.PrivateKey
}
