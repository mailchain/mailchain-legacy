// Copyright 2019 Finobo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package errs

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus" //nolint:depguard
	// TODO: pass stdout and stderr as params
)

// A function called whenever an error is encountered
// type errorHandler func(w http.ResponseWriter, r *http.Request, err string)
type ErrorWriter func(w http.ResponseWriter, code int, err error)

// errorf writes a swagger-compliant error response.
func JSONWriter(w http.ResponseWriter, code int, err error) {
	var out struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	out.Code = code
	out.Message = fmt.Sprint(err)

	b, err := json.Marshal(out)
	if err != nil {
		http.Error(w, `{"code": 500, "message": "Could not format JSON for original message."}`, 500)
		return
	}
	http.Error(w, string(b), code)

	switch out.Code {
	case http.StatusInternalServerError,
		http.StatusNotImplemented,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout,
		http.StatusHTTPVersionNotSupported,
		http.StatusVariantAlsoNegotiates,
		http.StatusInsufficientStorage,
		http.StatusLoopDetected,
		http.StatusNotExtended,
		http.StatusNetworkAuthenticationRequired:
		logrus.Error(err)
	default:
		logrus.Error(err)
	}
}
