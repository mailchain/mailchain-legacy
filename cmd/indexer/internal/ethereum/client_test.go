package ethereum

import (
	"context"
	"math/big"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/assert"
)

func TestNewRPC(t *testing.T) {
	type args struct {
		address string
	}
	tests := []struct {
		name    string
		args    args
		wantNil bool
		wantErr bool
	}{
		{
			"success",
			args{
				"https://client.url",
			},
			false,
			false,
		},
		{
			"err",
			args{
				"/client.url",
			},
			true,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewRPC(tt.args.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRPC() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got == nil) != tt.wantNil {
				t.Errorf("NewRPC() = %v, want %v", got, tt.wantNil)
			}
		})
	}
}

func TestClient_Get(t *testing.T) {
	server := httptest.NewServer(nil)
	defer server.Close()
	type fields struct {
		client *ethclient.Client
	}
	type args struct {
		ctx     context.Context
		blockNo uint64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantBlk interface{}
		wantErr bool
	}{
		{
			"err",
			fields{
				func() *ethclient.Client {
					c, _ := ethclient.Dial(server.URL)
					return c
				}(),
			},
			args{
				context.Background(),
				1,
			},
			(*types.Block)(nil),
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				client: tt.fields.client,
			}
			gotBlk, err := c.Get(tt.args.ctx, tt.args.blockNo)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.wantBlk, gotBlk) {
				t.Errorf("Client.Get() = %v, want %v", gotBlk, tt.wantBlk)
			}
		})
	}
}

func TestClient_NetworkID(t *testing.T) {
	server := httptest.NewServer(nil)
	defer server.Close()
	type fields struct {
		client *ethclient.Client
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *big.Int
		wantErr bool
	}{
		{
			"err",
			fields{
				func() *ethclient.Client {
					c, _ := ethclient.Dial(server.URL)
					return c
				}(),
			},
			args{
				context.Background(),
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				client: tt.fields.client,
			}
			got, err := c.NetworkID(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.NetworkID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.NetworkID() = %v, want %v", got, tt.want)
			}
		})
	}
}
