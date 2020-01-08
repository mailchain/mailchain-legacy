package substraterpc

import (
	"fmt"
	gsrpc "github.com/centrifuge/go-substrate-rpc-client"
	"github.com/centrifuge/go-substrate-rpc-client/types"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_GetMetadata(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		server *httptest.Server
		hash   types.Hash
	}
	tests := []struct {
		name    string
		args    args
		want    *types.Metadata
		wantErr bool
	}{
		{
			"success-latest",
			args{
				func() *httptest.Server {
					s := httptest.NewServer(
						http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							w.Write([]byte(fmt.Sprintf("{\"result\":\"%v\"}", types.ExamplaryMetadataV4String)))
						}),
					)
					return s
				}(),
				types.Hash{},
			},
			types.ExamplaryMetadataV4,
			false,
		},
		{
			"success-specific",
			args{
				func() *httptest.Server {
					s := httptest.NewServer(
						http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							w.Write([]byte(fmt.Sprintf("{\"result\":\"%v\"}", types.ExamplaryMetadataV4String)))
						}),
					)
					return s
				}(),
				types.NewHash([]byte("test")),
			},
			types.ExamplaryMetadataV4,
			false,
		},
		{
			"error-latest",
			args{
				func() *httptest.Server {
					s := httptest.NewServer(
						http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							w.WriteHeader(http.StatusBadRequest)
						}),
					)
					return s
				}(),
				types.Hash{},
			},
			nil,
			true,
		},
		{
			"error-specific",
			args{
				func() *httptest.Server {
					s := httptest.NewServer(
						http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							w.WriteHeader(http.StatusBadRequest)
						}),
					)
					return s
				}(),
				types.NewHash([]byte("test")),
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api, _ := gsrpc.NewSubstrateAPI(tt.args.server.URL)
			client := NewClient(api)
			got, err := client.GetMetadata(tt.args.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetMetadata() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(tt.want, got) {
				t.Errorf("Client.GetMetadata() got = %v, want %v", got, tt.want)
			}
		})
	}
}
