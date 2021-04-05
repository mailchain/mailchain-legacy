package settings

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mailchain/mailchain/cmd/internal/settings/values"
	"github.com/mailchain/mailchain/cmd/internal/settings/values/valuestest"
	"github.com/mailchain/mailchain/internal/clients/blockscout"
	"github.com/mailchain/mailchain/internal/mailbox"
	"github.com/stretchr/testify/assert"
)

func Test_blockscoutBalanceFinderNoAuth(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type args struct {
		s values.Store
	}
	tests := []struct {
		name                        string
		args                        args
		wantEnabledProtocolNetworks []string
	}{
		{
			"success",
			args{
				func() values.Store {
					m := valuestest.NewMockStore(mockCtrl)
					m.EXPECT().IsSet("public-key-finders.blockscout-no-auth.enabled-networks").Return(false)
					return m
				}(),
			},
			[]string{"ethereum/mainnet"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := blockscoutPublicKeyFinderNoAuth(tt.args.s)
			assert.Equal(t, tt.wantEnabledProtocolNetworks, got.EnabledProtocolNetworks.Get())
		})
	}
}

func TestBlockscoutBalanceFinder_Supports(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type fields struct {
		EnabledProtocolNetworks values.StringSlice
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]bool
	}{
		{
			"success",
			fields{
				func() values.StringSlice {
					m := valuestest.NewMockStringSlice(mockCtrl)
					m.EXPECT().Get().Return([]string{"ethereum/mainnet", "ethereum/ropsten"})
					return m
				}(),
			},
			map[string]bool{"ethereum/mainnet": true, "ethereum/ropsten": true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := BlockscoutPublicKeyFinder{
				EnabledProtocolNetworks: tt.fields.EnabledProtocolNetworks,
			}
			if got := r.Supports(); !assert.Equal(t, tt.want, got) {
				t.Errorf("BlockscoutPublicKeyFinder.Supports() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBlockscoutBalanceFinder_Produce(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type fields struct {
		kind                    string
		EnabledProtocolNetworks values.StringSlice
	}
	tests := []struct {
		name    string
		fields  fields
		want    mailbox.PubKeyFinder
		wantErr bool
	}{
		{
			"success",
			fields{
				"test",
				func() values.StringSlice {
					m := valuestest.NewMockStringSlice(mockCtrl)
					return m
				}(),
			},
			&blockscout.APIClient{},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := BlockscoutPublicKeyFinder{
				kind:                    tt.fields.kind,
				EnabledProtocolNetworks: tt.fields.EnabledProtocolNetworks,
			}
			got, err := r.Produce()
			if (err != nil) != tt.wantErr {
				t.Errorf("BlockscoutBalanceFinder.Produce() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.IsType(t, tt.want, got) {
				t.Errorf("BlockscoutBalanceFinder.Produce() = %v, want %v", got, tt.want)
			}
		})
	}
}
