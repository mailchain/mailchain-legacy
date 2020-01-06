package nacl

import (
	"crypto/rand"
	"io"

	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/cipher"
	"github.com/mailchain/mailchain/crypto/ed25519"
	"github.com/mailchain/mailchain/crypto/sr25519"
	"github.com/pkg/errors"
)

// NewEncrypter creates a new encrypter with crypto rand for reader,
// and attaching the public key to the encrypter.
func NewEncrypter(publicKey crypto.PublicKey) (*Encrypter, error) {
	if err := validatePublicKeyType(publicKey); err != nil {
		return nil, errors.WithStack(err)
	}

	return &Encrypter{rand: rand.Reader, publicKey: publicKey}, nil
}

// Encrypter will encrypt data using AES256CBC method.
type Encrypter struct {
	rand      io.Reader
	publicKey crypto.PublicKey
}

// Encrypt encrypts the message with the key that was attached to it.
func (e Encrypter) Encrypt(message cipher.PlainContent) (cipher.EncryptedContent, error) {
	encrypted, err := easySeal(message, e.publicKey.Bytes(), e.rand)
	return bytesEncode(encrypted), err
}

func validatePublicKeyType(recipientPublicKey crypto.PublicKey) error {
	switch recipientPublicKey.(type) {
	case ed25519.PublicKey, *ed25519.PublicKey:
		return nil
	case sr25519.PublicKey:
		return nil
	default:
		return errors.Errorf("invalid public key type for nacl encryption")
	}
}
