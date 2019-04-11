package keys

// PrivateKey definition usable in all mailchain crypto operations
type PrivateKey interface {
	// Bytes returns the byte representation of the private key
	Bytes() []byte
	// PublicKey from the PrivateKey
	PublicKey() PublicKey
}
