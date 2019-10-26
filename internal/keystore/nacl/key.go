package nacl

import (
	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/multikey"
	"github.com/mailchain/mailchain/internal/keystore"
	"github.com/mailchain/mailchain/internal/keystore/kdf"
	"github.com/mailchain/mailchain/internal/keystore/kdf/multi"
	"github.com/mailchain/mailchain/internal/keystore/kdf/scrypt"
	"github.com/pkg/errors"
)

func deriveKey(ek *keystore.EncryptedKey, deriveKeyOptions multi.OptionsBuilders) ([]byte, error) {
	switch ek.KDF {
	case kdf.Scrypt:
		if ek.ScryptParams == nil {
			return nil, errors.New("scryptParams are required")
		}
		storageOpts := scrypt.FromEncryptedKey(ek.ScryptParams.Len, ek.ScryptParams.N, ek.ScryptParams.P, ek.ScryptParams.R, ek.ScryptParams.Salt)

		return scrypt.DeriveKey(append(deriveKeyOptions.Scrypt, storageOpts))
	default:
		return nil, errors.New("KDF is not supported")
	}
}

func (f FileStore) getPrivateKey(address []byte, deriveKeyOptions multi.OptionsBuilders) (crypto.PrivateKey, error) {
	encryptedKey, err := f.getEncryptedKey(address)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	storageKey, err := deriveKey(encryptedKey, deriveKeyOptions)
	if err != nil {
		return nil, errors.WithMessage(err, "storage key could not be derived")
	}
	pkBytes, err := easyOpen(encryptedKey.CipherText, storageKey)
	if err != nil {
		return nil, errors.WithMessage(err, "could not decrypt key file")
	}
	pk, err := multikey.PrivateKeyFromBytes(encryptedKey.CurveType, pkBytes)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return pk, nil
}
