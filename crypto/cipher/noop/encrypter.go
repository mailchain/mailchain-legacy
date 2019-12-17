package noop

import (
	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/cipher"
)

// NewEncrypter create a new encrypter with crypto rand for reader
func NewEncrypter(pubKey crypto.PublicKey) (*Encrypter, error) {
	return &Encrypter{publicKey: pubKey}, nil
}

// Encrypter will not perform any operation when encrypting the message.
//
// No operation (noop) encrypter is used when the contents of the message
// and envelope are intended to readable by the public.
type Encrypter struct {
	publicKey crypto.PublicKey
}

// Encrypt does not apply any encrption algortim.
// PlainContent will be return as EncryptedContent with the encryption method
// prepend as the first byte.
func (e Encrypter) Encrypt(message cipher.PlainContent) (cipher.EncryptedContent, error) {
	return bytesEncode(cipher.EncryptedContent(message)), nil
}
