package nacl

import (
	"crypto/rand"
	"io"

	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/cipher"
	"github.com/mailchain/mailchain/crypto/ed25519"
	"github.com/pkg/errors"
)

// NewEncrypter create a new encrypter with crypto rand for reader
func NewEncrypter() Encrypter {
	return Encrypter{rand: rand.Reader}
}

// Encrypter will encrypt data using AES256CBC method
type Encrypter struct {
	rand io.Reader
}

// Encrypt contents with the recipients public key.
func (e Encrypter) Encrypt(recipientPublicKey crypto.PublicKey, message cipher.PlainContent) (cipher.EncryptedContent, error) {
	if err := validatePublicKeyType(recipientPublicKey); err != nil {
		return nil, err
	}

	encrypted, err := easySeal(message, recipientPublicKey.Bytes(), e.rand)

	return bytesEncode(encrypted), err
}

func validatePublicKeyType(recipientPublicKey crypto.PublicKey) error {
	switch recipientPublicKey.(type) {
	case ed25519.PublicKey:
		return nil
	default:
		return errors.Errorf("invalid public key type for nacl encryption")
	}
}
