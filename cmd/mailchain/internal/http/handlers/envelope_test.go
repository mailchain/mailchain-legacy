package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mailchain/mailchain/stores"
	"github.com/mailchain/mailchain/stores/pinata"
	"github.com/mailchain/mailchain/stores/s3store"
	"github.com/stretchr/testify/assert"
)

func TestGetEnvelope(t *testing.T) {
	type args struct {
		sent stores.Sent
	}
	tests := []struct {
		name       string
		args       args
		wantBody   string
		wantStatus int
	}{
		{
			"",
			args{
				&s3store.Sent{},
			},
			"[{\"type\":\"0x01\",\"description\":\"Private Message Stored with MLI\"}]\n",
			200,
		},
		{
			"",
			args{
				&stores.SentStore{},
			},
			"[{\"type\":\"0x01\",\"description\":\"Private Message Stored with MLI\"}]\n",
			200,
		},
		{
			"",
			args{
				&pinata.Sent{},
			},
			"[{\"type\":\"0x02\",\"description\":\"Private Message Stored on IPFS\"}]\n",
			200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(GetEnvelope(tt.args.sent))

			req, err := http.NewRequest("GET", "/", nil)
			if err != nil {
				t.Fatal(err)
			}

			// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
			// directly and pass in our Request and ResponseRecorder.
			handler.ServeHTTP(rr, req)

			// Check the status code is what we expect.
			if !assert.Equal(t, tt.wantStatus, rr.Code) {
				t.Errorf("handler returned wrong status code: got %v want %v",
					rr.Code, tt.wantStatus)
			}
			if !assert.Equal(t, tt.wantBody, rr.Body.String()) {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), tt.wantBody)
			}
		})
	}
}
