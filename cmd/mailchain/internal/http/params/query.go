package params

import (
	"net/http"

	"github.com/pkg/errors"
)

func QueryRequireProtocol(r *http.Request) (string, error) {
	protocols := r.URL.Query()["protocol"]
	if len(protocols) != 1 {
		return "", errors.Errorf("protocol must be specified exactly once")
	}
	return protocols[0], nil
}

func QueryRequireNetwork(r *http.Request) (string, error) {
	networks := r.URL.Query()["network"]
	if len(networks) != 1 {
		return "", errors.Errorf("network must be specified exactly once")
	}
	return networks[0], nil
}
