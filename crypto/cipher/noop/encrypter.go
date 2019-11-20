package noop

import (
	"crypto"

	"github.com/mailchain/mailchain/crypto/cipher"
)

// NewEncrypter create a new encrypter with crypto rand for reader
func NewEncrypter() Encrypter {
	return Encrypter{}
}

// Encrypter will not perform any operation when encrypting the message.
//
// No operation (noop) encrypter is used when the contents of the message
// and envelope are intended to readable by the public.
type Encrypter struct {
}

// Encrypt does not apply any encrption algortim.
// PlainContent will be return as EncryptedContent with the encryption method
// prepend as the first byte.
func (e Encrypter) Encrypt(recipientPublicKey crypto.PublicKey, message cipher.PlainContent) (cipher.EncryptedContent, error) {
	return bytesEncode(cipher.EncryptedContent(message)), nil
}
