package mailbox

import (
	"github.com/pkg/errors"
	"gopkg.in/resty.v1"
)

func getHTTPMessage(location string) ([]byte, error) {
	res, err := resty.R().Get(location)
	if err != nil {
		return nil, errors.Wrap(err, "could not get message from `location`")
	}
	msg := res.Body()

	return msg, nil
}
