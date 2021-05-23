package settings

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mailchain/mailchain/cmd/internal/settings/values"
	"github.com/mailchain/mailchain/cmd/internal/settings/values/valuestest"
	"github.com/mailchain/mailchain/internal/mailbox"
	"github.com/mailchain/mailchain/internal/protocols/algorand"
	"github.com/stretchr/testify/assert"
)

func Test_algorandPublicKeyFinder(t *testing.T) {
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
					m.EXPECT().IsSet("public-key-finders.algorand-public-key-extractor.enabled-networks").Return(false)
					return m
				}(),
			},
			[]string{"algorand/mainnet", "algorand/betanet", "algorand/testnet"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := algorandPublicKeyFinder(tt.args.s)
			assert.Equal(t, tt.wantEnabledProtocolNetworks, got.EnabledProtocolNetworks.Get())
		})
	}
}

func TestAlgorandPublicKeyFinder_Supports(t *testing.T) {
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
					m.EXPECT().Get().Return([]string{"algorand/mainnet", "algorand/betanet"})
					return m
				}(),
			},
			map[string]bool{"algorand/mainnet": true, "algorand/betanet": true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := AlgorandPublicKeyFinder{
				EnabledProtocolNetworks: tt.fields.EnabledProtocolNetworks,
			}
			if got := r.Supports(); !assert.Equal(t, tt.want, got) {
				t.Errorf("AlgorandPublicKeyFinder.Supports() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAlgorandPublicKeyFinder_Produce(t *testing.T) {
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
			&algorand.PublicKeyFinder{},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := AlgorandPublicKeyFinder{
				kind:                    tt.fields.kind,
				EnabledProtocolNetworks: tt.fields.EnabledProtocolNetworks,
			}
			got, err := r.Produce()
			if (err != nil) != tt.wantErr {
				t.Errorf("AlgorandPublicKeyFinder.Produce() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.IsType(t, tt.want, got) {
				t.Errorf("AlgorandPublicKeyFinder.Produce() = %v, want %v", got, tt.want)
			}
		})
	}
}
