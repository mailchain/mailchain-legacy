package crypto

// KeyExchange agrees on a symmetric keys by performing a key exchange using asymmetric keys.
type KeyExchange interface {
	// EphemeralKey generates a private/public key pair.
	EphemeralKey() (private PrivateKey, err error)

	// SharedSecret computes a secret value from ephemeralKey private key and recipientKey public key.
	SharedSecret(ephemeralKey PrivateKey, recipientKey PublicKey) ([]byte, error)
}
