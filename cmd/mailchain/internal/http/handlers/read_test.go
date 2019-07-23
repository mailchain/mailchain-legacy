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

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/mailchain/mailchain/internal/mail"
	"github.com/mailchain/mailchain/stores"
	"github.com/mailchain/mailchain/stores/storestest"
	"github.com/pkg/errors"
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
		{
			"err_inbox_func",
			args{
				func(messageID mail.ID) error {
					return errors.Errorf("inbox error")
				},
				func() *http.Request {
					req := httptest.NewRequest("GET", "/message_id", nil)
					req = mux.SetURLVars(req, map[string]string{
						"message_id": "47eca011e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471",
					})
					return req
				}(),
			},
			"{\"code\":422,\"message\":\"inbox error\"}\n",
			http.StatusUnprocessableEntity,
		},
		{
			"success",
			args{
				func(messageID mail.ID) error {
					return nil
				},
				func() *http.Request {
					req := httptest.NewRequest("GET", "/message_id", nil)
					req = mux.SetURLVars(req, map[string]string{
						"message_id": "47eca011e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471",
					})
					return req
				}(),
			},
			"",
			http.StatusOK,
		},
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

func TestPutRead(t *testing.T) {
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	type args struct {
		store stores.State
	}
	tests := []struct {
		name             string
		args             args
		req              *http.Request
		expectedResponse string
		expectedStatus   int
	}{
		{
			"success",
			args{
				func() stores.State {
					m := storestest.NewMockState(mockCtrl)
					m.EXPECT().PutMessageRead(mail.ID([]byte{0x47, 0xec, 0xa0, 0x11, 0xe3, 0x2b, 0x52, 0xc7, 0x10, 0x5, 0xad, 0x8a, 0x8f, 0x75, 0xe1, 0xb4, 0x4c, 0x92, 0xc9, 0x9f, 0xd1, 0x2e, 0x43, 0xbc, 0xcf, 0xe5, 0x71, 0xe3, 0xc2, 0xd1, 0x3d, 0x2e, 0x9a, 0x82, 0x6a, 0x55, 0xf, 0x5f, 0xf6, 0x3b, 0x24, 0x7a, 0xf4, 0x71})).Return(nil)
					return m
				}(),
			},
			func() *http.Request {
				req := httptest.NewRequest("PUT", "/message_id", nil)
				req = mux.SetURLVars(req, map[string]string{
					"message_id": "47eca011e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471",
				})
				return req
			}(),
			"",
			200,
		},
		{
			"err-put-failed",
			args{
				func() stores.State {
					m := storestest.NewMockState(mockCtrl)
					m.EXPECT().PutMessageRead(mail.ID([]byte{0x47, 0xec, 0xa0, 0x11, 0xe3, 0x2b, 0x52, 0xc7, 0x10, 0x5, 0xad, 0x8a, 0x8f, 0x75, 0xe1, 0xb4, 0x4c, 0x92, 0xc9, 0x9f, 0xd1, 0x2e, 0x43, 0xbc, 0xcf, 0xe5, 0x71, 0xe3, 0xc2, 0xd1, 0x3d, 0x2e, 0x9a, 0x82, 0x6a, 0x55, 0xf, 0x5f, 0xf6, 0x3b, 0x24, 0x7a, 0xf4, 0x71})).Return(errors.Errorf("failed"))
					return m
				}(),
			},
			func() *http.Request {
				req := httptest.NewRequest("PUT", "/message_id", nil)
				req = mux.SetURLVars(req, map[string]string{
					"message_id": "47eca011e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471",
				})
				return req
			}(),
			"{\"code\":422,\"message\":\"failed\"}\n",
			422,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(PutRead(tt.args.store))

			// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
			// directly and pass in our Request and ResponseRecorder.
			handler.ServeHTTP(rr, tt.req)

			// Check the status code is what we expect.
			if !assert.Equal(tt.expectedStatus, rr.Code) {
				t.Errorf("handler returned wrong status code: got %v want %v",
					rr.Code, tt.expectedStatus)
			}
			if !assert.Equal(tt.expectedResponse, rr.Body.String()) {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), tt.expectedResponse)
			}
		})
	}
}

