package envelope

import (
	"net/url"

	"github.com/mailchain/mailchain/crypto/cipher"
)

type Data interface {
	URL(decrypter cipher.Decrypter) (*url.URL, error)
	IntegrityHash(decrypter cipher.Decrypter) ([]byte, error)
	ContentsHash(decrypter cipher.Decrypter) ([]byte, error)
	Valid() error
}
