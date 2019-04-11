package keys

// PublicKey definition usable in all mailchain crypto operations
type PublicKey interface {
	// Bytes returns the byte representation of the public key
	Bytes() []byte
	// Address returns the byte representation of the address
	Address() []byte
}
