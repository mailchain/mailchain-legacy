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
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestJSONWriter(t *testing.T) {
	type args struct {
		w    http.ResponseWriter
		r    *http.Request
		code int
		err  error
	}
	tests := []struct {
		name   string
		args   args
		logOut string
	}{
		{
			"warn-error",
			args{
				&httptest.ResponseRecorder{},
				httptest.NewRequest("GET", "/address", nil),
				http.StatusPreconditionFailed,
				errors.New("pre-condition failed"),
			},
			"{\"level\":\"debug\",\"method\":\"GET\",\"url\":\"/address\",\"status\":412,\"error\":\"pre-condition failed\",\"caller\":\"/Users/robdefeo/development/GitHub/mailchain/errs/error_test.go:102\",\"message\":\"client error\"}\n",
		},
		{
			"warn-nil-error",
			args{
				&httptest.ResponseRecorder{},
				httptest.NewRequest("GET", "/address", nil),
				http.StatusPreconditionFailed,
				nil,
			},
			"{\"level\":\"debug\",\"method\":\"GET\",\"url\":\"/address\",\"status\":412,\"error\":\"no error specified\",\"caller\":\"/Users/robdefeo/development/GitHub/mailchain/errs/error_test.go:102\",\"message\":\"client error\"}\n",
		},
		{
			"err-error",
			args{
				&httptest.ResponseRecorder{},
				httptest.NewRequest("GET", "/address", nil),
				http.StatusInternalServerError,
				errors.New("internal error"),
			},
			"{\"level\":\"error\",\"method\":\"GET\",\"url\":\"/address\",\"status\":500,\"error\":\"internal error\",\"caller\":\"/Users/robdefeo/development/GitHub/mailchain/errs/error_test.go:102\",\"message\":\"server error\"}\n",
		},
		{
			"unknown-error",
			args{
				&httptest.ResponseRecorder{},
				httptest.NewRequest("GET", "/address", nil),
				600,
				errors.New("unknown error"),
			},
			"{\"level\":\"error\",\"method\":\"GET\",\"url\":\"/address\",\"status\":600,\"error\":\"unknown error\",\"caller\":\"/Users/robdefeo/development/GitHub/mailchain/errs/error_test.go:102\",\"message\":\"unknown status\"}\n",
		},
		{
			"unknown-error",
			args{
				&httptest.ResponseRecorder{},
				httptest.NewRequest("GET", "/address", nil),
				600,
				errors.New("unknown error"),
			},
			"{\"level\":\"error\",\"method\":\"GET\",\"url\":\"/address\",\"status\":600,\"error\":\"unknown error\",\"caller\":\"/Users/robdefeo/development/GitHub/mailchain/errs/error_test.go:102\",\"message\":\"unknown status\"}\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			log.Logger = zerolog.New(&buf)
			logrus.SetOutput(&buf)
			defer func() {
				log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
			}()
			JSONWriter(tt.args.w, tt.args.r, tt.args.code, tt.args.err)
			gotLogOut := buf.String()
			if !assert.Equal(t, tt.logOut, gotLogOut) {
				t.Errorf("logOut = %v, want %v", gotLogOut, tt.logOut)
			}
		})
	}
}
