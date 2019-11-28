package nacl

import (
	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/cipher"
	"github.com/mailchain/mailchain/crypto/ed25519"
	"github.com/pkg/errors"
)

// NewDecrypter create a new decrypter attaching the private key to it
func NewDecrypter(privateKey crypto.PrivateKey) (*Decrypter, error) {
	_, err := validatePrivateKeyType(privateKey)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &Decrypter{privateKey: privateKey}, nil
}

// Decrypter will decrypt data using AES256CBC method
type Decrypter struct {
	privateKey crypto.PrivateKey
}

// Decrypt data using recipient private key with AES in CBC mode.
func (d Decrypter) Decrypt(data cipher.EncryptedContent) (cipher.PlainContent, error) {
	privKeyBytes, err := validatePrivateKeyType(d.privateKey)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	data, err = bytesDecode(data)
	if err != nil {
		return nil, err
	}

	return easyOpen(data, privKeyBytes)
}

func validatePrivateKeyType(pk crypto.PrivateKey) ([]byte, error) {
	switch pk := pk.(type) {
	case ed25519.PrivateKey, *ed25519.PrivateKey:
		return pk.Bytes()[32:], nil
	default:
		return nil, errors.Errorf("invalid public key type for nacl encryption")
	}
}
