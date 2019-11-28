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

package stores

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getHTTPMessage(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name    string
		server  *httptest.Server
		want    []byte
		wantErr bool
	}{
		{
			"success",
			httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte("contents"))
				}),
			),
			[]byte{0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x73},
			false,
		},
		{
			"not-found",
			httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusNotFound)
				}),
			),
			nil,
			true,
		},
		{
			"error",
			func() *httptest.Server {
				s := httptest.NewServer(
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(http.StatusNotFound)
					}),
				)
				s.URL = "http://somethignnotvalid:133443"
				return s
			}(),
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer tt.server.Close()
			got, err := getHTTPMessage(tt.server.URL)
			if (err != nil) != tt.wantErr {
				t.Errorf("getHTTPMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(tt.want, got) {
				t.Errorf("getHTTPMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getAnyMessage(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		location string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			"err-parse-location",
			args{
				"un-parseable-location" + string(0x7f),
			},
			nil,
			true,
		},
		{
			"http",
			args{
				func() string {
					s := httptest.NewServer(
						http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							w.Write([]byte("contents"))
						}),
					)
					return s.URL
				}(),
			},
			[]byte{0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x73},
			false,
		},
		{
			"https",
			args{
				func() string {
					s := httptest.NewServer(
						http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							w.Write([]byte("contents"))
						}),
					)
					return s.URL
				}(),
			},
			[]byte{0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x73},
			false,
		},
		{
			"test",
			args{
				"test://hostvalue",
			},
			[]byte{0x68, 0x6f, 0x73, 0x74, 0x76, 0x61, 0x6c, 0x75, 0x65},
			false,
		},
		{
			"file",
			args{
				"file://./testdata/contents.txt",
			},
			[]byte{0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x73, 0xa},
			false,
		},
		{
			"err-not-supported",
			args{
				"no-sup://location",
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getAnyMessage(tt.args.location)
			if (err != nil) != tt.wantErr {
				t.Errorf("getAnyMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(tt.want, got) {
				t.Errorf("getAnyMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetMessage(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		location      string
		integrityHash []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			"success-hash",
			args{
				"file://./testdata/simple.golden.eml-2204ec872b32",
				[]byte{0x22, 0x04, 0xec, 0x87, 0x2b, 0x32},
			},
			func() []byte {
				contents, _ := ioutil.ReadFile("./testdata/simple.golden.eml-2204ec872b32")
				return contents
			}(),
			false,
		},
		{
			"success-no-hash",
			args{
				"file://./testdata/simple.golden.eml-2204ec872b32",
				nil,
			},
			func() []byte {
				contents, _ := ioutil.ReadFile("./testdata/simple.golden.eml-2204ec872b32")
				return contents
			}(),
			false,
		},
		{
			"err-no-schema",
			args{
				"invalid://./testdata/simple.golden.eml-2204ec872b32",
				nil,
			},
			nil,
			true,
		},
		{
			"hash-part-does-not-match",
			args{
				"test://hash.does.not.match-2204f3d89e5a",
				[]byte{0x22, 0x04, 0xf3, 0xd8, 0x9e, 0x5a},
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetMessage(tt.args.location, tt.args.integrityHash)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(tt.want, got) {
				t.Errorf("GetMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}
