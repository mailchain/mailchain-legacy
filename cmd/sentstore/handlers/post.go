package handlers

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/mailchain/mailchain/cmd/sentstore/storage"
	"github.com/mailchain/mailchain/errs"
	"github.com/mailchain/mailchain/internal/hash"
	"github.com/mailchain/mailchain/internal/mail"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// PostHandler stores message in the configured file storage
func PostHandler(base string, store storage.Store, maxContents int) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		messageID := []byte{}
		if len(r.URL.Query()["contents-hash"]) != 1 {
			errs.JSONWriter(w, r, http.StatusPreconditionFailed, errors.Errorf("`contents-hash` must be specified once"), log.Logger)
			return
		}
		contentsHash, err := hex.DecodeString(r.URL.Query()["contents-hash"][0])
		if err != nil {
			errs.JSONWriter(w, r, http.StatusUnprocessableEntity, errors.WithMessage(err, "contents-hash invalid"), log.Logger)
			return
		}

		if len(r.URL.Query()["hash"]) != 1 {
			errs.JSONWriter(w, r, http.StatusPreconditionFailed, errors.Errorf("`hash` must be specified once"), log.Logger)
			return
		}
		hash := r.URL.Query()["hash"][0]
		if hash == "" {
			errs.JSONWriter(w, r, http.StatusPreconditionFailed, errors.Errorf("`hash` must not be empty"), log.Logger)
			return
		}
		contents, status, err := getContents(r.Body, maxContents)
		if err != nil {
			errs.JSONWriter(w, r, status, err, log.Logger)
			return
		}

		if err := compHash(contents, hash); err != nil {
			errs.JSONWriter(w, r, http.StatusUnprocessableEntity, errors.WithMessage(err, "`hash` invalid"), log.Logger)
			return
		}
		if err := store.Exists(messageID, contentsHash, nil, contents); err != nil {
			errs.JSONWriter(w, r, http.StatusConflict, errors.WithMessage(err, "conflict"), log.Logger)
			return
		}

		address, resource, mli, err := store.Put(messageID, contentsHash, nil, contents)
		if err != nil {
			errs.JSONWriter(w, r, http.StatusInternalServerError, errors.WithMessage(err, "failed to create message"), log.Logger)
			return
		}
		if !strings.Contains(strings.TrimSuffix(address, "/"), resource) {
			errs.JSONWriter(w, r, http.StatusConflict, errors.Errorf("location does not contain resource"), log.Logger)
			return
		}
		w.Header().Set("Location", fmt.Sprintf("%s/%s", strings.TrimSuffix(base, "/"), resource))
		w.Header().Set("Resource", resource)
		w.Header().Set("Message-Location-Identifier", strconv.FormatUint(mli, 10))
		w.WriteHeader(http.StatusCreated)
	}
}

func getContents(body io.Reader, maxContents int) (contents []byte, status int, err error) {
	contents, err = ioutil.ReadAll(body)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.WithMessage(err, "contents can not be read")
	}
	if len(contents) == 0 {
		return nil, http.StatusPreconditionFailed, errors.Errorf("`contents` must not be empty")
	}
	if len(contents) > maxContents {
		return nil, http.StatusUnprocessableEntity, errors.Errorf("file size too large")
	}
	return contents, http.StatusOK, nil
}

func compHash(contents []byte, suppliedHex string) error {
	integrityHash := hash.CreateIntegrityHash(contents)
	suppliedHash, err := mail.FromHexString(suppliedHex)
	if err != nil {
		return err
	}

	if !bytes.Equal(integrityHash, suppliedHash) {
		return errors.Errorf("contents and supplied hash do not match")
	}
	return nil
}
