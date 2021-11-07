package settings

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mailchain/mailchain/cmd/internal/settings/values"
	"github.com/mailchain/mailchain/cmd/internal/settings/values/valuestest"
)

func Test_balanceFinders(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type args struct {
		s values.Store
	}
	tests := []struct {
		name    string
		args    args
		wantNil bool
	}{
		{
			"success",
			args{
				func() values.Store {
					m := valuestest.NewMockStore(mockCtrl)
					// m.EXPECT().IsSet("sentstore.kind")
					return m
				}(),
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := balanceFinders(tt.args.s)
			if (got == nil) != tt.wantNil {
				t.Errorf("balanceFinders() nil = %v, wantNil %v", got == nil, tt.wantNil)
				return
			}
		})
	}
}

func TestBalanceFinders_Produce(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type fields struct {
		clients map[string]BalanceFinderClient
	}
	type args struct {
		client string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantNil bool
		wantErr bool
	}{
		{
			"success",
			fields{
				map[string]BalanceFinderClient{
					"client": func() BalanceFinderClient {
						s := valuestest.NewMockStore(mockCtrl)
						s.EXPECT().IsSet("balance-finders.etherscan-no-auth.api-key").Return(false)
						return etherscanBalanceFinderNoAuth(s)
					}(),
				},
			},
			args{
				"client",
			},
			false,
			false,
		},
		{
			"err-nil-client",
			fields{
				nil,
			},
			args{
				"",
			},
			true,
			false,
		},
		{
			"err-no-client",
			fields{
				nil,
			},
			args{
				"invalid-client",
			},
			true,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := BalanceFinders{
				clients: tt.fields.clients,
			}
			got, err := s.Produce(tt.args.client)
			if (err != nil) != tt.wantErr {
				t.Errorf("BalanceFinders.Produce() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got == nil) != tt.wantNil {
				t.Errorf("BalanceFinders.Produce() nil = %v, wantNil %v", got == nil, tt.wantNil)
				return
			}
		})
	}
}
