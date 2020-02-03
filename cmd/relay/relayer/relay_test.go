package relayer

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func Test_copyHeader(t *testing.T) {
	type args struct {
		src http.Header
		dst http.Header
	}
	tests := []struct {
		name string
		args args
		want http.Header
	}{
		{
			"success",
			args{
				http.Header{
					"Origin": []string{"localhost"},
				},
				http.Header{
					"Content-Type": []string{"application/json"},
				},
			},
			http.Header{
				"Origin":       []string{"localhost"},
				"Content-Type": []string{"application/json"},
			},
		},
		{
			"append",
			args{
				http.Header{
					"Origin":       []string{"localhost"},
					"Content-Type": []string{"application/xml"},
				},
				http.Header{
					"Content-Type": []string{"application/json"},
				},
			},
			http.Header{
				"Origin":       []string{"localhost"},
				"Content-Type": []string{"application/json", "application/xml"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			copyHeader(tt.args.src, tt.args.dst)
			if !assert.Equal(t, tt.want, tt.args.dst) {
				t.Errorf("want %v dst %v", tt.want, tt.args.dst)
			}
		})
	}
}

func TestChangeURL(t *testing.T) {
	type args struct {
		url string
		req *http.Request
	}
	type want struct {
		url string
	}
	tests := []struct {
		name    string
		args    args
		want    *want
		wantErr bool
	}{
		{
			"success",
			args{
				"http://new-host.com",
				func() *http.Request {
					r, _ := http.NewRequest("GET", "http://location", nil)
					return r
				}(),
			},
			&want{
				"http://new-host.com",
			},
			false,
		},
		{
			"bad-request",
			args{
				"foo.html",
				func() *http.Request {
					r := &http.Request{
						Method: "\"GET",
					}
					return r
				}(),
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := ChangeURL(tt.args.url)(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChangeURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (req == nil) != (tt.want == nil) {
				t.Errorf("req = %v, want %v", req, tt.want)
				return
			}
			if tt.want != nil && req != nil {
				assert.Equal(t, tt.want.url, req.URL.String())
			}

		})
	}
}

func TestRelayFunc_HandleRequest(t *testing.T) {

	type args struct {
		w      *httptest.ResponseRecorder
		req    *http.Request
		f      RelayFunc
		server *httptest.Server
	}
	tests := []struct {
		name           string
		args           args
		expectedStatus int
	}{
		{
			"success",
			func() args {
				server := httptest.NewServer(
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						// w.Write([]byte("contents"))
					}))

				return args{
					httptest.NewRecorder(),
					httptest.NewRequest("", "http://random", nil),
					ChangeURL(server.URL),
					server,
				}
			}(),
			http.StatusOK,
		},
		{
			"err-do",
			func() args {
				server := httptest.NewServer(
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						// w.Write([]byte("contents"))
					}))

				return args{
					httptest.NewRecorder(),
					httptest.NewRequest("", "http://random", nil),
					func(req *http.Request) (*http.Request, error) {
						return httptest.NewRequest(req.Method, server.URL, nil), nil
					},
					server,
				}
			}(),
			http.StatusBadGateway,
		},
		{
			"err-relay-func",
			func() args {
				server := httptest.NewServer(
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						// w.Write([]byte("contents"))
					}))

				return args{
					httptest.NewRecorder(),
					httptest.NewRequest("", "http://random", nil),
					func(req *http.Request) (*http.Request, error) {
						return nil, errors.Errorf("failed")
					},
					server,
				}
			}(),
			http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		defer tt.args.server.Close()
		t.Run(tt.name, func(t *testing.T) {
			tt.args.f.HandleRequest(tt.args.w, tt.args.req)
		})
		if !assert.Equal(t, tt.expectedStatus, tt.args.w.Result().StatusCode) {
			t.Errorf("expectedStatus %v != result %v", tt.expectedStatus, tt.args.w.Result().StatusCode)
		}
	}
}
