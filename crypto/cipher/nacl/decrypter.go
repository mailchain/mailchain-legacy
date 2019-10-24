package nacl

import (
	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/cipher"
)

// NewDecrypter create a new decrypter attaching the private key to it
func NewDecrypter(privateKey crypto.PrivateKey) Decrypter {
	return Decrypter{privateKey: privateKey}
}

// Decrypter will decrypt data using AES256CBC method
type Decrypter struct {
	privateKey crypto.PrivateKey
}

// Decrypt data using recipient private key with AES in CBC mode.
func (d Decrypter) Decrypt(data cipher.EncryptedContent) (cipher.PlainContent, error) {
	return easyOpen(data, d.privateKey.Bytes())
}
