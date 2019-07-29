package params

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestQueryRequireProtocol(t *testing.T) {
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			"success",
			args{
				func() *http.Request {
					req := httptest.NewRequest("GET", "/?protocol=ethereum", nil)
					return req
				}(),
			},
			"ethereum",
			false,
		},
		{
			"err-missing",
			args{
				func() *http.Request {
					req := httptest.NewRequest("GET", "/", nil)
					return req
				}(),
			},
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := QueryRequireProtocol(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("QueryRequireProtocol() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("QueryRequireProtocol() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQueryRequireNetwork(t *testing.T) {
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			"success",
			args{
				func() *http.Request {
					req := httptest.NewRequest("GET", "/?network=mainnet", nil)
					return req
				}(),
			},
			"mainnet",
			false,
		},
		{
			"err-missing",
			args{
				func() *http.Request {
					req := httptest.NewRequest("GET", "/", nil)
					return req
				}(),
			},
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := QueryRequireNetwork(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("QueryRequireNetwork() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("QueryRequireNetwork() = %v, want %v", got, tt.want)
			}
		})
	}
}
