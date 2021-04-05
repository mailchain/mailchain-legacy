package handlers

import (
	"fmt"
	"net/http"

	"github.com/mailchain/mailchain/errs"
	"github.com/mailchain/mailchain/internal/keystore"
	"github.com/mailchain/mailchain/internal/mailbox"
	"github.com/mailchain/mailchain/stores"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func FetchMessages(inbox stores.State, receivers map[string]mailbox.Receiver, ks keystore.Store) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := parseGetMessagesRequest(r)
		if err != nil {
			errs.JSONWriter(w, r, http.StatusUnprocessableEntity, errors.WithStack(err), log.Logger)
			return
		}

		receiver, ok := receivers[fmt.Sprintf("%s/%s", req.Protocol, req.Network)]
		if !ok {
			errs.JSONWriter(w, r, http.StatusUnprocessableEntity, errors.Errorf("receiver not supported on \"%s/%s\"", req.Protocol, req.Network), log.Logger)
			return
		}

		if receiver == nil {
			errs.JSONWriter(w, r, http.StatusUnprocessableEntity, errors.Errorf("no receiver configured for \"%s/%s\"", req.Protocol, req.Network), log.Logger)
			return
		}

		if !ks.HasAddress(req.addressBytes, req.Protocol, req.Network) {
			errs.JSONWriter(w, r, http.StatusNotAcceptable, errors.Errorf("no private key found for address"), log.Logger)
			return
		}

		transactions, err := receiver.Receive(r.Context(), req.Protocol, req.Network, req.addressBytes)
		if mailbox.IsNetworkNotSupportedError(err) {
			errs.JSONWriter(w, r, http.StatusNotAcceptable, errors.Errorf("network `%s` does not have etherscan client configured", req.Network), log.Logger)
			return
		}

		if err != nil {
			errs.JSONWriter(w, r, http.StatusInternalServerError, errors.WithStack(err), log.Logger)
			return
		}

		for i := range transactions {
			tx := transactions[i]
			if err := inbox.PutTransaction(req.Protocol, req.Network, req.addressBytes, tx); err != nil {
				errs.JSONWriter(w, r, http.StatusInternalServerError, errors.WithStack(err), log.Logger)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
	}
}
