package nacl

import (
	"crypto/rand"
	"io"

	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/cipher"
	"github.com/pkg/errors"
)

// NewEncrypter creates a new encrypter with crypto rand for reader,
// and attaching the public key to the encrypter.
func NewEncrypter(publicKey crypto.PublicKey) (*Encrypter, error) {
	keyExchange, err := getPublicKeyExchange(publicKey)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &Encrypter{rand: rand.Reader, publicKey: publicKey, keyExchange: keyExchange}, nil
}

// Encrypter will encrypt data using AES256CBC method.
type Encrypter struct {
	rand        io.Reader
	publicKey   crypto.PublicKey
	keyExchange cipher.KeyExchange
}

// Encrypt encrypts the message with the key that was attached to it.
func (e Encrypter) Encrypt(message cipher.PlainContent) (cipher.EncryptedContent, error) {
	ephemeralKey, err := e.keyExchange.EphemeralKey()
	if err != nil {
		return nil, err
	}

	sharedSecret, err := e.keyExchange.SharedSecret(ephemeralKey, e.publicKey)
	if err != nil {
		return nil, err
	}

	encrypted, err := easySeal(message, sharedSecret, e.rand)
	if err != nil {
		return nil, err
	}

	return bytesEncode(encrypted, ephemeralKey.PublicKey())
}
