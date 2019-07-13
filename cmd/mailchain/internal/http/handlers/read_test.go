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

package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/mailchain/mailchain/internal/mail"
	"github.com/stretchr/testify/assert"
)

func Test_doRead(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		inboxFunc func(messageID mail.ID) error
		r         *http.Request
	}
	tests := []struct {
		name             string
		args             args
		expectedResponse string
		expectedStatus   int
	}{
		{
			"err_invalid_message_id",
			args{
				func(messageID mail.ID) error {
					return nil
				},
				func() *http.Request {
					req := httptest.NewRequest("GET", "/message_id", nil)
					req = mux.SetURLVars(req, map[string]string{
						"message_id": "47eca01e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471",
					})
					return req
				}(),
			},
			"{\"code\":406,\"message\":\"invalid `message_id`: encoding/hex: odd length hex string\"}\n",
			http.StatusNotAcceptable,
		},
		// {
		// 	"err_inbox_func",
		// 	args{
		// 		func(messageID mail.ID) error {
		// 			return errors.Errorf("inbox error")
		// 		},
		// 		func() *http.Request {
		// 			req := httptest.NewRequest("GET", "/message_id", nil)
		// 			req = mux.SetURLVars(req, map[string]string{
		// 				"message_id": "47eca011e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471",
		// 			})
		// 			return req
		// 		}(),
		// 	},
		// 	"{\"code\":422,\"message\":\"inbox error\"}\n",
		// 	http.StatusUnprocessableEntity,
		// },
		// {
		// 	"success",
		// 	args{
		// 		func(messageID mail.ID) error {
		// 			return nil
		// 		},
		// 		func() *http.Request {
		// 			req := httptest.NewRequest("GET", "/message_id", nil)
		// 			req = mux.SetURLVars(req, map[string]string{
		// 				"message_id": "47eca011e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471",
		// 			})
		// 			return req
		// 		}(),
		// 	},
		// 	"",
		// 	http.StatusOK,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &httptest.ResponseRecorder{Body: bytes.NewBuffer([]byte{})}
			doRead(tt.args.inboxFunc, resp, tt.args.r)
			// Check the response body is what we expect.
			// o := io.Reader

			if !assert.Equal(tt.expectedResponse, resp.Body.String()) {
				t.Errorf("handler returned unexpected body: got %v want %v",
					resp.Body.String(), tt.expectedResponse)
			}
			if resp.Code != tt.expectedStatus {
				t.Errorf("handler returned unexpected status: got %v want %v",
					resp.Code, tt.expectedStatus)
			}
		})
	}
}

// func TestPutRead(t *testing.T) {
// 	assert := assert.New(t)
// 	type args struct {
// 		store stores.State
// 	}
// 	tests := []struct {
// 		name             string
// 		args             args
// 		expectedResponse string
// 		expectedStatus   int
// 	}{
// 		{
// "success",
// args{}
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			req, err := http.NewRequest("GET", "/", nil)
// 			if err != nil {
// 				t.Fatal(err)
// 			}
// 			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
// 			rr := httptest.NewRecorder()
// 			handler := http.HandlerFunc(PutRead(tt.args.store))

// 			// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
// 			// directly and pass in our Request and ResponseRecorder.
// 			handler.ServeHTTP(rr, req)

// 			// Check the status code is what we expect.
// 			if !assert.Equal(tt.expectedStatus, rr.Code) {
// 				t.Errorf("handler returned wrong status code: got %v want %v",
// 					rr.Code, tt.expectedStatus)
// 			}
// 			if !assert.Equal(tt.expectedResponse, rr.Body.String()) {
// 				t.Errorf("handler returned unexpected body: got %v want %v",
// 					rr.Body.String(), tt.expectedResponse)
// 			}
// 		})
// 	}
// }
