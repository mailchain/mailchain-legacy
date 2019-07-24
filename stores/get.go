// Copyright 2019 Finobo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package stores

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/mailchain/mailchain/crypto"
	"github.com/pkg/errors"
	"gopkg.in/resty.v1"
)

// GetMessage get the message contents from the location and perform location hash check
func GetMessage(location string, integrityHash []byte) ([]byte, error) {
	msg, err := getAnyMessage(location)
	if err != nil {
		return nil, err
	}
	hash := crypto.CreateIntegrityHash(msg)
	if len(integrityHash) != 0 && !bytes.Equal(hash, integrityHash) {
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
	case "http", "https":
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
	if res.StatusCode() != http.StatusOK {
		return nil, errors.Errorf(res.Status())
	}
	return res.Body(), nil
}
