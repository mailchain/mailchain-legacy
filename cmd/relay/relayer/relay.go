package relayer

import (
	"io"
	"net/http"

	"github.com/mailchain/mailchain/errs"
)

// RelayFunc definition of a relay.
type RelayFunc func(req *http.Request) (*http.Request, error)

// HandleRequest accepts all the HTTP calls and relay's them.
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

// ChangeURL changes the URL of the incoming request.
func ChangeURL(url string) RelayFunc {
	return func(inReq *http.Request) (*http.Request, error) {
		outReq, err := http.NewRequest(inReq.Method, url, inReq.Body)
		if err != nil {
			return nil, err
		}
		copyHeader(inReq.Header, outReq.Header)
		outReq.ContentLength = inReq.ContentLength
		return outReq, nil
	}
}

func copyHeader(src, dst http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
