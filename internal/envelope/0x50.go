package envelope

import (
	"encoding/hex"
	"net/url"
	"strings"

	"github.com/mailchain/mailchain/crypto/cipher"
	"github.com/pkg/errors"
)

func (d *ZeroX50) URL(decrypter cipher.Decrypter) (*url.URL, error) {
	loc, err := decrypter.Decrypt(d.EncryptedURL)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return url.Parse(string(loc))
}

func (d *ZeroX50) ContentsHash(decrypter cipher.Decrypter) ([]byte, error) {
	return d.DecryptedHash, nil
}

func (d *ZeroX50) IntegrityHash(decrypter cipher.Decrypter) ([]byte, error) {
	loc, err := decrypter.Decrypt(d.EncryptedURL)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	parts := strings.Split(string(loc), "-")
	if len(parts) < 2 {
		return nil, errors.Errorf("could not safely extract hash from location")
	}
	return hex.DecodeString(parts[len(parts)-1])
}

func (d *ZeroX50) Valid() error {
	if len(d.EncryptedURL) == 0 {
		return errors.Errorf("EncryptedURL must not be empty")
	}

	if len(d.DecryptedHash) == 0 {
		return errors.Errorf("DecryptedHash must not be empty")
	}

	return nil
}
