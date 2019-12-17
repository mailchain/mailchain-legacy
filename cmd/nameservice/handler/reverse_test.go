package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/mailchain/mailchain/encoding/encodingtest"
	"github.com/mailchain/mailchain/nameservice"
	"github.com/mailchain/mailchain/nameservice/nameservicetest"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestReverse(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("contents"))
		}),
	)
	defer server.Close()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	assert := assert.New(t)
	type args struct {
		resolver nameservice.ReverseLookup
	}
	tests := []struct {
		name       string
		args       args
		req        *http.Request
		wantStatus int
		wantBody   string
	}{
		{
			"success",
			args{
				func() nameservice.ReverseLookup {
					m := nameservicetest.NewMockReverseLookup(mockCtrl)
					m.EXPECT().ResolveAddress(gomock.Any(), "ethereum", "mainnet", encodingtest.MustDecodeHex("5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761")).Return("test.eth", nil)
					return m
				}(),
			},
			func() *http.Request {
				req := httptest.NewRequest("GET", "/?address=0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761", nil)
				req = mux.SetURLVars(req, map[string]string{
					"protocol": "ethereum",
					"network":  "mainnet",
				})
				return req
			}(),
			200,
			"{\"name\":\"test.eth\"}\n",
		},
		{
			"err-unknown",
			args{
				func() nameservice.ReverseLookup {
					m := nameservicetest.NewMockReverseLookup(mockCtrl)
					m.EXPECT().ResolveAddress(gomock.Any(), "ethereum", "mainnet", encodingtest.MustDecodeHex("5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761")).Return("", errors.Errorf("failed"))
					return m
				}(),
			},
			func() *http.Request {
				req := httptest.NewRequest("GET", "/?address=0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761", nil)
				req = mux.SetURLVars(req, map[string]string{
					"protocol": "ethereum",
					"network":  "mainnet",
				})
				return req
			}(),
			500,
			"{\"code\":500,\"message\":\"failed\"}\n",
		},
		{
			"err-not-found",
			args{
				func() nameservice.ReverseLookup {
					m := nameservicetest.NewMockReverseLookup(mockCtrl)
					m.EXPECT().ResolveAddress(gomock.Any(), "ethereum", "mainnet", encodingtest.MustDecodeHex("5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761")).Return("", nameservice.ErrNXDomain)
					return m
				}(),
			},
			func() *http.Request {
				req := httptest.NewRequest("GET", "/?address=0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761", nil)
				req = mux.SetURLVars(req, map[string]string{
					"protocol": "ethereum",
					"network":  "mainnet",
				})
				return req
			}(),
			200,
			"{\"name\":\"\",\"status\":3}\n",
		},
		{
			"err-invalid-address",
			args{
				func() nameservice.ReverseLookup {
					m := nameservicetest.NewMockReverseLookup(mockCtrl)
					m.EXPECT().ResolveAddress(gomock.Any(), "ethereum", "mainnet", encodingtest.MustDecodeHex("5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761")).Return("", nameservice.ErrFormat)
					return m
				}(),
			},
			func() *http.Request {
				req := httptest.NewRequest("GET", "/?address=0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761", nil)
				req = mux.SetURLVars(req, map[string]string{
					"protocol": "ethereum",
					"network":  "mainnet",
				})
				return req
			}(),
			200,
			"{\"name\":\"\",\"status\":1}\n",
		},
		{
			"err-invalid-address-query",
			args{
				func() nameservice.ReverseLookup {
					m := nameservicetest.NewMockReverseLookup(mockCtrl)
					return m
				}(),
			},
			func() *http.Request {
				req := httptest.NewRequest("GET", "/?address=0x560", nil)
				req = mux.SetURLVars(req, map[string]string{
					"protocol": "ethereum",
					"network":  "mainnet",
				})
				return req
			}(),
			200,
			"{\"name\":\"\",\"status\":1}\n",
		},
		{
			"err-missing-address-query",
			args{
				func() nameservice.ReverseLookup {
					m := nameservicetest.NewMockReverseLookup(mockCtrl)
					return m
				}(),
			},
			func() *http.Request {
				req := httptest.NewRequest("GET", "/?", nil)
				req = mux.SetURLVars(req, map[string]string{
					"protocol": "ethereum",
					"network":  "mainnet",
				})
				return req
			}(),
			412,
			"{\"code\":412,\"message\":\"address must be specified exactly once\"}\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(Reverse(tt.args.resolver))

			// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
			// directly and pass in our Request and ResponseRecorder.
			handler.ServeHTTP(rr, tt.req)

			// Check the status code is what we expect.
			if !assert.Equal(tt.wantStatus, rr.Code) {
				t.Errorf("Reverse() returned wrong status code: got %v want %v",
					rr.Code, tt.wantStatus)
			}
			if !assert.Equal(tt.wantBody, rr.Body.String()) {
				t.Errorf("Reverse() returned unexpected body: got %v want %v",
					rr.Body.String(), tt.wantBody)
			}
		})
	}
}
