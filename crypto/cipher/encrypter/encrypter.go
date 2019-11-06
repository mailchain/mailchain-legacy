package encrypter

import (
	"strings"

	crypto "github.com/mailchain/mailchain/crypto/cipher"
	"github.com/mailchain/mailchain/crypto/cipher/aes256cbc"
	"github.com/mailchain/mailchain/crypto/cipher/nacl"
	"github.com/pkg/errors"
)

// GetEncrypter is an `Encrypter` factory that returns an encrypter
func GetEncrypter(encryption string) (crypto.Encrypter, error) {
	method := strings.ToLower(encryption)
	if method == "aes256cbc" {
		return aes256cbc.NewEncrypter(), nil
	}
	if method == "nalc" {
		return nacl.NewEncrypter(), nil
	}
	return nil, errors.Errorf("string provided is invalid")
}
