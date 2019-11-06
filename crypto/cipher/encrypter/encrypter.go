package encrypter

import (
	crypto "github.com/mailchain/mailchain/crypto/cipher"
	"github.com/mailchain/mailchain/crypto/cipher/aes256cbc"
	"github.com/mailchain/mailchain/crypto/cipher/nacl"
	"github.com/pkg/errors"
)

// GetEncrypter is an `Encrypter` factory that returns an encrypter
func GetEncrypter(encryption byte) (crypto.Encrypter, error) {
	switch encryption {
	case crypto.AES256CBC:
		return aes256cbc.NewEncrypter(), nil
	case crypto.NACL:
		return nacl.NewEncrypter(), nil
	default:
		return nil, errors.Errorf("string provided is invalid")
	}
}