func TestDeleteRead(t *testing.T) {
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	type args struct {
		store stores.State
	}
	tests := []struct {
		name             string
		args             args
		req              *http.Request
		expectedResponse string
		expectedStatus   int
	}{
		{
			"success",
			args{
				func() stores.State {
					m := storestest.NewMockState(mockCtrl)
					m.EXPECT().DeleteMessageRead(mail.ID([]byte{0x47, 0xec, 0xa0, 0x11, 0xe3, 0x2b, 0x52, 0xc7, 0x10, 0x5, 0xad, 0x8a, 0x8f, 0x75, 0xe1, 0xb4, 0x4c, 0x92, 0xc9, 0x9f, 0xd1, 0x2e, 0x43, 0xbc, 0xcf, 0xe5, 0x71, 0xe3, 0xc2, 0xd1, 0x3d, 0x2e, 0x9a, 0x82, 0x6a, 0x55, 0xf, 0x5f, 0xf6, 0x3b, 0x24, 0x7a, 0xf4, 0x71})).Return(nil)
					return m
				}(),
			},
			func() *http.Request {
				req := httptest.NewRequest("DELETE", "/message_id", nil)
				req = mux.SetURLVars(req, map[string]string{
					"message_id": "47eca011e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471",
				})
				return req
			}(),
			"",
			200,
		},
		{
			"err-delete-failed",
			args{
				func() stores.State {
					m := storestest.NewMockState(mockCtrl)
					m.EXPECT().DeleteMessageRead(mail.ID([]byte{0x47, 0xec, 0xa0, 0x11, 0xe3, 0x2b, 0x52, 0xc7, 0x10, 0x5, 0xad, 0x8a, 0x8f, 0x75, 0xe1, 0xb4, 0x4c, 0x92, 0xc9, 0x9f, 0xd1, 0x2e, 0x43, 0xbc, 0xcf, 0xe5, 0x71, 0xe3, 0xc2, 0xd1, 0x3d, 0x2e, 0x9a, 0x82, 0x6a, 0x55, 0xf, 0x5f, 0xf6, 0x3b, 0x24, 0x7a, 0xf4, 0x71})).Return(errors.Errorf("failed"))
					return m
				}(),
			},
			func() *http.Request {
				req := httptest.NewRequest("DELETE", "/message_id", nil)
				req = mux.SetURLVars(req, map[string]string{
					"message_id": "47eca011e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471",
				})
				return req
			}(),
			"{\"code\":422,\"message\":\"failed\"}\n",
			422,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(DeleteRead(tt.args.store))

			// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
			// directly and pass in our Request and ResponseRecorder.
			handler.ServeHTTP(rr, tt.req)

			// Check the status code is what we expect.
			if !assert.Equal(tt.expectedStatus, rr.Code) {
				t.Errorf("handler returned wrong status code: got %v want %v",
					rr.Code, tt.expectedStatus)
			}
			if !assert.Equal(tt.expectedResponse, rr.Body.String()) {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), tt.expectedResponse)
			}
		})
	}
}

