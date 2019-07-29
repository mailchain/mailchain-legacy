package settings

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/values"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/values/valuestest"
	"github.com/mailchain/mailchain/nameservice"
	"github.com/stretchr/testify/assert"
)

func Test_DomainNameServices(t *testing.T) {
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
			got := domainNameServices(tt.args.s)
			if (got == nil) != tt.wantNil {
				t.Errorf("domainNameServices() nil = %v, wantNil %v", got == nil, tt.wantNil)
				return
			}
		})
	}
}

func TestDomainNameServices_Produce(t *testing.T) {
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type fields struct {
		clients map[string]NameServiceDomainClient
	}
	type args struct {
		client string
	}

	tests := []struct {
		name     string
		fields   fields
		args     args
		wantType nameservice.ForwardLookup
		wantErr  bool
	}{
		{
			"success",
			fields{
				map[string]NameServiceDomainClient{
					"client": func() NameServiceDomainClient {
						s := valuestest.NewMockStore(mockCtrl)
						s.EXPECT().IsSet("name-service-domain-name.base-url").Return(false)
						return mailchainDomainNameServices(s)
					}(),
				},
			},
			args{
				"client",
			},
			&nameservice.LookupService{},
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
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := DomainNameServices{
				clients: tt.fields.clients,
			}
			got, err := s.Produce(tt.args.client)
			if (err != nil) != tt.wantErr {
				t.Errorf("DomainNameServices.Produce() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.IsType(tt.wantType, got) {
				t.Errorf("DomainNameServices.Produce() = %v, want %v", got, tt.wantType)
			}
		})
	}
}

func TestMailchainDomainNameServices_Supports(t *testing.T) {
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type fields struct {
		BaseURL                 values.String
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
				nil,
				func() values.StringSlice {
					m := valuestest.NewMockStringSlice(mockCtrl)
					m.EXPECT().Get().Return([]string{"ethereum/mainnet", "ethereum/goerli"})
					return m
				}(),
			},
			map[string]bool{"ethereum/goerli": true, "ethereum/mainnet": true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := MailchainDomainNameServices{
				BaseURL:                 tt.fields.BaseURL,
				EnabledProtocolNetworks: tt.fields.EnabledProtocolNetworks,
			}
			if got := s.Supports(); !assert.Equal(tt.want, got) {
				t.Errorf("MailchainDomainNameServices.Supports() = %v, want %v", got, tt.want)
			}
		})
	}
}
