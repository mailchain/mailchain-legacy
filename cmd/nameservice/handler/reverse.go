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

// Reverse handle forward domain lookups where a domain name is looked up to find an address.
func Reverse(resolver nameservice.ReverseLookup) func(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Name string `json:"name"`

		// The rFC1035 error status, if present
		// Since 0 status belongs to 'No Error', it's safe to use 'omitempty'
		//
		// Required: false
		// example: 3
		Status int `json:"status,omitempty"`
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
			_ = json.NewEncoder(w).Encode(response{
				Status: nameservice.RFC1035StatusMap[nameservice.ErrFormat],
			})
			return
		}
		name, err := resolver.ResolveAddress(r.Context(), protocol, network, address)
		if nameservice.IsRFC1035Error(err) {
			_ = json.NewEncoder(w).Encode(response{
				Status: nameservice.RFC1035StatusMap[err],
			})
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
