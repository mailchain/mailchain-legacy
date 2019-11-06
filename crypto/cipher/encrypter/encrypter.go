package encrypter

import (
	"strings"

	crypto "github.com/mailchain/mailchain/crypto/cipher"
	"github.com/mailchain/mailchain/crypto/cipher/aes256cbc"
	"github.com/pkg/errors"
)

// GetEncrypter is an `Encrypter` factory that returns an encrypter
func GetEncrypter(encryption string) (crypto.Encrypter, error) {
	switch strings.ToLower(encryption) {
	case "aes256cbc":
		return aes256cbc.NewEncrypter(), nil
	default:
		return nil, errors.Errorf("string provided is invalid")
	}
}
