package settings

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/values"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/values/valuestest"
	"github.com/stretchr/testify/assert"
)

func Test_network(t *testing.T) {
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type args struct {
		s        values.Store
		protocol string
		network  string
	}
	tests := []struct {
		name                      string
		args                      args
		wantNameServiceAddress    string
		wantNameServiceDomainName string
		wantPublicKeyFinder       string
		wantReceiver              string
		wantSender                string
		wantDisabled              bool
	}{
		{
			"success",
			args{
				func() values.Store {
					m := valuestest.NewMockStore(mockCtrl)
					m.EXPECT().IsSet("protocols.ethereum.networks.mainnet.name-service-address").Return(false)
					m.EXPECT().IsSet("protocols.ethereum.networks.mainnet.name-service-domain-name").Return(false)
					m.EXPECT().IsSet("protocols.ethereum.networks.mainnet.public-key-finder").Return(false)
					m.EXPECT().IsSet("protocols.ethereum.networks.mainnet.sender").Return(false)
					m.EXPECT().IsSet("protocols.ethereum.networks.mainnet.receiver").Return(false)
					m.EXPECT().IsSet("protocols.ethereum.networks.mainnet.disabled").Return(false)
					return m
				}(),
				"ethereum",
				"mainnet",
			},
			"mailchain",
			"mailchain",
			"etherscan-no-auth",
			"etherscan-no-auth",
			"ethereum-relay",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := network(tt.args.s, tt.args.protocol, tt.args.network)
			assert.Equal(tt.wantNameServiceAddress, got.NameServiceAddress.Get())
			assert.Equal(tt.wantNameServiceDomainName, got.NameServiceDomainName.Get())
			assert.Equal(tt.wantPublicKeyFinder, got.PublicKeyFinder.Get())
			assert.Equal(tt.wantReceiver, got.Receiver.Get())
			assert.Equal(tt.wantSender, got.Sender.Get())
			assert.Equal(tt.wantDisabled, got.Disabled.Get())
		})
	}
}

func TestNetwork_ProduceSender(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type fields struct {
		PublicKeyFinder values.String
		Receiver        values.String
		Sender          values.String
		Disabled        values.Bool
	}
	type args struct {
		senders *Senders
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
				nil,
				nil,
				func() values.String {
					m := valuestest.NewMockString(mockCtrl)
					m.EXPECT().Get().Return("ethereum-relay")
					return m
				}(),
				nil,
			},
			args{
				senders(func() values.Store {
					m := valuestest.NewMockStore(mockCtrl)
					m.EXPECT().IsSet("senders.ethereum-relay.base-url").Return(false)
					return m
				}()),
			},
			false,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Network{
				PublicKeyFinder: tt.fields.PublicKeyFinder,
				Receiver:        tt.fields.Receiver,
				Sender:          tt.fields.Sender,
			}
			got, err := s.ProduceSender(tt.args.senders)
			if (err != nil) != tt.wantErr {
				t.Errorf("Network.ProduceSender() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got == nil) != tt.wantNil {
				t.Errorf("Network.ProduceSender() nil = %v, wantErr %v", got == nil, tt.wantNil)
				return
			}
		})
	}
}

func TestNetwork_ProduceReceiver(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type fields struct {
		PublicKeyFinder values.String
		Receiver        values.String
		Sender          values.String
	}
	type args struct {
		receivers *Receivers
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
				nil,
				func() values.String {
					m := valuestest.NewMockString(mockCtrl)
					m.EXPECT().Get().Return("etherscan-no-auth")
					return m
				}(),
				nil,
			},
			args{
				receivers(func() values.Store {
					m := valuestest.NewMockStore(mockCtrl)
					m.EXPECT().IsSet("receivers.etherscan-no-auth.api-key").Return(false)
					return m
				}()),
			},
			false,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Network{
				PublicKeyFinder: tt.fields.PublicKeyFinder,
				Receiver:        tt.fields.Receiver,
				Sender:          tt.fields.Sender,
			}
			got, err := s.ProduceReceiver(tt.args.receivers)
			if (err != nil) != tt.wantErr {
				t.Errorf("Network.ProduceReceiver() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got == nil) != tt.wantNil {
				t.Errorf("Network.ProduceReceiver() nil = %v, wantErr %v", got == nil, tt.wantNil)
				return
			}
		})
	}
}

func TestNetwork_ProducePublicKeyFinders(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type fields struct {
		PublicKeyFinder values.String
		Receiver        values.String
		Sender          values.String
	}
	type args struct {
		publicKeyFinders *PublicKeyFinders
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
				func() values.String {
					m := valuestest.NewMockString(mockCtrl)
					m.EXPECT().Get().Return("etherscan-no-auth")
					return m
				}(),
				nil,
				nil,
			},
			args{
				publicKeyFinders(func() values.Store {
					m := valuestest.NewMockStore(mockCtrl)
					m.EXPECT().IsSet("public-key-finders.etherscan-no-auth.api-key").Return(false)
					return m
				}()),
			},
			false,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Network{
				PublicKeyFinder: tt.fields.PublicKeyFinder,
				Receiver:        tt.fields.Receiver,
				Sender:          tt.fields.Sender,
			}
			got, err := s.ProducePublicKeyFinders(tt.args.publicKeyFinders)
			if (err != nil) != tt.wantErr {
				t.Errorf("Network.ProducePublicKeyFinders() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got == nil) != tt.wantNil {
				t.Errorf("Network.ProducePublicKeyFinders() nil = %v, wantErr %v", got == nil, tt.wantNil)
				return
			}
		})
	}
}

func TestNetwork_ProduceNameServiceAddress(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type fields struct {
		NameServiceAddress values.String
	}
	type args struct {
		ans *AddressNameServices
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
				func() values.String {
					m := valuestest.NewMockString(mockCtrl)
					m.EXPECT().Get().Return("mailchain")
					return m
				}(),
			},
			args{
				addressNameServices(func() values.Store {
					m := valuestest.NewMockStore(mockCtrl)
					m.EXPECT().IsSet("name-service-address.base-url").Return(false)
					return m
				}()),
			},
			false,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Network{
				NameServiceAddress: tt.fields.NameServiceAddress,
			}
			got, err := s.ProduceNameServiceAddress(tt.args.ans)
			if (err != nil) != tt.wantErr {
				t.Errorf("Network.ProduceNameServiceAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got == nil) != tt.wantNil {
				t.Errorf("Network.ProduceNameServiceAddress() = %v, want %v", got, tt.wantNil)
				return
			}
		})
	}
}