func TestGetRead(t *testing.T) {
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	type args struct {
		store stores.State
	}
	tests := []struct {
		name             string
		args             args
		req              *http.Request
		expectedResponse string
		expectedStatus   int
	}{
		{
			"success-read",
			args{
				func() stores.State {
					m := storestest.NewMockState(mockCtrl)
					m.EXPECT().GetReadStatus(mail.ID([]byte{0x47, 0xec, 0xa0, 0x11, 0xe3, 0x2b, 0x52, 0xc7, 0x10, 0x5, 0xad, 0x8a, 0x8f, 0x75, 0xe1, 0xb4, 0x4c, 0x92, 0xc9, 0x9f, 0xd1, 0x2e, 0x43, 0xbc, 0xcf, 0xe5, 0x71, 0xe3, 0xc2, 0xd1, 0x3d, 0x2e, 0x9a, 0x82, 0x6a, 0x55, 0xf, 0x5f, 0xf6, 0x3b, 0x24, 0x7a, 0xf4, 0x71})).Return(true, nil)
					return m
				}(),
			},
			func() *http.Request {
				req := httptest.NewRequest("GET", "/message_id", nil)
				req = mux.SetURLVars(req, map[string]string{
					"message_id": "47eca011e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471",
				})
				return req
			}(),
			"{\"read\":true}\n",
			200,
		},
		{
			"success-unread",
			args{
				func() stores.State {
					m := storestest.NewMockState(mockCtrl)
					m.EXPECT().GetReadStatus(mail.ID([]byte{0x47, 0xec, 0xa0, 0x11, 0xe3, 0x2b, 0x52, 0xc7, 0x10, 0x5, 0xad, 0x8a, 0x8f, 0x75, 0xe1, 0xb4, 0x4c, 0x92, 0xc9, 0x9f, 0xd1, 0x2e, 0x43, 0xbc, 0xcf, 0xe5, 0x71, 0xe3, 0xc2, 0xd1, 0x3d, 0x2e, 0x9a, 0x82, 0x6a, 0x55, 0xf, 0x5f, 0xf6, 0x3b, 0x24, 0x7a, 0xf4, 0x71})).Return(false, nil)
					return m
				}(),
			},
			func() *http.Request {
				req := httptest.NewRequest("GET", "/message_id", nil)
				req = mux.SetURLVars(req, map[string]string{
					"message_id": "47eca011e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471",
				})
				return req
			}(),
			"{\"read\":false}\n",
			200,
		},
		{
			"err-get-read-status",
			args{
				func() stores.State {
					m := storestest.NewMockState(mockCtrl)
					m.EXPECT().GetReadStatus(mail.ID([]byte{0x47, 0xec, 0xa0, 0x11, 0xe3, 0x2b, 0x52, 0xc7, 0x10, 0x5, 0xad, 0x8a, 0x8f, 0x75, 0xe1, 0xb4, 0x4c, 0x92, 0xc9, 0x9f, 0xd1, 0x2e, 0x43, 0xbc, 0xcf, 0xe5, 0x71, 0xe3, 0xc2, 0xd1, 0x3d, 0x2e, 0x9a, 0x82, 0x6a, 0x55, 0xf, 0x5f, 0xf6, 0x3b, 0x24, 0x7a, 0xf4, 0x71})).Return(false, errors.Errorf("failed"))
					return m
				}(),
			},
			func() *http.Request {
				req := httptest.NewRequest("GET", "/message_id", nil)
				req = mux.SetURLVars(req, map[string]string{
					"message_id": "47eca011e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471",
				})
				return req
			}(),
			"{\"code\":500,\"message\":\"failed\"}\n",
			500,
		},
		{
			"err-not-found",
			args{
				func() stores.State {
					m := storestest.NewMockState(mockCtrl)
					m.EXPECT().GetReadStatus(mail.ID([]byte{0x47, 0xec, 0xa0, 0x11, 0xe3, 0x2b, 0x52, 0xc7, 0x10, 0x5, 0xad, 0x8a, 0x8f, 0x75, 0xe1, 0xb4, 0x4c, 0x92, 0xc9, 0x9f, 0xd1, 0x2e, 0x43, 0xbc, 0xcf, 0xe5, 0x71, 0xe3, 0xc2, 0xd1, 0x3d, 0x2e, 0x9a, 0x82, 0x6a, 0x55, 0xf, 0x5f, 0xf6, 0x3b, 0x24, 0x7a, 0xf4, 0x71})).Return(false, errors.Errorf("not found"))
					return m
				}(),
			},
			func() *http.Request {
				req := httptest.NewRequest("GET", "/message_id", nil)
				req = mux.SetURLVars(req, map[string]string{
					"message_id": "47eca011e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471",
				})
				return req
			}(),
			"",
			404,
		},
		{
			"err-message-id",
			args{
				func() stores.State {
					m := storestest.NewMockState(mockCtrl)

					return m
				}(),
			},
			func() *http.Request {
				req := httptest.NewRequest("GET", "/message_id", nil)
				req = mux.SetURLVars(req, map[string]string{
					"message_id": "47eca",
				})
				return req
			}(),
			"{\"code\":406,\"message\":\"invalid `message_id`: encoding/hex: odd length hex string\"}\n",
			406,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(GetRead(tt.args.store))

			// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
			// directly and pass in our Request and ResponseRecorder.
			handler.ServeHTTP(rr, tt.req)

			// Check the status code is what we expect.
			if !assert.Equal(tt.expectedStatus, rr.Code) {
				t.Errorf("handler returned wrong status code: got %v want %v",
					rr.Code, tt.expectedStatus)
			}
			if !assert.Equal(tt.expectedResponse, rr.Body.String()) {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), tt.expectedResponse)
			}
		})
	}
}
