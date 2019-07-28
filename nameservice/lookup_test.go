// Copyright 2019 Finobo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package nameservice

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestNewLookupService(t *testing.T) {
	type args struct {
		baseURL string
	}
	tests := []struct {
		name    string
		args    args
		wantNil bool
	}{
		{
			"success",
			args{
				"https://client.url",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewLookupService(tt.args.baseURL)
			if (got == nil) != tt.wantNil {
				t.Errorf("NewLookupService() = %v, want %v", got, tt.wantNil)
			}
		})
	}
}

func TestLookupService_ResolveName(t *testing.T) {
	assert := assert.New(t)
	type fields struct {
		newRequest func(method string, url string, body io.Reader) (*http.Request, error)
		doRequest  func(req *http.Request) (*http.Response, error)
	}
	type args struct {
		ctx        context.Context
		protocol   string
		network    string
		domainName string
	}
	tests := []struct {
		name    string
		server  *httptest.Server
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		{
			"success",
			httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte("{\"address\":\"0x4ad2b251246aafc2f3bdf3b690de3bf906622c51\"}"))
				}),
			),
			fields{
				http.NewRequest,
				func() func(req *http.Request) (*http.Response, error) {
					c := http.Client{
						Timeout: 1 * time.Second,
					}
					return c.Do
				}(),
			},
			args{
				context.Background(),
				"ethereum",
				"mainnet",
				"test.eth",
			},
			[]byte{0x4a, 0xd2, 0xb2, 0x51, 0x24, 0x6a, 0xaf, 0xc2, 0xf3, 0xbd, 0xf3, 0xb6, 0x90, 0xde, 0x3b, 0xf9, 0x6, 0x62, 0x2c, 0x51},
			false,
		},
		{
			"err-invalid-200",
			httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte("{\"address:\"0x4ad2b251246aafc2f3bdf3b690de3bf906622c51\"}"))
				}),
			),
			fields{
				http.NewRequest,
				func() func(req *http.Request) (*http.Response, error) {
					c := http.Client{
						Timeout: 1 * time.Second,
					}
					return c.Do
				}(),
			},
			args{
				context.Background(),
				"ethereum",
				"mainnet",
				"test.eth",
			},
			nil,
			true,
		},
		{
			"err-404",
			httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusNotFound)
					w.Write([]byte("{\"code\":404,\"message\":\"not found: unregistered name\"}"))
				}),
			),
			fields{
				http.NewRequest,
				func() func(req *http.Request) (*http.Response, error) {
					c := http.Client{
						Timeout: 1 * time.Second,
					}
					return c.Do
				}(),
			},
			args{
				context.Background(),
				"ethereum",
				"mainnet",
				"test.eth",
			},
			nil,
			true,
		},
		{
			"err-invalid-404",
			httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusNotFound)
					w.Write([]byte("{\"code\":404,\"message:\"not found: unregistered name\"}"))
				}),
			),
			fields{
				http.NewRequest,
				func() func(req *http.Request) (*http.Response, error) {
					c := http.Client{
						Timeout: 1 * time.Second,
					}
					return c.Do
				}(),
			},
			args{
				context.Background(),
				"ethereum",
				"mainnet",
				"test.eth",
			},
			nil,
			true,
		},
		{
			"err-do-request",
			httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusNotFound)
					w.Write([]byte("{\"code\":404,\"message:\"not found: unregistered name\"}"))
				}),
			),
			fields{
				http.NewRequest,
				func() func(req *http.Request) (*http.Response, error) {
					return func(req *http.Request) (*http.Response, error) {
						return nil, errors.Errorf("do request failed")
					}
				}(),
			},
			args{
				context.Background(),
				"ethereum",
				"mainnet",
				"test.eth",
			},
			nil,
			true,
		},
		{
			"err-do-request",
			httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusNotFound)
					w.Write([]byte("{\"code\":404,\"message:\"not found: unregistered name\"}"))
				}),
			),
			fields{
				func() func(method string, url string, body io.Reader) (*http.Request, error) {
					return func(method string, url string, body io.Reader) (*http.Request, error) {
						return nil, errors.Errorf("failed to create request")
					}
				}(),
				nil,
			},
			args{
				context.Background(),
				"ethereum",
				"mainnet",
				"test.eth",
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := LookupService{
				baseURL:    tt.server.URL,
				newRequest: tt.fields.newRequest,
				doRequest:  tt.fields.doRequest,
			}
			got, err := s.ResolveName(tt.args.ctx, tt.args.protocol, tt.args.network, tt.args.domainName)
			if (err != nil) != tt.wantErr {
				t.Errorf("LookupService.ResolveName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(tt.want, got) {
				t.Errorf("LookupService.ResolveName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLookupService_ResolveAddress(t *testing.T) {
	type fields struct {
		newRequest func(method string, url string, body io.Reader) (*http.Request, error)
		doRequest  func(req *http.Request) (*http.Response, error)
	}
	type args struct {
		ctx      context.Context
		protocol string
		network  string
		address  []byte
	}
	tests := []struct {
		name    string
		server  *httptest.Server
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			"success",
			httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte("{\"name\":\"test.eth\"}"))
				}),
			),
			fields{
				http.NewRequest,
				func() func(req *http.Request) (*http.Response, error) {
					c := http.Client{
						Timeout: 1 * time.Second,
					}
					return c.Do
				}(),
			},
			args{
				context.Background(),
				"ethereum",
				"mainnet",
				[]byte{0x4a, 0xd2, 0xb2, 0x51, 0x24, 0x6a, 0xaf, 0xc2, 0xf3, 0xbd, 0xf3, 0xb6, 0x90, 0xde, 0x3b, 0xf9, 0x6, 0x62, 0x2c, 0x51},
			},
			"test.eth",
			false,
		},
		{
			"err-invalid-200",
			httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte("{\"address:\"0x4ad2b251246aafc2f3bdf3b690de3bf906622c51\"}"))
				}),
			),
			fields{
				http.NewRequest,
				func() func(req *http.Request) (*http.Response, error) {
					c := http.Client{
						Timeout: 1 * time.Second,
					}
					return c.Do
				}(),
			},
			args{
				context.Background(),
				"ethereum",
				"mainnet",
				[]byte{0x4a, 0xd2, 0xb2, 0x51, 0x24, 0x6a, 0xaf, 0xc2, 0xf3, 0xbd, 0xf3, 0xb6, 0x90, 0xde, 0x3b, 0xf9, 0x6, 0x62, 0x2c, 0x51},
			},
			"",
			true,
		},
		{
			"err-404",
			httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusNotFound)
					w.Write([]byte("{\"code\":404,\"message\":\"not found: unregistered name\"}"))
				}),
			),
			fields{
				http.NewRequest,
				func() func(req *http.Request) (*http.Response, error) {
					c := http.Client{
						Timeout: 1 * time.Second,
					}
					return c.Do
				}(),
			},
			args{
				context.Background(),
				"ethereum",
				"mainnet",
				[]byte{0x4a, 0xd2, 0xb2, 0x51, 0x24, 0x6a, 0xaf, 0xc2, 0xf3, 0xbd, 0xf3, 0xb6, 0x90, 0xde, 0x3b, 0xf9, 0x6, 0x62, 0x2c, 0x51},
			},
			"",
			true,
		},
		{
			"err-invalid-404",
			httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusNotFound)
					w.Write([]byte("{\"code\":404,\"message:\"not found: unregistered name\"}"))
				}),
			),
			fields{
				http.NewRequest,
				func() func(req *http.Request) (*http.Response, error) {
					c := http.Client{
						Timeout: 1 * time.Second,
					}
					return c.Do
				}(),
			},
			args{
				context.Background(),
				"ethereum",
				"mainnet",
				[]byte{0x4a, 0xd2, 0xb2, 0x51, 0x24, 0x6a, 0xaf, 0xc2, 0xf3, 0xbd, 0xf3, 0xb6, 0x90, 0xde, 0x3b, 0xf9, 0x6, 0x62, 0x2c, 0x51},
			},
			"",
			true,
		},
		{
			"err-do-request",
			httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusNotFound)
					w.Write([]byte("{\"code\":404,\"message:\"not found: unregistered name\"}"))
				}),
			),
			fields{
				http.NewRequest,
				func() func(req *http.Request) (*http.Response, error) {
					return func(req *http.Request) (*http.Response, error) {
						return nil, errors.Errorf("do request failed")
					}
				}(),
			},
			args{
				context.Background(),
				"ethereum",
				"mainnet",
				[]byte{0x4a, 0xd2, 0xb2, 0x51, 0x24, 0x6a, 0xaf, 0xc2, 0xf3, 0xbd, 0xf3, 0xb6, 0x90, 0xde, 0x3b, 0xf9, 0x6, 0x62, 0x2c, 0x51},
			},
			"",
			true,
		},
		{
			"err-do-request",
			httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusNotFound)
					w.Write([]byte("{\"code\":404,\"message:\"not found: unregistered name\"}"))
				}),
			),
			fields{
				func() func(method string, url string, body io.Reader) (*http.Request, error) {
					return func(method string, url string, body io.Reader) (*http.Request, error) {
						return nil, errors.Errorf("failed to create request")
					}
				}(),
				nil,
			},
			args{
				context.Background(),
				"ethereum",
				"mainnet",
				[]byte{0x4a, 0xd2, 0xb2, 0x51, 0x24, 0x6a, 0xaf, 0xc2, 0xf3, 0xbd, 0xf3, 0xb6, 0x90, 0xde, 0x3b, 0xf9, 0x6, 0x62, 0x2c, 0x51},
			},
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := LookupService{
				baseURL:    tt.server.URL,
				newRequest: tt.fields.newRequest,
				doRequest:  tt.fields.doRequest,
			}
			got, err := s.ResolveAddress(tt.args.ctx, tt.args.protocol, tt.args.network, tt.args.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("LookupService.ResolveAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("LookupService.ResolveAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}
