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
		{
			"err-missing",
			args{
				func() *http.Request {
					req := httptest.NewRequest("GET", "/?protocol=", nil)
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
		{
			"err-empty",
			args{
				func() *http.Request {
					req := httptest.NewRequest("GET", "/?network=", nil)
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

func TestQueryRequireAddresses(t *testing.T) {
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
					req := httptest.NewRequest("GET", "/?address=0xde0b295669a9fd93d5f28d9ec85e40f4cb697bae", nil)
					return req
				}(),
			},
			"0xde0b295669a9fd93d5f28d9ec85e40f4cb697bae",
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
		{
			"err-empty",
			args{
				func() *http.Request {
					req := httptest.NewRequest("GET", "/?address=", nil)
					return req
				}(),
			},
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := QueryRequireAddress(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestQueryRequireAddresses() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("TestQueryRequireAddresses() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQueryOptionalProtocol(t *testing.T) {
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name string
		args args
		want string
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
		},
		{
			"missing",
			args{
				func() *http.Request {
					req := httptest.NewRequest("GET", "/", nil)
					return req
				}(),
			},
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := QueryOptionalProtocol(tt.args.r); got != tt.want {
				t.Errorf("QueryOptionalProtocol() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQueryOptionalNetwork(t *testing.T) {
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name string
		args args
		want string
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
		},
		{
			"missing",
			args{
				func() *http.Request {
					req := httptest.NewRequest("GET", "/", nil)
					return req
				}(),
			},
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := QueryOptionalNetwork(tt.args.r); got != tt.want {
				t.Errorf("QueryOptionalNetwork() = %v, want %v", got, tt.want)
			}
		})
	}
}
