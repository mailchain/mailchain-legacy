package settings

import (
	"sort"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/defaults"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/values"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/values/valuestest"
	"github.com/mailchain/mailchain/internal/protocols"
	"github.com/mailchain/mailchain/internal/protocols/ethereum"
	"github.com/stretchr/testify/assert"
)

func TestProtocol_GetSenders(t *testing.T) {
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type fields struct {
		Networks map[string]NetworkClient
		Kind     string
	}
	type args struct {
		senders *Senders
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantKeys []string
		wantErr  bool
	}{
		{
			"success",
			fields{
				map[string]NetworkClient{
					"mainnet": network(
						func() values.Store {
							m := valuestest.NewMockStore(mockCtrl)
							m.EXPECT().IsSet("protocols.ethereum.networks.mainnet.sender").Return(false)
							return m
						}(), "ethereum", "mainnet", defaults.EthereumNetworkAny()),
					"ropsten": network(
						func() values.Store {
							m := valuestest.NewMockStore(mockCtrl)
							m.EXPECT().IsSet("protocols.ethereum.networks.ropsten.sender").Return(false)
							return m
						}(), "ethereum", "ropsten", defaults.EthereumNetworkAny())},
				"ethereum",
			},
			args{
				senders(func() values.Store {
					m := valuestest.NewMockStore(mockCtrl)
					m.EXPECT().IsSet("senders.ethereum-relay.base-url").Return(false).Times(2)
					return m
				}()),
			},
			[]string{"ethereum/mainnet", "ethereum/ropsten"},
			false,
		},
		{
			"err-produce",
			fields{
				map[string]NetworkClient{
					"mainnet": network(
						func() values.Store {
							m := valuestest.NewMockStore(mockCtrl)
							m.EXPECT().IsSet("protocols.ethereum.networks.mainnet.sender").Return(true)
							m.EXPECT().GetString("protocols.ethereum.networks.mainnet.sender").Return("invalid")
							return m
						}(), "ethereum", "mainnet", defaults.EthereumNetworkAny()),
				},
				"ethereum",
			},
			args{
				senders(func() values.Store {
					m := valuestest.NewMockStore(mockCtrl)
					return m
				}()),
			},
			[]string{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Protocol{
				Networks: tt.fields.Networks,
				Kind:     tt.fields.Kind,
			}
			got, err := p.GetSenders(tt.args.senders)
			if (err != nil) != tt.wantErr {
				t.Errorf("Protocol.GetSenders() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			gotKeys := []string{}
			for k := range got {
				gotKeys = append(gotKeys, k)
			}
			sort.Strings(gotKeys)
			if !assert.EqualValues(gotKeys, tt.wantKeys) {
				t.Errorf("Protocol.GetSenders() = %v, want %v", got, tt.wantKeys)
			}
		})
	}
}

func TestProtocol_GetReceivers(t *testing.T) {
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type fields struct {
		Networks map[string]NetworkClient
		Kind     string
	}
	type args struct {
		receivers *Receivers
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantKeys []string
		wantErr  bool
	}{
		{
			"success",
			fields{
				map[string]NetworkClient{
					"mainnet": network(
						func() values.Store {
							m := valuestest.NewMockStore(mockCtrl)
							m.EXPECT().IsSet("protocols.ethereum.networks.mainnet.receiver").Return(false)
							return m
						}(), "ethereum", "mainnet", defaults.EthereumNetworkAny()),
					"ropsten": network(
						func() values.Store {
							m := valuestest.NewMockStore(mockCtrl)
							m.EXPECT().IsSet("protocols.ethereum.networks.ropsten.receiver").Return(false)
							return m
						}(), "ethereum", "ropsten", defaults.EthereumNetworkAny())},
				"ethereum",
			},
			args{
				receivers(func() values.Store {
					m := valuestest.NewMockStore(mockCtrl)
					m.EXPECT().IsSet("receivers.etherscan-no-auth.api-key").Return(false).Times(2)
					return m
				}()),
			},
			[]string{"ethereum/mainnet", "ethereum/ropsten"},
			false,
		},
		{
			"err-produce",
			fields{
				map[string]NetworkClient{
					"mainnet": network(
						func() values.Store {
							m := valuestest.NewMockStore(mockCtrl)
							m.EXPECT().IsSet("protocols.ethereum.networks.mainnet.receiver").Return(true)
							m.EXPECT().GetString("protocols.ethereum.networks.mainnet.receiver").Return("invalid")
							return m
						}(), "ethereum", "mainnet", defaults.EthereumNetworkAny()),
				},
				"ethereum",
			},
			args{
				receivers(func() values.Store {
					m := valuestest.NewMockStore(mockCtrl)
					return m
				}()),
			},
			[]string{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Protocol{
				Networks: tt.fields.Networks,
				Kind:     tt.fields.Kind,
			}
			got, err := p.GetReceivers(tt.args.receivers)
			if (err != nil) != tt.wantErr {
				t.Errorf("Protocol.GetReceivers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			gotKeys := []string{}
			for k := range got {
				gotKeys = append(gotKeys, k)
			}
			sort.Strings(gotKeys)
			if !assert.EqualValues(gotKeys, tt.wantKeys) {
				t.Errorf("Protocol.GetReceivers() = %v, want %v", got, tt.wantKeys)
			}
		})
	}
}

func TestProtocol_GetPublicKeyFinders(t *testing.T) {
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type fields struct {
		Networks map[string]NetworkClient
		Kind     string
	}
	type args struct {
		publicKeyFinders *PublicKeyFinders
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantKeys []string
		wantErr  bool
	}{
		{
			"success",
			fields{
				map[string]NetworkClient{
					"mainnet": network(
						func() values.Store {
							m := valuestest.NewMockStore(mockCtrl)
							m.EXPECT().IsSet("protocols.ethereum.networks.mainnet.public-key-finder").Return(false)
							return m
						}(), "ethereum", "mainnet", defaults.EthereumNetworkAny()),
					"ropsten": network(
						func() values.Store {
							m := valuestest.NewMockStore(mockCtrl)
							m.EXPECT().IsSet("protocols.ethereum.networks.ropsten.public-key-finder").Return(false)
							return m
						}(), "ethereum", "ropsten", defaults.EthereumNetworkAny())},
				"ethereum",
			},
			args{
				publicKeyFinders(func() values.Store {
					m := valuestest.NewMockStore(mockCtrl)
					m.EXPECT().IsSet("public-key-finders.etherscan-no-auth.api-key").Return(false).Times(2)
					return m
				}()),
			},
			[]string{"ethereum/mainnet", "ethereum/ropsten"},
			false,
		},
		{
			"err-produce",
			fields{
				map[string]NetworkClient{
					"mainnet": network(
						func() values.Store {
							m := valuestest.NewMockStore(mockCtrl)
							m.EXPECT().IsSet("protocols.ethereum.networks.mainnet.public-key-finder").Return(true)
							m.EXPECT().GetString("protocols.ethereum.networks.mainnet.public-key-finder").Return("invalid")
							return m
						}(), "ethereum", "mainnet", defaults.EthereumNetworkAny()),
				},
				"ethereum",
			},
			args{
				publicKeyFinders(func() values.Store {
					m := valuestest.NewMockStore(mockCtrl)
					return m
				}()),
			},
			[]string{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Protocol{
				Networks: tt.fields.Networks,
				Kind:     tt.fields.Kind,
			}
			got, err := p.GetPublicKeyFinders(tt.args.publicKeyFinders)
			if (err != nil) != tt.wantErr {
				t.Errorf("Protocol.GetPublicKeyFinders() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			gotKeys := []string{}
			for k := range got {
				gotKeys = append(gotKeys, k)
			}
			sort.Strings(gotKeys)
			if !assert.EqualValues(gotKeys, tt.wantKeys) {
				t.Errorf("Protocol.GetPublicKeyFinders() = %v, want %v", got, tt.wantKeys)
			}
		})
	}
}

func Test_protocol(t *testing.T) {
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type args struct {
		s              values.Store
		protocol       string
		networkClients map[string]NetworkClient
	}
	tests := []struct {
		name         string
		args         args
		wantKeys     []string
		wantDisabled bool
	}{
		{
			"success",
			args{
				func() values.Store {
					m := valuestest.NewMockStore(mockCtrl)
					m.EXPECT().IsSet("protocols.ethereum.disabled").Return(false)
					return m
				}(),
				"ethereum",
				map[string]NetworkClient{
					ethereum.Goerli:  network(valuestest.NewMockStore(mockCtrl), protocols.Ethereum, ethereum.Goerli, defaults.EthereumNetworkAny()),
					ethereum.Kovan:   network(valuestest.NewMockStore(mockCtrl), protocols.Ethereum, ethereum.Kovan, defaults.EthereumNetworkAny()),
					ethereum.Mainnet: network(valuestest.NewMockStore(mockCtrl), protocols.Ethereum, ethereum.Mainnet, defaults.EthereumNetworkAny()),
					ethereum.Rinkeby: network(valuestest.NewMockStore(mockCtrl), protocols.Ethereum, ethereum.Rinkeby, defaults.EthereumNetworkAny()),
					ethereum.Ropsten: network(valuestest.NewMockStore(mockCtrl), protocols.Ethereum, ethereum.Ropsten, defaults.EthereumNetworkAny()),
				},
			},
			[]string{"goerli", "kovan", "mainnet", "rinkeby", "ropsten"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := protocol(tt.args.s, tt.args.protocol, tt.args.networkClients)
			gotKeys := []string{}
			for k := range got.Networks {
				gotKeys = append(gotKeys, k)
			}
			sort.Strings(gotKeys)
			if !assert.EqualValues(gotKeys, tt.wantKeys) {
				t.Errorf("protocol().Networks = %v, want %v", got, tt.wantKeys)
			}
			assert.Equal(tt.wantDisabled, got.Disabled.Get())
		})
	}
}

func TestProtocol_GetAddressNameServices(t *testing.T) {
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	type fields struct {
		Networks map[string]NetworkClient
		Kind     string
		Disabled values.Bool
	}
	type args struct {
		ans *AddressNameServices
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantKeys []string
		wantErr  bool
	}{
		{
			"success",
			fields{
				map[string]NetworkClient{
					"mainnet": network(
						func() values.Store {
							m := valuestest.NewMockStore(mockCtrl)
							m.EXPECT().IsSet("protocols.ethereum.networks.mainnet.nameservice-address").Return(false)
							return m
						}(), "ethereum", "mainnet", defaults.EthereumNetworkAny()),
					"ropsten": network(
						func() values.Store {
							m := valuestest.NewMockStore(mockCtrl)
							m.EXPECT().IsSet("protocols.ethereum.networks.ropsten.nameservice-address").Return(false)
							return m
						}(), "ethereum", "ropsten", defaults.EthereumNetworkAny())},
				"ethereum",
				nil,
			},
			args{
				addressNameServices(func() values.Store {
					m := valuestest.NewMockStore(mockCtrl)
					m.EXPECT().IsSet("nameservice-address.base-url").Return(false).Times(2)
					return m
				}()),
			},
			[]string{"ethereum/mainnet", "ethereum/ropsten"},
			false,
		},
		{
			"err-produce",
			fields{
				map[string]NetworkClient{
					"mainnet": network(
						func() values.Store {
							m := valuestest.NewMockStore(mockCtrl)
							m.EXPECT().IsSet("protocols.ethereum.networks.mainnet.nameservice-address").Return(true)
							m.EXPECT().GetString("protocols.ethereum.networks.mainnet.nameservice-address").Return("invalid")
							return m
						}(), "ethereum", "mainnet", defaults.EthereumNetworkAny()),
				},
				"ethereum",
				nil,
			},
			args{
				addressNameServices(func() values.Store {
					m := valuestest.NewMockStore(mockCtrl)
					return m
				}()),
			},
			[]string{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Protocol{
				Networks: tt.fields.Networks,
				Kind:     tt.fields.Kind,
				Disabled: tt.fields.Disabled,
			}
			got, err := p.GetAddressNameServices(tt.args.ans)
			if (err != nil) != tt.wantErr {
				t.Errorf("Protocol.GetAddressNameServices() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			gotKeys := []string{}
			for k := range got {
				gotKeys = append(gotKeys, k)
			}
			sort.Strings(gotKeys)
			if !assert.EqualValues(gotKeys, tt.wantKeys) {
				t.Errorf("Protocol.GetPublicKeyFinders() = %v, want %v", got, tt.wantKeys)
			}
		})
	}
}

func TestProtocol_GetDomainNameServices(t *testing.T) {
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type fields struct {
		Networks map[string]NetworkClient
		Kind     string
		Disabled values.Bool
	}
	type args struct {
		ans *DomainNameServices
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantKeys []string
		wantErr  bool
	}{
		{
			"success",
			fields{
				map[string]NetworkClient{
					"mainnet": network(
						func() values.Store {
							m := valuestest.NewMockStore(mockCtrl)
							m.EXPECT().IsSet("protocols.ethereum.networks.mainnet.nameservice-domain-name").Return(false)
							return m
						}(), "ethereum", "mainnet", defaults.EthereumNetworkAny()),
					"ropsten": network(
						func() values.Store {
							m := valuestest.NewMockStore(mockCtrl)
							m.EXPECT().IsSet("protocols.ethereum.networks.ropsten.nameservice-domain-name").Return(false)
							return m
						}(), "ethereum", "ropsten", defaults.EthereumNetworkAny())},
				"ethereum",
				nil,
			},
			args{
				domainNameServices(func() values.Store {
					m := valuestest.NewMockStore(mockCtrl)
					m.EXPECT().IsSet("nameservice-domain-name.base-url").Return(false).Times(2)
					return m
				}()),
			},
			[]string{"ethereum/mainnet", "ethereum/ropsten"},
			false,
		},
		{
			"err-produce",
			fields{
				map[string]NetworkClient{
					"mainnet": network(
						func() values.Store {
							m := valuestest.NewMockStore(mockCtrl)
							m.EXPECT().IsSet("protocols.ethereum.networks.mainnet.nameservice-domain-name").Return(true)
							m.EXPECT().GetString("protocols.ethereum.networks.mainnet.nameservice-domain-name").Return("invalid")
							return m
						}(), "ethereum", "mainnet", defaults.EthereumNetworkAny()),
				},
				"ethereum",
				nil,
			},
			args{
				domainNameServices(func() values.Store {
					m := valuestest.NewMockStore(mockCtrl)
					return m
				}()),
			},
			[]string{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Protocol{
				Networks: tt.fields.Networks,
				Kind:     tt.fields.Kind,
				Disabled: tt.fields.Disabled,
			}
			got, err := p.GetDomainNameServices(tt.args.ans)
			if (err != nil) != tt.wantErr {
				t.Errorf("Protocol.GetDomainNameServices() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			gotKeys := []string{}
			for k := range got {
				gotKeys = append(gotKeys, k)
			}
			sort.Strings(gotKeys)
			if !assert.EqualValues(gotKeys, tt.wantKeys) {
				t.Errorf("Protocol.GetPublicKeyFinders() = %v, want %v", got, tt.wantKeys)
			}
		})
	}
}
