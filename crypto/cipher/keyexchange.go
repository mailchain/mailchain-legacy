package cipher

import "github.com/mailchain/mailchain/crypto"

//go:generate mockgen -source=keyexchange.go -package=ciphertest -destination=./ciphertest/keyexchange_mock.go

// KeyExchange agrees on a symmetric keys by performing a key exchange using asymmetric keys.
type KeyExchange interface {
	// EphemeralKey generates a private/public key pair.
	EphemeralKey() (private crypto.PrivateKey, err error)

	// SharedSecret computes a secret value from a private / public key pair.
	// On sending a message the private key should be an ephemeralKey or generated private key,
	// the public key is the recipient public key.
	// On reading a message the private key is the recipient private key, the public key is the
	// ephemeralKey or generated public key.
	SharedSecret(privateKey crypto.PrivateKey, publicKey crypto.PublicKey) ([]byte, error)
}
