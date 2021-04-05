package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetVersion(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/version", nil)
	handler := http.HandlerFunc(GetVersion())
	handler.ServeHTTP(rr, req)
	wantStatus := http.StatusOK
	wantBody := "{\"version\":\"dev\",\"commit\":\"none\",\"time\":\"unknown\"}\n"
	if !assert.Equal(t, wantStatus, rr.Code) {
		t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, wantStatus)
	}
	if !assert.Equal(t, wantBody, rr.Body.String()) {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), wantBody)
	}
}
