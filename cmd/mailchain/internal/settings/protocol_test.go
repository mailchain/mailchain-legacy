package settings

import (
	"sort"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/values"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/values/valuestest"
	"github.com/stretchr/testify/assert"
)

func TestProtocol_GetSenders(t *testing.T) {
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type fields struct {
		Networks map[string]*Network
		protocol string
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
				map[string]*Network{
					"mainnet": network(
						func() values.Store {
							m := valuestest.NewMockStore(mockCtrl)
							m.EXPECT().IsSet("protocols.ethereum.networks.mainnet.sender").Return(false)
							return m
						}(), "ethereum", "mainnet"),
					"ropsten": network(
						func() values.Store {
							m := valuestest.NewMockStore(mockCtrl)
							m.EXPECT().IsSet("protocols.ethereum.networks.ropsten.sender").Return(false)
							return m
						}(), "ethereum", "ropsten")},
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
				map[string]*Network{
					"mainnet": network(
						func() values.Store {
							m := valuestest.NewMockStore(mockCtrl)
							m.EXPECT().IsSet("protocols.ethereum.networks.mainnet.sender").Return(true)
							m.EXPECT().GetString("protocols.ethereum.networks.mainnet.sender").Return("invalid")
							return m
						}(), "ethereum", "mainnet"),
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
				protocol: tt.fields.protocol,
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
		Networks map[string]*Network
		protocol string
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
				map[string]*Network{
					"mainnet": network(
						func() values.Store {
							m := valuestest.NewMockStore(mockCtrl)
							m.EXPECT().IsSet("protocols.ethereum.networks.mainnet.receiver").Return(false)
							return m
						}(), "ethereum", "mainnet"),
					"ropsten": network(
						func() values.Store {
							m := valuestest.NewMockStore(mockCtrl)
							m.EXPECT().IsSet("protocols.ethereum.networks.ropsten.receiver").Return(false)
							return m
						}(), "ethereum", "ropsten")},
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
				map[string]*Network{
					"mainnet": network(
						func() values.Store {
							m := valuestest.NewMockStore(mockCtrl)
							m.EXPECT().IsSet("protocols.ethereum.networks.mainnet.receiver").Return(true)
							m.EXPECT().GetString("protocols.ethereum.networks.mainnet.receiver").Return("invalid")
							return m
						}(), "ethereum", "mainnet"),
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
				protocol: tt.fields.protocol,
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
		Networks map[string]*Network
		protocol string
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
				map[string]*Network{
					"mainnet": network(
						func() values.Store {
							m := valuestest.NewMockStore(mockCtrl)
							m.EXPECT().IsSet("protocols.ethereum.networks.mainnet.public-key-finder").Return(false)
							return m
						}(), "ethereum", "mainnet"),
					"ropsten": network(
						func() values.Store {
							m := valuestest.NewMockStore(mockCtrl)
							m.EXPECT().IsSet("protocols.ethereum.networks.ropsten.public-key-finder").Return(false)
							return m
						}(), "ethereum", "ropsten")},
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
				map[string]*Network{
					"mainnet": network(
						func() values.Store {
							m := valuestest.NewMockStore(mockCtrl)
							m.EXPECT().IsSet("protocols.ethereum.networks.mainnet.public-key-finder").Return(true)
							m.EXPECT().GetString("protocols.ethereum.networks.mainnet.public-key-finder").Return("invalid")
							return m
						}(), "ethereum", "mainnet"),
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
				protocol: tt.fields.protocol,
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
		s        values.Store
		protocol string
	}
	tests := []struct {
		name     string
		args     args
		wantKeys []string
	}{
		{
			"success",
			args{
				func() values.Store {
					m := valuestest.NewMockStore(mockCtrl)
					return m
				}(),
				"ethereum",
			},
			[]string{"goerli", "kovan", "mainnet", "rinkeby", "ropsten"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := protocol(tt.args.s, tt.args.protocol)
			gotKeys := []string{}
			for k := range got.Networks {
				gotKeys = append(gotKeys, k)
			}
			sort.Strings(gotKeys)
			if !assert.EqualValues(gotKeys, tt.wantKeys) {
				t.Errorf("protocol().Networks = %v, want %v", got, tt.wantKeys)
			}
		})
	}
}
