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
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/errs"
	"github.com/mailchain/mailchain/internal/envelope"
	"github.com/mailchain/mailchain/internal/mail"
	"github.com/pkg/errors"
)

// NewSentStore create Mailchain sent store.
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

// SentStore type for storing sent Mailchain messages.
type SentStore struct {
	domain     string
	newRequest func(method string, url string, body io.Reader) (*http.Request, error)
	doRequest  func(req *http.Request) (*http.Response, error)
}

// Key gets the key of a Mailchain message.
func (s SentStore) Key(messageID mail.ID, contentsHash, msg []byte) string {
	return hex.EncodeToString(contentsHash)
}

// PutMessage stores message contents.
func (s SentStore) PutMessage(messageID mail.ID, contentsHash, msg []byte, headers map[string]string) (
	address, resource string, mli uint64, err error) {
	hash := crypto.CreateIntegrityHash(msg)
	url := fmt.Sprintf("%s?hash=%s&contents-hash=%s", s.domain, hash.HexString(), hex.EncodeToString(contentsHash))

	req, err := s.newRequest("POST", url, bytes.NewReader(msg))
	if err != nil {
		return "", "", envelope.MLIMailchain, err
	}

	req.Header.Add("Content-Type", "application/octet-stream")

	resp, err := s.doRequest(req)
	if err != nil {
		return "", "", envelope.MLIMailchain, err
	}

	if err := responseAsError(resp); err != nil {
		return "", "", envelope.MLIMailchain, err
	}

	loc := resp.Header.Get("Location")
	if loc == "" {
		return "", "", envelope.MLIMailchain, errors.Errorf("missing `Location` header")
	}

	mli, err = strconv.ParseUint(resp.Header.Get("Message-Location-Identifier"), 10, 0)
	if err != nil {
		return "", "", envelope.MLIMailchain, errors.Errorf("%q is not valid for `Message-Location-Identifier` header must be %v",
			resp.Header.Get("Message-Location-Identifier"), envelope.MLIMailchain)
	}
	if mli != envelope.MLIMailchain {
		return "", "", envelope.MLIMailchain, errors.Errorf("mismatch `Message-Location-Identifier` header")
	}

	return loc, hex.EncodeToString(contentsHash), envelope.MLIMailchain, nil
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
