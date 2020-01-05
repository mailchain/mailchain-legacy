package cipher

import "github.com/mailchain/mailchain/crypto"

// KeyExchange agrees on a symmetric keys by performing a key exchange using asymmetric keys.
type KeyExchange interface {
	// EphemeralKey generates a private/public key pair.
	EphemeralKey() (private crypto.PrivateKey, err error)

	// SharedSecret computes a secret value from ephemeralKey private key and recipientKey public key.
	SharedSecret(ephemeralKey crypto.PrivateKey, recipientKey crypto.PublicKey) ([]byte, error)
}
