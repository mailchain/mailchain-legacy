package mailbox

import (
	"io/ioutil"
	"net/url"

	"github.com/pkg/errors"
	"gopkg.in/resty.v1"
)

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
