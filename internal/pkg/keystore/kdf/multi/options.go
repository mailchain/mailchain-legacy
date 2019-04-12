package multi

import "github.com/mailchain/mailchain/internal/pkg/keystore/kdf/scrypt"

// OptionsBuilders contains all the builders for different key derivation functions
type OptionsBuilders struct {
	Scrypt []scrypt.DeriveOptionsBuilder
}
