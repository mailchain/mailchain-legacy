package crypto

type KeyExchange interface {
	// EphemeralKey generates a private/public key pair.
	EphemeralKey() (private PrivateKey, public PublicKey, err error)

	// SharedSecret computes a secret value from ephemeralKey private key and recipientKey public key.
	SharedSecret(ephemeralKey PrivateKey, recipientKey PublicKey) ([]byte, error)
}
