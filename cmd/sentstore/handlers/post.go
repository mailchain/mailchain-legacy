package handlers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/mailchain/mailchain/cmd/sentstore/storage"
	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/errs"
	"github.com/mailchain/mailchain/internal/mail"
	"github.com/pkg/errors"
)

func PostHandler(base string, store storage.Store, maxContents int) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: check message contents don't already exist
		if len(r.URL.Query()["message-id"]) != 1 {
			errs.JSONWriter(w, http.StatusPreconditionFailed, errors.Errorf("`message-id` must be specified once"))
			return
		}
		messageID, err := mail.FromHexString(r.URL.Query()["message-id"][0])
		if err != nil {
			errs.JSONWriter(w, http.StatusUnprocessableEntity, errors.WithMessage(err, "message-id invalid"))
			return
		}

		if len(r.URL.Query()["hash"]) != 1 {
			errs.JSONWriter(w, http.StatusPreconditionFailed, errors.Errorf("`hash` must be specified once"))
			return
		}
		hash := r.URL.Query()["hash"][0]
		if hash == "" {
			errs.JSONWriter(w, http.StatusPreconditionFailed, errors.Errorf("`hash` must not be empty"))
			return
		}

		contents, err := ioutil.ReadAll(r.Body)
		if err != nil {
			errs.JSONWriter(w, http.StatusInternalServerError, errors.WithMessage(err, "contents can not be read"))
			return
		}
		if len(contents) == 0 {
			errs.JSONWriter(w, http.StatusPreconditionFailed, errors.Errorf("`contents` must not be empty"))
			return
		}
		if len(contents) > maxContents {
			errs.JSONWriter(w, http.StatusUnprocessableEntity, errors.Errorf("file size too large"))
			return
		}
		if err := compHash(contents, hash); err != nil {
			errs.JSONWriter(w, http.StatusUnprocessableEntity, errors.WithMessage(err, "`hash` invalid"))
			return
		}
		if err := store.Exists(messageID, contents, hash); err != nil {
			errs.JSONWriter(w, http.StatusConflict, errors.WithMessage(err, "conflict"))
			return
		}

		loc, err := store.Put(messageID, contents, hash)
		if err != nil {
			errs.JSONWriter(w, http.StatusInternalServerError, errors.WithMessage(err, "failed to create message"))
			return
		}

		w.Header().Set("Location", fmt.Sprintf("%s/%s", strings.TrimSuffix(base, "/"), loc))
		w.WriteHeader(http.StatusCreated)
	}
}

func compHash(contents []byte, suppliedHex string) error {
	contentsLocationHash := crypto.CreateLocationHash(contents)
	suppliedHash, err := mail.FromHexString(suppliedHex)
	if err != nil {
		return err
	}

	if !bytes.Equal(contentsLocationHash, suppliedHash) {
		return errors.Errorf("contents and supplied hash do not match")
	}
	return nil
}
