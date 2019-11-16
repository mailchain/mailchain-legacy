package crypto

const (
	// SECP256K1 value.
	SECP256K1 = "secp256k1"
	// ED25519 value.
	ED25519 = "ed25519"
)

// KeyTypes available key types.
func KeyTypes() map[string]bool {
	return map[string]bool{
		SECP256K1: true,
		ED25519:   true,
	}
}
