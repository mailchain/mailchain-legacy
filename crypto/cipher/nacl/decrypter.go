package nacl

import (
	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/cipher"
	"github.com/pkg/errors"
)

// NewDecrypter create a new decrypter attaching the private key to it
func NewDecrypter(privateKey crypto.PrivateKey) (*Decrypter, error) {
	keyExchange, err := getPrivateKeyExchange(privateKey)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &Decrypter{privateKey: privateKey, keyExchange: keyExchange}, nil
}

// Decrypter will decrypt data using NACL with ECDH key exchange
type Decrypter struct {
	privateKey  crypto.PrivateKey
	keyExchange cipher.KeyExchange
}

// Decrypt data using recipient private key with AES in CBC mode.
func (d Decrypter) Decrypt(data cipher.EncryptedContent) (cipher.PlainContent, error) {
	data, pubKey, err := deserializeSecret(data)
	if err != nil {
		return nil, err
	}

	sharedSecret, err := d.keyExchange.SharedSecret(d.privateKey, pubKey)
	if err != nil {
		return nil, err
	}

	return easyOpen(data, sharedSecret)
}
