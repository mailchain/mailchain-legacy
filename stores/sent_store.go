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

//go:generate mockgen -source=sent.go -package=storestest -destination=./storestest/sent_mock.go

package stores

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/errs"
	"github.com/mailchain/mailchain/internal/mail"
	"github.com/pkg/errors"
)

func NewSentStore() *SentStore {
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	return &SentStore{
		domain:     "https://mcx.mx",
		newRequest: http.NewRequest,
		doRequest:  client.Do,
	}
}

type SentStore struct {
	domain     string
	newRequest func(method string, url string, body io.Reader) (*http.Request, error)
	doRequest  func(req *http.Request) (*http.Response, error)
}

func (s SentStore) Key(messageID mail.ID, msg []byte) string {
	hash := crypto.CreateLocationHash(msg)
	return fmt.Sprintf("%s-%s", messageID.HexString(), hash.HexString())
}

func (s SentStore) PutMessage(messageID mail.ID, msg []byte, headers map[string]string) (string, error) {
	hash := crypto.CreateLocationHash(msg)
	url := fmt.Sprintf("%s?hash=%s&message-id=%s", s.domain, hash.HexString(), messageID.HexString())

	req, err := s.newRequest("POST", url, bytes.NewReader(msg))
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/octet-stream")

	resp, err := s.doRequest(req)
	if err != nil {
		return "", err
	}

	if err := responseAsError(resp); err != nil {
		return "", err
	}

	loc := resp.Header.Get("Location")
	if loc == "" {
		return "", errors.Errorf("missing `Location` header")
	}

	return loc, nil
}

func responseAsError(r *http.Response) error {
	var httpError errs.HTTPError
	if r.StatusCode != http.StatusCreated {
		if err := json.NewDecoder(r.Body).Decode(&httpError); err != nil {
			return errors.WithMessage(err, "failed to read response")
		}
		return errors.Errorf("%v: %s", httpError.Code, httpError.Message)
	}
	return nil
}
