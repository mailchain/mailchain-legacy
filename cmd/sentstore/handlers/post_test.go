package handlers

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mailchain/mailchain/cmd/sentstore/storage"
	"github.com/mailchain/mailchain/cmd/sentstore/storage/storagetest"
	"github.com/mailchain/mailchain/internal/mail"
	"github.com/mailchain/mailchain/internal/encoding/encodingtest"
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
					s.EXPECT().Exists(mail.ID([]byte{}), encodingtest.MustDecodeHex("47eca011e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471"), nil, []byte("body")).Return(nil)
					s.EXPECT().Put(mail.ID([]byte{}), encodingtest.MustDecodeHex("47eca011e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471"), nil, []byte("body")).Return("https://domain.com/47eca011e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471", "47eca011e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471", uint64(1), nil)
					return s
				}(),
				stores.SizeMegabyte * 2,
			},
			func() *http.Request {
				req := httptest.NewRequest("POST", "/?contents-hash=47eca011e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471&hash=2204a9590878", nil)
				req.Body = ioutil.NopCloser(strings.NewReader("body"))
				return req
			}(),
			"https://test.com/47eca011e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471",
			http.StatusCreated,
			[]byte{},
		},
		{
			"success-base-with-slash",
			args{
				"https://test.com/",
				func() storage.Store {
					s := storagetest.NewMockStore(mockCtrl)
					s.EXPECT().Exists(mail.ID([]byte{}), encodingtest.MustDecodeHex("47eca011e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471"), nil, []byte("body")).Return(nil)
					s.EXPECT().Put(mail.ID([]byte{}), encodingtest.MustDecodeHex("47eca011e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471"), nil, []byte("body")).Return("https://domain.com/47eca011e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471", "47eca011e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471", uint64(1), nil)
					return s
				}(),
				stores.SizeMegabyte * 2,
			},
			func() *http.Request {
				req := httptest.NewRequest("POST", "/?contents-hash=47eca011e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471&hash=2204a9590878", nil)
				req.Body = ioutil.NopCloser(strings.NewReader("body"))
				return req
			}(),
			"https://test.com/47eca011e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471",
			http.StatusCreated,
			[]byte{},
		},
		{
			"err-location-does-not-match-resource",
			args{
				"https://test.com/",
				func() storage.Store {
					s := storagetest.NewMockStore(mockCtrl)
					s.EXPECT().Exists(mail.ID([]byte{}), encodingtest.MustDecodeHex("47eca011e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471"), nil, []byte("body")).Return(nil)
					s.EXPECT().Put(mail.ID([]byte{}), encodingtest.MustDecodeHex("47eca011e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471"), nil, []byte("body")).Return("https://domain.com/d12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471", "47eca011e32b52c71005ad8a8f75e1b44c92c99f", uint64(1), nil)
					return s
				}(),
				stores.SizeMegabyte * 2,
			},
			func() *http.Request {
				req := httptest.NewRequest("POST", "/?contents-hash=47eca011e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471&hash=2204a9590878", nil)
				req.Body = ioutil.NopCloser(strings.NewReader("body"))
				return req
			}(),
			"",
			http.StatusConflict,
			[]byte("{\"code\":409,\"message\":\"location does not contain resource\"}"),
		},
		{
			"err-no-contents-hash",
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
			[]byte("{\"code\":412,\"message\":\"`contents-hash` must be specified once\"}"),
		},
		{
			"err-invalid-contents-hash",
			args{
				"https://test.com/",
				func() storage.Store {
					s := storagetest.NewMockStore(mockCtrl)
					return s
				}(),
				stores.SizeMegabyte * 2,
			},
			func() *http.Request {
				req := httptest.NewRequest("POST", "/?contents-hash=47eca011e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af47&hash=2204a9590878", nil)
				req.Body = ioutil.NopCloser(strings.NewReader("body"))
				return req
			}(),
			"",
			http.StatusUnprocessableEntity,
			[]byte("{\"code\":422,\"message\":\"contents-hash invalid: encoding/hex: odd length hex string\"}"),
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
				req := httptest.NewRequest("POST", "/?contents-hash=47eca011e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471", nil)
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
				req := httptest.NewRequest("POST", "/?contents-hash=47eca011e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471&hash=", nil)
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
				req := httptest.NewRequest("POST", "/?contents-hash=47eca011e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471&hash=2204a9590878", nil)
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
				req := httptest.NewRequest("POST", "/?contents-hash=47eca011e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471&hash=2204a9590878", nil)
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
				req := httptest.NewRequest("POST", "/?contents-hash=47eca011e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471&hash=2204a9590877", nil)
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
					s.EXPECT().Exists(mail.ID([]byte{}), encodingtest.MustDecodeHex("47eca011e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471"), nil, []byte("body")).Return(errors.Errorf("already exists"))
					return s
				}(),
				stores.SizeMegabyte * 2,
			},
			func() *http.Request {
				req := httptest.NewRequest("POST", "/?contents-hash=47eca011e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471&hash=2204a9590878", nil)
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
					s.EXPECT().Exists(mail.ID([]byte{}), encodingtest.MustDecodeHex("47eca011e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471"), nil, []byte("body")).Return(nil)
					s.EXPECT().Put(mail.ID([]byte{}), encodingtest.MustDecodeHex("47eca011e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471"), nil, []byte("body")).Return("", "", uint64(1), errors.Errorf("put failed"))
					return s
				}(),
				stores.SizeMegabyte * 2,
			},
			func() *http.Request {
				req := httptest.NewRequest("POST", "/?contents-hash=47eca011e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471&hash=2204a9590878", nil)
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

func Test_getContents(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		body        io.Reader
		maxContents int
	}
	tests := []struct {
		name       string
		args       args
		want       []byte
		wantStatus int
		wantErr    bool
	}{
		{
			"success",
			args{
				ioutil.NopCloser(strings.NewReader("body")),
				stores.SizeMegabyte * 2,
			},
			[]byte{0x62, 0x6f, 0x64, 0x79},
			200,
			false,
		},
		{
			"err-empty-contents",
			args{
				ioutil.NopCloser(strings.NewReader("")),
				stores.SizeMegabyte * 2,
			},
			nil,
			412,
			true,
		},
		{
			"err-contents-too-big",
			args{
				ioutil.NopCloser(strings.NewReader("body")),
				2,
			},
			nil,
			422,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotStatus, err := getContents(tt.args.body, tt.args.maxContents)
			if (err != nil) != tt.wantErr {
				t.Errorf("getContents() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(tt.want, got) {
				t.Errorf("getContents() got = %v, want %v", got, tt.want)
			}
			if gotStatus != tt.wantStatus {
				t.Errorf("getContents() got1 = %v, wantStatus %v", gotStatus, tt.wantStatus)
			}
		})
	}
}
