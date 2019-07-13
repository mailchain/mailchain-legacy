package handlers

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mailchain/mailchain/cmd/sentstore/storage"
	"github.com/mailchain/mailchain/cmd/sentstore/storage/storagetest"
	"github.com/mailchain/mailchain/internal/mail"
	"github.com/mailchain/mailchain/internal/testutil"
	"github.com/mailchain/mailchain/stores"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func Test_compHash(t *testing.T) {
	type args struct {
		contents     []byte
		suppliedHash string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"success",
			args{
				func() []byte {
					c, err := ioutil.ReadFile("./testdata/ce9fb2801637d78b9fdc4b028fc437ac477fdc35f0409fb381757b82bc232fff27421b0cec283a65fbb81a65-22048866a9d5")
					if err != nil {
						t.Error(err)
					}
					return c
				}(),
				"22048866a9d5",
			},
			false,
		},
		{
			"err-hash-does-not-match-contents",
			args{
				[]byte("invalid-hash-data"),
				"22048866a9d5",
			},
			true,
		},
		{
			"err-invalid-hex",
			args{
				[]byte("invalid-hash-data"),
				"22048866aTd5",
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := compHash(tt.args.contents, tt.args.suppliedHash); (err != nil) != tt.wantErr {
				t.Errorf("compHash() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPostHandler(t *testing.T) {
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type args struct {
		base    string
		store   storage.Store
		maxSize int
	}
	tests := []struct {
		name               string
		args               args
		req                *http.Request
		wantLocationHeader string
		wantStatus         int
		wantBody           []byte
	}{
		{
			"success",
			args{
				"https://test.com",
				func() storage.Store {
					s := storagetest.NewMockStore(mockCtrl)
					var id mail.ID
					id = testutil.MustHexDecodeString("47eca011e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471")

					s.EXPECT().Exists(id, []byte("body"), "2204a9590878").Return(nil)
					s.EXPECT().Put(id, []byte("body"), "2204a9590878").Return("47eca011e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471-2204a9590878", nil)
					return s
				}(),
				stores.SizeMegabyte * 2,
			},
			func() *http.Request {
				req := httptest.NewRequest("POST", "/?message-id=47eca011e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471&hash=2204a9590878", nil)
				req.Body = ioutil.NopCloser(strings.NewReader("body"))
				return req
			}(),
			"https://test.com/47eca011e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471-2204a9590878",
			http.StatusCreated,
			[]byte{},
		},
		{
			"success-base-with-slash",
			args{
				"https://test.com/",
				func() storage.Store {
					s := storagetest.NewMockStore(mockCtrl)
					var id mail.ID
					id = testutil.MustHexDecodeString("47eca011e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471")

					s.EXPECT().Exists(id, []byte("body"), "2204a9590878").Return(nil)
					s.EXPECT().Put(id, []byte("body"), "2204a9590878").Return("47eca011e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471-2204a9590878", nil)
					return s
				}(),
				stores.SizeMegabyte * 2,
			},
			func() *http.Request {
				req := httptest.NewRequest("POST", "/?message-id=47eca011e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471&hash=2204a9590878", nil)
				req.Body = ioutil.NopCloser(strings.NewReader("body"))
				return req
			}(),
			"https://test.com/47eca011e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471-2204a9590878",
			http.StatusCreated,
			[]byte{},
		},
		{
			"err-no-message-id",
			args{
				"https://test.com/",
				func() storage.Store {
					s := storagetest.NewMockStore(mockCtrl)
					return s
				}(),
				stores.SizeMegabyte * 2,
			},
			func() *http.Request {
				req := httptest.NewRequest("POST", "/?hash=2204a9590878", nil)
				req.Body = ioutil.NopCloser(strings.NewReader("body"))
				return req
			}(),
			"",
			http.StatusPreconditionFailed,
			[]byte("{\"code\":412,\"message\":\"`message-id` must be specified once\"}"),
		},
		{
			"err-invalid-message-id",
			args{
				"https://test.com/",
				func() storage.Store {
					s := storagetest.NewMockStore(mockCtrl)
					return s
				}(),
				stores.SizeMegabyte * 2,
			},
			func() *http.Request {
				req := httptest.NewRequest("POST", "/?message-id=47eca011e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af47&hash=2204a9590878", nil)
				req.Body = ioutil.NopCloser(strings.NewReader("body"))
				return req
			}(),
			"",
			http.StatusUnprocessableEntity,
			[]byte("{\"code\":422,\"message\":\"message-id invalid: encoding/hex: odd length hex string\"}"),
		},
		{
			"err-missing-hash",
			args{
				"https://test.com",
				func() storage.Store {
					s := storagetest.NewMockStore(mockCtrl)
					return s
				}(),
				stores.SizeMegabyte * 2,
			},
			func() *http.Request {
				req := httptest.NewRequest("POST", "/?message-id=47eca011e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471", nil)
				req.Body = ioutil.NopCloser(strings.NewReader("body"))
				return req
			}(),
			"",
			http.StatusPreconditionFailed,
			[]byte("{\"code\":412,\"message\":\"`hash` must be specified once\"}"),
		},
		{
			"err-empty-hash",
			args{
				"https://test.com",
				func() storage.Store {
					s := storagetest.NewMockStore(mockCtrl)
					return s
				}(),
				stores.SizeMegabyte * 2,
			},
			func() *http.Request {
				req := httptest.NewRequest("POST", "/?message-id=47eca011e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471&hash=", nil)
				req.Body = ioutil.NopCloser(strings.NewReader("body"))
				return req
			}(),
			"",
			http.StatusPreconditionFailed,
			[]byte("{\"code\":412,\"message\":\"`hash` must not be empty\"}"),
		},
		{
			"err-empty-contents",
			args{
				"https://test.com",
				func() storage.Store {
					s := storagetest.NewMockStore(mockCtrl)
					return s
				}(),
				stores.SizeMegabyte * 2,
			},
			func() *http.Request {
				req := httptest.NewRequest("POST", "/?message-id=47eca011e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471&hash=2204a9590878", nil)
				return req
			}(),
			"",
			http.StatusPreconditionFailed,
			[]byte("{\"code\":412,\"message\":\"`contents` must not be empty\"}"),
		},
		{
			"err-contents-too-big",
			args{
				"https://test.com/",
				func() storage.Store {
					s := storagetest.NewMockStore(mockCtrl)

					return s
				}(),
				2,
			},
			func() *http.Request {
				req := httptest.NewRequest("POST", "/?message-id=47eca011e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471&hash=2204a9590878", nil)
				req.Body = ioutil.NopCloser(strings.NewReader("body"))
				return req
			}(),
			"",
			http.StatusUnprocessableEntity,
			[]byte("{\"code\":422,\"message\":\"file size too large\"}"),
		},
		{
			"err-comp-hash",
			args{
				"https://test.com",
				func() storage.Store {
					s := storagetest.NewMockStore(mockCtrl)

					return s
				}(),
				stores.SizeMegabyte * 2,
			},
			func() *http.Request {
				req := httptest.NewRequest("POST", "/?message-id=47eca011e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471&hash=2204a9590877", nil)
				req.Body = ioutil.NopCloser(strings.NewReader("body"))
				return req
			}(),
			"",
			http.StatusUnprocessableEntity,
			[]byte("{\"code\":422,\"message\":\"`hash` invalid: contents and supplied hash do not match\"}"),
		},
		{
			"err-exists",
			args{
				"https://test.com",
				func() storage.Store {
					s := storagetest.NewMockStore(mockCtrl)
					var id mail.ID
					id = testutil.MustHexDecodeString("47eca011e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471")

					s.EXPECT().Exists(id, []byte("body"), "2204a9590878").Return(errors.Errorf("already exists"))
					return s
				}(),
				stores.SizeMegabyte * 2,
			},
			func() *http.Request {
				req := httptest.NewRequest("POST", "/?message-id=47eca011e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471&hash=2204a9590878", nil)
				req.Body = ioutil.NopCloser(strings.NewReader("body"))
				return req
			}(),
			"",
			http.StatusConflict,
			[]byte("{\"code\":409,\"message\":\"conflict: already exists\"}"),
		},
		{
			"err-put-failed",
			args{
				"https://test.com",
				func() storage.Store {
					s := storagetest.NewMockStore(mockCtrl)
					var id mail.ID
					id = testutil.MustHexDecodeString("47eca011e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471")

					s.EXPECT().Exists(id, []byte("body"), "2204a9590878").Return(nil)
					s.EXPECT().Put(id, []byte("body"), "2204a9590878").Return("", errors.Errorf("put failed"))
					return s
				}(),
				stores.SizeMegabyte * 2,
			},
			func() *http.Request {
				req := httptest.NewRequest("POST", "/?message-id=47eca011e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471&hash=2204a9590878", nil)
				req.Body = ioutil.NopCloser(strings.NewReader("body"))
				return req
			}(),
			"",
			http.StatusInternalServerError,
			[]byte("{\"code\":500,\"message\":\"failed to create message: put failed\"}"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(PostHandler(tt.args.base, tt.args.store, tt.args.maxSize))

			// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
			// directly and pass in our Request and ResponseRecorder.
			handler.ServeHTTP(rr, tt.req)

			// Check the status code is what we expect.
			if !assert.Equal(tt.wantStatus, rr.Code) {
				t.Errorf("handler returned wrong status code: got %v want %v",
					rr.Code, tt.wantStatus)
			}
			if !assert.Equal(tt.wantLocationHeader, rr.Header().Get("Location")) {
				t.Errorf("handler returned unexpected header: got %v want %v",
					rr.Header().Get("Location"), tt.wantLocationHeader)
			}
			body, _ := ioutil.ReadAll(rr.Body)
			if !assert.EqualValues(bytes.TrimSuffix(tt.wantBody, []byte{0x2, 0x2}), bytes.TrimSuffix(body, []byte{0xa})) {
				t.Errorf("handler returned unexpected body: got %v want %v",
					string(body), string(tt.wantBody))
			}
		})
	}
}
