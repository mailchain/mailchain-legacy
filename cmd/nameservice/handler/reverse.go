package handler

import (
	"encoding/json"
	"net/http"

	"github.com/mailchain/mailchain/errs"
	"github.com/mailchain/mailchain/internal/addressing"
	"github.com/mailchain/mailchain/nameservice"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// Reverse handle forward domain lookups where a domain name is looked up to find an address.
func Reverse(resolver nameservice.ReverseLookup, protocol, network string) func(w http.ResponseWriter, r *http.Request) {
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
		if len(r.URL.Query()["address"]) != 1 {
			errs.JSONWriter(w, r, http.StatusPreconditionFailed, errors.Errorf("address must be specified exactly once"), log.Logger)
			return
		}

		addr, err := addressing.DecodeByProtocol(r.URL.Query()["address"][0], protocol)
		if err != nil {
			_ = json.NewEncoder(w).Encode(response{
				Status: nameservice.ErrorToRFC1035Status(nameservice.ErrFormat),
			})

			return
		}

		name, err := resolver.ResolveAddress(r.Context(), protocol, network, addr)
		if nameservice.ErrorToRFC1035Status(err) > 0 {
			_ = json.NewEncoder(w).Encode(response{
				Status: nameservice.ErrorToRFC1035Status(err),
			})

			return
		}

		if err != nil {
			errs.JSONWriter(w, r, http.StatusInternalServerError, err, log.Logger)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response{
			Name: name,
		})
	}
}
