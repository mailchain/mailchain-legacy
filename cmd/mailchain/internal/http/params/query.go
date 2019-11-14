package params

import (
	"net/http"

	"github.com/pkg/errors"
)

// QueryRequireProtocol verify presence and return value of `protocol` url query parameter
func QueryRequireProtocol(r *http.Request) (string, error) {
	protocols := r.URL.Query()["protocol"]
	if len(protocols) != 1 {
		return "", errors.Errorf("'protocol' must be specified exactly once")
	}
	if protocols[0] == "" {
		return "", errors.Errorf("'protocol' must not be empty")
	}
	return protocols[0], nil
}

// QueryRequireNetwork verify presence and return value of `network` url query parameter
func QueryRequireNetwork(r *http.Request) (string, error) {
	networks := r.URL.Query()["network"]
	if len(networks) != 1 {
		return "", errors.Errorf("'network' must be specified exactly once")
	}
	if networks[0] == "" {
		return "", errors.Errorf("'network' must not be empty")
	}
	return networks[0], nil
}

// QueryRequireAddress verify presence and return value of `address` url query parameter
func QueryRequireAddress(r *http.Request) (string, error) {
	addresses := r.URL.Query()["address"]
	if len(addresses) != 1 {
		return "", errors.Errorf("'address' must be specified exactly once")
	}
	if addresses[0] == "" {
		return "", errors.Errorf("'address' must not be empty")
	}
	return addresses[0], nil
}
