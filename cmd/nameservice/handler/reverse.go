package handler

import (
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/mailchain/mailchain/errs"
	"github.com/mailchain/mailchain/nameservice"
	"github.com/pkg/errors"
)

func Reverse(resolver nameservice.ReverseLookup) func(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Name string `json:"name"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		protocol := strings.ToLower(mux.Vars(r)["protocol"])
		network := strings.ToLower(mux.Vars(r)["network"])
		if len(r.URL.Query()["address"]) != 1 {
			errs.JSONWriter(w, http.StatusPreconditionFailed, errors.Errorf("address must be specified exactly once"))
			return
		}
		address, err := hex.DecodeString(strings.TrimPrefix(r.URL.Query()["address"][0], "0x"))
		if err != nil {
			errs.JSONWriter(w, http.StatusPreconditionFailed, nameservice.ErrInvalidAddress)
			return
		}
		name, err := resolver.ResolveAddress(r.Context(), protocol, network, address)
		if nameservice.IsInvalidAddressError(err) {
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

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response{
			Name: name,
		})
	}
}
