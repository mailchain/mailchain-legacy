package relayer

import (
	"io"
	"net/http"

	"github.com/mailchain/mailchain/errs"
)

type RelayFunc func(req *http.Request) (*http.Request, error)

// ServeHTTP calls f(w, r).
func (f RelayFunc) HandleRequest(w http.ResponseWriter, req *http.Request) {
	r, err := f(req)
	if err != nil {
		errs.JSONWriter(w, http.StatusBadRequest, err)
		return
	}
	client := http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		errs.JSONWriter(w, http.StatusBadGateway, err)
		return
	}
	defer resp.Body.Close()
	copyHeader(resp.Header, w.Header())
	w.WriteHeader(resp.StatusCode)
	_, _ = io.Copy(w, resp.Body)
}

func ChangeURL(url string) RelayFunc {
	return func(req *http.Request) (*http.Request, error) {
		proxyReq, err := http.NewRequest(req.Method, url, req.Body)
		if err != nil {
			return nil, err
		}
		copyHeader(req.Header, proxyReq.Header)
		return proxyReq, nil
	}
}

func copyHeader(src, dst http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
