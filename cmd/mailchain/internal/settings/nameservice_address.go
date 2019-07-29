package settings

import (
	"github.com/mailchain/mailchain"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/values"
	"github.com/mailchain/mailchain/internal/chains/ethereum"
	"github.com/mailchain/mailchain/nameservice"
	"github.com/pkg/errors"
)

func addressNameServices(s values.Store) *AddressNameServices {
	return &AddressNameServices{
		clients: map[string]NameServiceAddressClient{
			mailchain.Mailchain: mailchainAddressNameServices(s),
		},
	}
}

type AddressNameServices struct {
	clients map[string]NameServiceAddressClient
}

func (s AddressNameServices) Produce(client string) (nameservice.ReverseLookup, error) {
	m, ok := s.clients[client]
	if !ok {
		return nil, errors.Errorf("%s not a supported address name service", client)
	}
	return m.Produce()
}

func mailchainAddressNameServices(s values.Store) *MailchainAddressNameServices {
	enabledNetworks := []string{}
	for _, n := range ethereum.Networks() {
		enabledNetworks = append(enabledNetworks, "ethereum/"+n)
	}
	return &MailchainAddressNameServices{
		BaseURL: values.NewDefaultString("https://ns.mailchain.xyz/", s, "name-service-address.base-url"),
		EnabledProtocolNetworks: values.NewDefaultStringSlice(
			enabledNetworks,
			s,
			"name-service-address.mailchain.enabled-networks",
		),
	}
}

type MailchainAddressNameServices struct {
	BaseURL                 values.String
	EnabledProtocolNetworks values.StringSlice
}

func (s MailchainAddressNameServices) Produce() (nameservice.ReverseLookup, error) {
	return nameservice.NewLookupService(s.BaseURL.Get()), nil
}

func (s MailchainAddressNameServices) Supports() map[string]bool {
	m := map[string]bool{}
	for _, np := range s.EnabledProtocolNetworks.Get() {
		m[np] = true
	}
	return m
}
