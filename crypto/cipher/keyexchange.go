package cipher

import (
	"github.com/mailchain/mailchain/crypto"
)

type KeyExchange interface {
	// EphemeralKey generates a private/public key pair.
	EphemeralKey() (private crypto.PrivateKey, public crypto.PublicKey, err error)

	// SharedSecret computes a secret value from ephemeralKey private key and recipientKey public key.
	SharedSecret(ephemeralKey crypto.PrivateKey, recipientKey crypto.PublicKey) ([]byte, error)
}
