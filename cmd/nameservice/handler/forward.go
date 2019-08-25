package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/mailchain/mailchain/errs"
	"github.com/mailchain/mailchain/internal/address"
	"github.com/mailchain/mailchain/nameservice"
	"github.com/pkg/errors"
)

func Forward(resolver nameservice.ForwardLookup) func(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Address string `json:"address"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		protocol := strings.ToLower(mux.Vars(r)["protocol"])
		network := strings.ToLower(mux.Vars(r)["network"])
		if len(r.URL.Query()["domain-name"]) != 1 {
			errs.JSONWriter(w, http.StatusPreconditionFailed, errors.Errorf("domain-name must be specified exactly once"))
			return
		}

		resolvedAddress, err := resolver.ResolveName(r.Context(), protocol, network, r.URL.Query()["domain-name"][0])
		if nameservice.IsInvalidNameError(err) {
			errs.JSONWriter(w, http.StatusPreconditionFailed, err)
			return
		}
		if nameservice.IsNoResolverError(err) {
			errs.JSONWriter(w, http.StatusNotFound, err)
			return
		}
		if nameservice.IsNotFoundError(err) {
			errs.JSONWriter(w, http.StatusNotFound, err)
			return
		}
		if err != nil {
			errs.JSONWriter(w, http.StatusInternalServerError, err)
			return
		}
		encAddress, err := address.EncodeByProtocol(resolvedAddress, protocol)
		if err != nil {
			errs.JSONWriter(w, http.StatusInternalServerError, errors.WithMessage(err, "failed to encode address"))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response{
			Address: encAddress,
		})
	}
}
