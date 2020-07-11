package settings

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mailchain/go-substrate-rpc-client/types"
	"github.com/mailchain/mailchain/cmd/internal/settings/values"
	"github.com/mailchain/mailchain/cmd/internal/settings/values/valuestest"
	"github.com/mailchain/mailchain/internal/protocols/substrate"
	"github.com/stretchr/testify/assert"
)

func Test_substrateRPCSender(t *testing.T) {
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
					m.EXPECT().IsSet("senders.substrate-rpc-edgeware-berlin.address").Return(false)
					return m
				}(),
				substrate.EdgewareBerlin,
			},
			"ws://berlin1.edgewa.re:9944",
			"edgeware-berlin",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := substrateRPCSender(tt.args.s, tt.args.network)
			assert.Equal(t, tt.wantAddress, got.Address.Get())
			assert.Equal(t, tt.wantNetwork, got.network)
		})
	}
}

func TestSubstrateRPC_Produce(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("{\"result\":\"%v\"}", types.ExamplaryMetadataV11SubstrateString)))
	}))
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
				substrate.EdgewareBerlin,
			},
			false,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := SubstrateRPC{
				Address: tt.fields.Address,
				network: tt.fields.network,
			}
			got, err := s.Produce()
			if (err != nil) != tt.wantErr {
				t.Errorf("SubstrateRPC.Produce() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got != nil) == tt.wantErr {
				t.Errorf("SubstrateRPC.Produce() = %v, want %v", got == nil, tt.wantNil)
			}
		})
	}
}

func TestSubstrateRPC_Supports(t *testing.T) {
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
				substrate.EdgewareBerlin,
			},
			map[string]bool{"substrate/edgeware-berlin": true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := SubstrateRPC{
				Address: tt.fields.Address,
				network: tt.fields.network,
			}
			if got := s.Supports(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SubstrateRPC.Supports() = %v, want %v", got, tt.want)
			}
		})
	}
}
