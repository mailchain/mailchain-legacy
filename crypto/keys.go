package crypto

const (
	// KindSECP256K1 string identifier for secp256k1 keys.
	KindSECP256K1 = "secp256k1"
	// KindED25519 string identifier for ed25519 keys.
	KindED25519 = "ed25519"
	// KindSR25519 string identifier for sr25519 keys.
	KindSR25519 = "sr25519"
)

const (
	// ByteSECP256K1 byte identifier for secp256k1 keys.
	ByteSECP256K1 = 0xe1
	// ByteED25519 byte identifier for ed25519 keys.
	ByteED25519 = 0xe2
	// ByteSR25519 byte identifier for sr25519 keys.
	ByteSR25519 = 0xe3
)

// KeyTypes available key types.
func KeyTypes() map[string]bool {
	return map[string]bool{
		KindSECP256K1: true,
		KindED25519:   true,
		KindSR25519:   true,
	}
}
