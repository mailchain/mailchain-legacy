package handlers

import (
	"encoding/json"
	"github.com/mailchain/mailchain"
	"net/http"
)

func GetResolveVersion() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(GetResolveVersionResponseBody{
			VersionTag:    mailchain.Version,
			VersionCommit: mailchain.Commit,
			VersionDate:   mailchain.Date,
		})
	}
}

// GetBody body response
//
// swagger:model GetResolveVersionResponseBody
type GetResolveVersionResponseBody struct {
	// The resolved version tag
	// Required: true
	// example: 1.0.0
	VersionTag string `json:"version"`
	// The resolved version commit
	// Required: true
	VersionCommit string `json:"commit"`
	// The resolved version release date
	// Required: true
	// example: 2019-09-04T21:59:26Z
	VersionDate string `json:"time"`
}
