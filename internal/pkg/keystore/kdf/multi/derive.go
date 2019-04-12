package multi

import (
	"github.com/mailchain/mailchain/internal/pkg/keystore/kdf/scrypt"
	"github.com/pkg/errors"
)

func DeriveKey(options OptionsBuilders) (storageKey []byte, kdf string, err error) {
	if options.Scrypt != nil {
		storageKey, err = scrypt.DeriveKey(options.Scrypt)
		return storageKey, "scrypt", err
	}
	return nil, "unknown", errors.Errorf("unknown `kdf`")
}
