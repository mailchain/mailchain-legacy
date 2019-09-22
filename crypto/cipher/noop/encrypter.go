package noop

import (
	"crypto"

	"github.com/mailchain/mailchain/crypto/cipher"
)

// NewEncrypter create a new encrypter with crypto rand for reader
func NewEncrypter() Encrypter {
	return Encrypter{}
}

// Encrypter will encrypt data using AES256CBC method
type Encrypter struct {
}

// Encrypt noop (no operation) encrypter returns the plain content
func (e Encrypter) Encrypt(recipientPublicKey crypto.PublicKey, message cipher.PlainContent) (cipher.EncryptedContent, error) {
	return cipher.EncryptedContent(message), nil
}
