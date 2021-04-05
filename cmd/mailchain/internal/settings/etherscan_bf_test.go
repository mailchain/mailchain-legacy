package settings

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mailchain/mailchain/cmd/internal/settings/values"
	"github.com/mailchain/mailchain/cmd/internal/settings/values/valuestest"
	"github.com/stretchr/testify/assert"
)

func Test_etherscanBalanceFinderAny(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type args struct {
		s    values.Store
		kind string
	}
	tests := []struct {
		name                        string
		args                        args
		wantEnabledProtocolNetworks []string
		wantAPIKey                  string
	}{
		{
			"success",
			args{
				func() values.Store {
					m := valuestest.NewMockStore(mockCtrl)
					m.EXPECT().IsSet("public-key-finders.type.enabled-networks").Return(false)
					m.EXPECT().IsSet("public-key-finders.type.api-key").Return(false)
					return m
				}(),
				"type",
			},
			[]string{"ethereum/goerli", "ethereum/kovan", "ethereum/mainnet", "ethereum/rinkeby", "ethereum/ropsten"},
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := etherscanPublicKeyFinderAny(tt.args.s, tt.args.kind)
			assert.Equal(t, tt.wantEnabledProtocolNetworks, got.EnabledProtocolNetworks.Get())
			assert.Equal(t, tt.wantAPIKey, got.APIKey.Get())
		})
	}
}

func Test_etherscanBalanceFinder(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type args struct {
		s values.Store
	}
	tests := []struct {
		name                        string
		args                        args
		wantEnabledProtocolNetworks []string
		wantAPIKey                  string
	}{
		{
			"success",
			args{
				func() values.Store {
					m := valuestest.NewMockStore(mockCtrl)
					m.EXPECT().IsSet("public-key-finders.etherscan.enabled-networks").Return(false)
					m.EXPECT().IsSet("public-key-finders.etherscan.api-key").Return(false)
					return m
				}(),
			},
			[]string{"ethereum/goerli", "ethereum/kovan", "ethereum/mainnet", "ethereum/rinkeby", "ethereum/ropsten"},
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := etherscanPublicKeyFinder(tt.args.s)
			assert.Equal(t, tt.wantEnabledProtocolNetworks, got.EnabledProtocolNetworks.Get())
			assert.Equal(t, tt.wantAPIKey, got.APIKey.Get())
		})
	}
}

func Test_etherscanBalanceFinderNoAuth(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type args struct {
		s values.Store
	}
	tests := []struct {
		name                        string
		args                        args
		wantEnabledProtocolNetworks []string
		wantAPIKey                  string
	}{
		{
			"success",
			args{
				func() values.Store {
					m := valuestest.NewMockStore(mockCtrl)
					m.EXPECT().IsSet("public-key-finders.etherscan-no-auth.enabled-networks").Return(false)
					m.EXPECT().IsSet("public-key-finders.etherscan-no-auth.api-key").Return(false)
					return m
				}(),
			},
			[]string{"ethereum/goerli", "ethereum/kovan", "ethereum/mainnet", "ethereum/rinkeby", "ethereum/ropsten"},
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := etherscanPublicKeyFinderNoAuth(tt.args.s)
			assert.Equal(t, tt.wantEnabledProtocolNetworks, got.EnabledProtocolNetworks.Get())
			assert.Equal(t, tt.wantAPIKey, got.APIKey.Get())
		})
	}
}

func Test_etherscanBalanceFinder_Supports(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type fields struct {
		EnabledProtocolNetworks values.StringSlice
		APIKey                  values.String
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
				nil,
			},
			map[string]bool{"ethereum/mainnet": true, "ethereum/ropsten": true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := EtherscanPublicKeyFinder{
				EnabledProtocolNetworks: tt.fields.EnabledProtocolNetworks,
				APIKey:                  tt.fields.APIKey,
			}
			if got := r.Supports(); !assert.Equal(t, tt.want, got) {
				t.Errorf("EtherscanPublicKeyFinder.Supports() = %v, want %v", got, tt.want)
			}
		})
	}
}
