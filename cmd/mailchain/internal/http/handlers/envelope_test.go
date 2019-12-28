package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnvelope(t *testing.T) {
	assert := assert.New(t)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/envelope", nil)
	handler := http.HandlerFunc(GetEnvelope())
	handler.ServeHTTP(rr, req)
	wantStatus := http.StatusOK
	wantBody := "[{\"type\":\"0x01\",\"description\":\"Private Message Stored with MLI\"}]\n"
	if !assert.Equal(wantStatus, rr.Code) {
		t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, wantStatus)
	}
	if !assert.Equal(wantBody, rr.Body.String()) {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), wantBody)
	}
}
