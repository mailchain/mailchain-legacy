package encrypter

import (
	crypto "github.com/mailchain/mailchain/crypto/cipher"
	"github.com/mailchain/mailchain/crypto/cipher/aes256cbc"
	"github.com/mailchain/mailchain/crypto/cipher/nacl"
	"github.com/pkg/errors"
)

// Cipher Name lookup
const (
	NoOperation string = "noop"
	NACL        string = "nacl"
	AES256CBC   string = "aes256cbc"
)

// GetEncrypter is an `Encrypter` factory that returns an encrypter
func GetEncrypter(encryption string) (crypto.Encrypter, error) {
	switch encryption {
	case AES256CBC:
		return aes256cbc.NewEncrypter(), nil
	case NACL:
		return nacl.NewEncrypter(), nil
	case "":
		return nil, errors.Errorf("`encryption` provided is set to empty")
	default:
		return nil, errors.Errorf("`encryption` provided is invalid")
	}
}
