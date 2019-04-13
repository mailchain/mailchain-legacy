package mailbox

import (
	"io/ioutil"
	"net/url"
	"strings"

	"github.com/mailchain/mailchain/internal/pkg/crypto"
	"github.com/mailchain/mailchain/internal/pkg/crypto/cipher"
	"github.com/mailchain/mailchain/internal/pkg/mail"
	"github.com/pkg/errors"
	"gopkg.in/resty.v1"
)

// decryptLocation return the location in readable form
func decryptLocation(d *mail.Data, decrypter cipher.Decrypter) (string, error) {
	decryptedLocation, err := decrypter.Decrypt(d.EncryptedLocation)
	if err != nil {
		return "", errors.WithMessage(err, "could not decrypt location")
	}
	return string(decryptedLocation), nil
}

// getMessage get the message contents from the location and perform location hash check
func getMessage(location string) ([]byte, error) {
	msg, err := getAnyMessage(location)
	if err != nil {
		return nil, err
	}

	hash, err := crypto.CreateLocationHash(msg)
	if err != nil {
		return nil, err
	}
	parts := strings.Split(location, "-")
	if len(parts) < 1 {
		return nil, errors.Errorf("could not safely extract hash from location")
	}
	if hash.String() != parts[len(parts)-1] {
		return nil, errors.Errorf("hash does not match contents")
	}
	return msg, nil
}

func getAnyMessage(location string) ([]byte, error) {
	parsed, err := url.Parse(location)
	if err != nil {
		return nil, err
	}

	switch parsed.Scheme {
	case "http":
		return getHTTPMessage(location)
	case "file":
		return ioutil.ReadFile(parsed.Host + parsed.Path)
	case "test":
		return []byte(parsed.Host), nil
	default:
		return nil, errors.Errorf("unsupported scheme")
	}
}
func getHTTPMessage(location string) ([]byte, error) {
	res, err := resty.R().Get(location)
	if err != nil {
		return nil, errors.Wrap(err, "could not get message from `location`")
	}
	msg := res.Body()

	return msg, nil
}
