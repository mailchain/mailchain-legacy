package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mailchain/mailchain/cmd/relay/relayer"
	"github.com/stretchr/testify/assert"
)

func TestHandleRequest(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("contents"))
		}),
	)
	defer server.Close()

	assert := assert.New(t)
	type args struct {
		relayers map[string]relayer.RelayFunc
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
				map[string]relayer.RelayFunc{
					"ethereum/mainnet": relayer.ChangeURL(server.URL),
				},
			},
			httptest.NewRequest("GET", "/ethereum/mainnet", nil),
			200,
			"contents",
		},
		{
			"success-trailing-slash",
			args{
				map[string]relayer.RelayFunc{
					"ethereum/mainnet": relayer.ChangeURL(server.URL),
				},
			},
			httptest.NewRequest("GET", "/ethereum/mainnet/", nil),
			200,
			"contents",
		},
		{
			"err-no-relay",
			args{
				map[string]relayer.RelayFunc{
					"ethereum/testnet": relayer.ChangeURL(server.URL),
				},
			},
			httptest.NewRequest("GET", "/ethereum/mainnet/", nil),
			500,
			"{\"code\":500,\"message\":\"unknown relay destination for \\\"ethereum/mainnet\\\"\"}\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(HandleRequest(tt.args.relayers))

			// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
			// directly and pass in our Request and ResponseRecorder.
			handler.ServeHTTP(rr, tt.req)

			// Check the status code is what we expect.
			if !assert.Equal(tt.wantStatus, rr.Code) {
				t.Errorf("handler returned wrong status code: got %v want %v",
					rr.Code, tt.wantStatus)
			}
			if !assert.Equal(tt.wantBody, rr.Body.String()) {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), tt.wantBody)
			}
		})
	}
}
