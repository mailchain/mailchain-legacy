package settings

import (
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mailchain/mailchain/cmd/internal/settings/values"
	"github.com/mailchain/mailchain/cmd/internal/settings/values/valuestest"
	"github.com/mailchain/mailchain/internal/clients/blockscout"
	"github.com/mailchain/mailchain/internal/mailbox"
	"github.com/stretchr/testify/assert"
)

func Test_blockscoutReceiverAny(t *testing.T) {
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
					m.EXPECT().IsSet("receivers.type.enabled-networks").Return(false)
					m.EXPECT().IsSet("receivers.type.api-key").Return(false)
					return m
				}(),
				"type",
			},
			[]string{"ethereum/mainnet"},
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := blockscoutReceiverAny(tt.args.s, tt.args.kind)
			assert.Equal(t, tt.wantEnabledProtocolNetworks, got.EnabledProtocolNetworks.Get())
			assert.Equal(t, tt.wantAPIKey, got.APIKey.Get())
		})
	}
}

func Test_blockscoutReceiverNoAuth(t *testing.T) {
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
					m.EXPECT().IsSet("receivers.blockscout-no-auth.enabled-networks").Return(false)
					m.EXPECT().IsSet("receivers.blockscout-no-auth.api-key").Return(false)
					return m
				}(),
			},
			[]string{"ethereum/mainnet"},
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := blockscoutReceiverNoAuth(tt.args.s)
			assert.Equal(t, tt.wantEnabledProtocolNetworks, got.EnabledProtocolNetworks.Get())
			assert.Equal(t, tt.wantAPIKey, got.APIKey.Get())
		})
	}
}

func TestBlockscoutReceiver_Supports(t *testing.T) {
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
			r := BlockscoutReceiver{
				EnabledProtocolNetworks: tt.fields.EnabledProtocolNetworks,
				APIKey:                  tt.fields.APIKey,
			}
			if got := r.Supports(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("blockscoutReceiver.Supports() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBlockscoutReceiver_Produce(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type fields struct {
		kind                    string
		EnabledProtocolNetworks values.StringSlice
		APIKey                  values.String
	}
	tests := []struct {
		name    string
		fields  fields
		want    mailbox.Receiver
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
				func() values.String {
					m := valuestest.NewMockString(mockCtrl)
					m.EXPECT().Get().Return("apikey")
					return m
				}(),
			},
			&blockscout.APIClient{},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := BlockscoutReceiver{
				kind:                    tt.fields.kind,
				EnabledProtocolNetworks: tt.fields.EnabledProtocolNetworks,
				APIKey:                  tt.fields.APIKey,
			}
			got, err := r.Produce()
			if (err != nil) != tt.wantErr {
				t.Errorf("BlockscoutReceiver.Produce() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.IsType(t, tt.want, got) {
				t.Errorf("BlockscoutReceiver.Produce() = %v, want %v", got, tt.want)
			}
		})
	}
}
