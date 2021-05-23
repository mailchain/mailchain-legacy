package settings

import (
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mailchain/mailchain/cmd/internal/settings/values"
	"github.com/mailchain/mailchain/cmd/internal/settings/values/valuestest"
	"github.com/mailchain/mailchain/internal/protocols/algorand"
	"github.com/stretchr/testify/assert"
)

func Test_algodSender(t *testing.T) {
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
					m.EXPECT().IsSet("senders.algod-testnet.address").Return(false)
					return m
				}(),
				algorand.Testnet,
			},
			"https://api.testnet.algoexplorer.io",
			"testnet",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := algodSender(tt.args.s, tt.args.network)
			assert.Equal(t, tt.wantAddress, got.Address.Get())
			assert.Equal(t, tt.wantNetwork, got.network)
		})
	}
}

func TestAlgod_Produce(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type fields struct {
		Address values.String
		Token   values.String
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
					m.EXPECT().Get().Return("anything")
					return m
				}(),
				func() values.String {
					m := valuestest.NewMockString(mockCtrl)
					m.EXPECT().Get().Return("")
					return m
				}(),
				algorand.Testnet,
			},
			false,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := AlgodSender{
				Address: tt.fields.Address,
				Token:   tt.fields.Token,
				network: tt.fields.network,
			}
			got, err := s.Produce()
			if (err != nil) != tt.wantErr {
				t.Errorf("AlgodSender.Produce() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got != nil) == tt.wantErr {
				t.Errorf("AlgodSender.Produce() = %v, want %v", got == nil, tt.wantNil)
			}
		})
	}
}

func TestAlgod_Supports(t *testing.T) {
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
				algorand.Testnet,
			},
			map[string]bool{"algorand/testnet": true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := AlgodSender{
				Address: tt.fields.Address,
				network: tt.fields.network,
			}
			if got := s.Supports(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AlgodSender.Supports() = %v, want %v", got, tt.want)
			}
		})
	}
}
