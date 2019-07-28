package nameservice

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

type Lookup interface {
	ForwardLookup
	ReverseLookup
}

func NewLookupService(baseURL string) Lookup {
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	return &LookupService{
		baseURL:    strings.TrimSuffix(baseURL, "/"),
		newRequest: http.NewRequest,
		doRequest:  client.Do,
	}
}

type LookupService struct {
	baseURL    string
	newRequest func(method string, url string, body io.Reader) (*http.Request, error)
	doRequest  func(req *http.Request) (*http.Response, error)
}

func (s LookupService) ResolveName(ctx context.Context, protocol, network, domainName string) ([]byte, error) {
	req, err := s.newRequest("GET", fmt.Sprintf("%s/%s/%s/name?domain-name=%s", s.baseURL, protocol, network, domainName), nil)
	if err != nil {
		return nil, err
	}
	res, err := s.doRequest(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode == http.StatusOK {
		type response struct {
			Address string `json:"address"`
		}
		var okRes response
		if err := json.NewDecoder(res.Body).Decode(&okRes); err != nil {
			return nil, err
		}

		return common.FromHex(okRes.Address), nil
	}

	type response struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
	}
	var errRes response
	if err := json.NewDecoder(res.Body).Decode(&errRes); err != nil {
		return nil, err
	}

	return nil, errors.Errorf(errRes.Message)
}

func (s LookupService) ResolveAddress(ctx context.Context, protocol, network string, address []byte) (string, error) {
	// {protocol}/{network}/address?address={address}
	req, err := s.newRequest("GET", fmt.Sprintf("%s/%s/%s/address?address=%s", s.baseURL, protocol, network, common.BytesToAddress(address).Hex()), nil)
	if err != nil {
		return "", err
	}
	res, err := s.doRequest(req)
	if err != nil {
		return "", err
	}
	if res.StatusCode == http.StatusOK {
		type response struct {
			Name string `json:"name"`
		}
		var okRes response
		if err := json.NewDecoder(res.Body).Decode(&okRes); err != nil {
			return "", err
		}

		return okRes.Name, nil
	}

	type response struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
	}
	var errRes response
	if err := json.NewDecoder(res.Body).Decode(&errRes); err != nil {
		return "", err
	}

	return "", errors.Errorf(errRes.Message)
}
