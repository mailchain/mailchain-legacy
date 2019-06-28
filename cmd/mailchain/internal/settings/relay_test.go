package settings

import (
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/values"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/values/valuestest"
	"github.com/stretchr/testify/assert"
)

func Test_relaySender(t *testing.T) {
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type args struct {
		s       values.Store
		network string
	}
	tests := []struct {
		name                        string
		args                        args
		wantEnabledProtocolNetworks []string
		wantBaseURL                 string
	}{
		{
			"success",
			args{
				func() values.Store {
					m := valuestest.NewMockStore(mockCtrl)
					m.EXPECT().IsSet("senders.ethereum-relay.enabled-networks").Return(false)
					m.EXPECT().IsSet("senders.ethereum-relay.base-url").Return(false)
					return m
				}(),
				"ethereum",
			},
			[]string{"ethereum/goerli", "ethereum/kovan", "ethereum/mainnet", "ethereum/rinkeby", "ethereum/ropsten"},
			"https://relay.mailchain.xyz/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := relaySender(tt.args.s, tt.args.network)
			assert.Equal(tt.wantEnabledProtocolNetworks, got.EnabledProtocolNetworks.Get())
			assert.Equal(tt.wantBaseURL, got.BaseURL.Get())
		})
	}
}

func TestRelaySender_Supports(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type fields struct {
		EnabledProtocolNetworks values.StringSlice
		BaseURL                 values.String
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
			r := RelaySender{
				EnabledProtocolNetworks: tt.fields.EnabledProtocolNetworks,
				BaseURL:                 tt.fields.BaseURL,
			}
			if got := r.Supports(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RelaySender.Supports() = %v, want %v", got, tt.want)
			}
		})
	}
}
