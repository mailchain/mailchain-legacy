package envelope

import (
	"net/url"

	"github.com/mailchain/mailchain/crypto/cipher"
)

const (
	Kind0x01 byte = 0x01 // Message locator
	Kind0x50 byte = 0x50 // Alpha
)

type Data interface {
	URL(decrypter cipher.Decrypter) (*url.URL, error)
	IntegrityHash(decrypter cipher.Decrypter) ([]byte, error)
	ContentsHash(decrypter cipher.Decrypter) ([]byte, error)
	Valid() error
}
