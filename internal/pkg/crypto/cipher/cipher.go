//go:generate mockgen -source=crypto.go -package=mocks -destination=$PACKAGE_PATH/internal/pkg/testutil/mocks/cipher.go
package cipher

import (
	"io"

	"github.com/mailchain/mailchain/internal/pkg/crypto/keys"
)

// EncryptedContent typed version of byte array that holds encrypted data
type EncryptedContent []byte

// PlainContent typed version of byte array that holds plain data
type PlainContent []byte

// Decrypter will decrypt data using specified method
type Decrypter interface {
	Decrypt(EncryptedContent) (PlainContent, error)
}

// Encrypter will encrypt data using public key
type Encrypter interface {
	Encrypt(rand io.Reader, pub keys.PublicKey, plain PlainContent) (EncryptedContent, error)
}
