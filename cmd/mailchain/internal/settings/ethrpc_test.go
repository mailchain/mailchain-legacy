package settings

import (
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mailchain/mailchain/cmd/internal/settings/values"
	"github.com/mailchain/mailchain/cmd/internal/settings/values/valuestest"
	"github.com/stretchr/testify/assert"
)

func TestEthereumRPC2_Produce(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	server := httptest.NewServer(nil)
	defer server.Close()
	type fields struct {
		Address values.String
		network string
	}
	tests := []struct {
		name    string
		fields  fields
		wantNil bool
		wantErr bool
	}{
		{
			"success",
			fields{
				func() values.String {
					m := valuestest.NewMockString(mockCtrl)
					m.EXPECT().Get().Return(server.URL)
					return m
				}(),
				"mainnet",
			},
			false,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := EthereumRPC2{
				Address: tt.fields.Address,
				network: tt.fields.network,
			}
			got, err := e.Produce()
			if (err != nil) != tt.wantErr {
				t.Errorf("EthereumRPC2.Produce() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got != nil) == tt.wantErr {
				t.Errorf("EthereumRPC2.Produce() = %v, want %v", got == nil, tt.wantNil)
			}
		})
	}
}

func TestEthereumRPC2_Supports(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type fields struct {
		Address values.String
		network string
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]bool
	}{
		{
			"success",
			fields{
				func() values.String {
					m := valuestest.NewMockString(mockCtrl)
					return m
				}(),
				"mainnet",
			},
			map[string]bool{"ethereum/mainnet": true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := EthereumRPC2{
				Address: tt.fields.Address,
				network: tt.fields.network,
			}
			if got := e.Supports(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EthereumRPC2.Supports() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ethereumRPC2Sender(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type args struct {
		s       values.Store
		network string
	}
	tests := []struct {
		name        string
		args        args
		wantAddress string
		wantNetwork string
	}{
		{
			"success",
			args{
				func() values.Store {
					m := valuestest.NewMockStore(mockCtrl)
					m.EXPECT().IsSet("senders.ethereum-rpc2-mainnet.address").Return(false)
					return m
				}(),
				"mainnet",
			},
			"https://relay.mailchain.xyz/json-rpc/ethereum/mainnet",
			"mainnet",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ethereumRPC2Sender(tt.args.s, tt.args.network)
			assert.Equal(t, tt.wantAddress, got.Address.Get())
			assert.Equal(t, tt.wantNetwork, got.network)
		})
	}
}
