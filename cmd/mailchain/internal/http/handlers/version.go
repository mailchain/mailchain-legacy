package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/mailchain/mailchain"
)

// GetVersion handler that the running version.
func GetVersion() func(w http.ResponseWriter, r *http.Request) {
	// Get swagger:route GET /version Version GetVersion
	//
	// Get version
	//
	// Get version of the running mailchain client application and API.
	// This method be used to determine what version of the API and client is being used and what functionality.
	//
	// Responses:
	//   200: GetVersionResponseBody
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(GetVersionResponseBody{
			VersionTag:    mailchain.Version,
			VersionCommit: mailchain.Commit,
			VersionDate:   mailchain.Date,
		})
	}
}

// GetVersionResponse version response
//
// swagger:response GetVersionResponse
type GetVersionResponse struct {
	// in: body
	Body GetVersionResponseBody
}

// GetVersionResponseBody response
//
// swagger:model GetVersionResponseBody
type GetVersionResponseBody struct {
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
