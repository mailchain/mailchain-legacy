package params

import (
	"net/http"
	"strconv"

	"github.com/pkg/errors"
)

// QueryRequireProtocol verify presence and return value of `protocol` url query parameter
func QueryRequireProtocol(r *http.Request) (string, error) {
	protocols := r.URL.Query()["protocol"]
	if len(protocols) != 1 {
		return "", errors.New("'protocol' must be specified exactly once")
	}

	if protocols[0] == "" {
		return "", errors.New("'protocol' must not be empty")
	}

	return protocols[0], nil
}

// QueryRequireNetwork verify presence and return value of `network` url query parameter
func QueryRequireNetwork(r *http.Request) (string, error) {
	networks := r.URL.Query()["network"]
	if len(networks) != 1 {
		return "", errors.New("'network' must be specified exactly once")
	}

	if networks[0] == "" {
		return "", errors.New("'network' must not be empty")
	}

	return networks[0], nil
}

// QueryRequireAddress verify presence and return value of `address` url query parameter
func QueryRequireAddress(r *http.Request) (string, error) {
	addresses := r.URL.Query()["address"]
	if len(addresses) != 1 {
		return "", errors.New("'address' must be specified exactly once")
	}

	if addresses[0] == "" {
		return "", errors.New("'address' must not be empty")
	}

	return addresses[0], nil
}

// QueryOptionalProtocol verify presence and return value of `protocol` url query parameter
func QueryOptionalProtocol(r *http.Request) string {
	protocols := r.URL.Query()["protocol"]
	if len(protocols) != 1 {
		return ""
	}

	if protocols[0] == "" {
		return ""
	}

	return protocols[0]
}

// QueryOptionalNetwork verify presence and return value of `network` url query parameter
func QueryOptionalNetwork(r *http.Request) string {
	networks := r.URL.Query()["network"]
	if len(networks) != 1 {
		return ""
	}

	if networks[0] == "" {
		return ""
	}

	return networks[0]
}

// QueryDefaultInt verify presence and return value of parameter if empty return default value
func QueryDefaultInt(r *http.Request, name string, defaultValue int32) (int32, error) {
	val := r.URL.Query()[name]
	if len(val) != 1 {
		return defaultValue, nil
	}

	i64, err := strconv.ParseInt(val[0], 10, 32)
	if err != nil {
		return 0, err
	}

	return int32(i64), nil
}
