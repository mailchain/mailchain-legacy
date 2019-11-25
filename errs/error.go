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

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus" //nolint:depguard
)

// HTTPError definition
type HTTPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ErrorWriter is the function definition called when writing a HTTP error.
type ErrorWriter func(w http.ResponseWriter, code int, err error)

// JSONWriter writes a swagger-compliant error response.
func JSONWriter(w http.ResponseWriter, code int, err error) {
	if err == nil {
		err = errors.Errorf("no error specified")
	}

	out := HTTPError{
		Code:    code,
		Message: fmt.Sprint(err),
	}

	// this can not fail as the error is a string
	b, _ := json.Marshal(out)
	http.Error(w, string(b), code)

	switch out.Code {
	case http.StatusPreconditionFailed,
		http.StatusMethodNotAllowed,
		http.StatusNotFound:
		logrus.Warnf("status %v: %v", out.Code, err)
	case http.StatusInternalServerError,
		http.StatusNotImplemented,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout,
		http.StatusHTTPVersionNotSupported,
		http.StatusVariantAlsoNegotiates,
		http.StatusInsufficientStorage,
		http.StatusLoopDetected,
		http.StatusNotExtended, http.StatusNetworkAuthenticationRequired:
		logrus.Errorf("status %v: %+v", out.Code, err)
	default:
		logrus.Errorf("unknown status %v: %+v", out.Code, err)
	}
}
