package handler

import (
	"net/http"
	"strings"

	"github.com/mailchain/mailchain/cmd/relay/relayer"
	"github.com/mailchain/mailchain/errs"
	"github.com/pkg/errors"
)

// HandleRequest accepts all relay requests and routes then to the new URL as required.
func HandleRequest(relayers map[string]relayer.RelayFunc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		path := strings.Trim(req.URL.Path, "/")
		relay, ok := relayers[path]
		if !ok {
			errs.JSONWriter(w, http.StatusInternalServerError, errors.Errorf("unknown relay destination for %q", path))
			return
		}
		relay.HandleRequest(w, req)
	}
}
