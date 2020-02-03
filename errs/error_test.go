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

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestJSONWriter(t *testing.T) {
	type args struct {
		w    http.ResponseWriter
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
				http.StatusPreconditionFailed,
				errors.New("pre-condition failed"),
			},
			"level=warning msg=\"status 412: pre-condition failed\"\n",
		},
		{
			"warn-nil-error",
			args{
				&httptest.ResponseRecorder{},
				http.StatusPreconditionFailed,
				nil,
			},
			"level=warning msg=\"status 412: no error specified\"\n",
		},
		{
			"err-error",
			args{
				&httptest.ResponseRecorder{},
				http.StatusInternalServerError,
				errors.New("internal error"),
			},
			"level=error msg=\"status 500: internal error\"\n",
		},
		{
			"unknown-error",
			args{
				&httptest.ResponseRecorder{},
				600,
				errors.New("unknown error"),
			},
			"level=error msg=\"unknown status 600: unknown error\"\n",
		},
		{
			"unknown-error",
			args{
				&httptest.ResponseRecorder{},
				600,
				errors.New("unknown error"),
			},
			"level=error msg=\"unknown status 600: unknown error\"\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logrus.SetFormatter(&logrus.TextFormatter{
				DisableColors:    true,
				DisableTimestamp: true,
			})
			logrus.SetOutput(&buf)
			defer func() {
				logrus.SetOutput(os.Stderr)
				logrus.SetFormatter(&logrus.TextFormatter{})
			}()
			JSONWriter(tt.args.w, tt.args.code, tt.args.err)
			gotLogOut := buf.String()
			if !assert.Equal(t, tt.logOut, gotLogOut) {
				t.Errorf("logOut = %v, want %v", gotLogOut, tt.logOut)
			}
		})
	}
}
