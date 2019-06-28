package settings

import (
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/values"
	"github.com/mailchain/mailchain/internal/chains/ethereum"
	"github.com/mailchain/mailchain/sender"
	relayer "github.com/mailchain/mailchain/sender/relay"
)

type RelaySender struct {
	EnabledProtocolNetworks values.StringSlice
	BaseURL                 values.String
}

func relaySender(s values.Store, network string) *RelaySender {
	enabledNetworks := []string{}
	for _, n := range ethereum.Networks() {
		enabledNetworks = append(enabledNetworks, "ethereum/"+n)
	}
	return &RelaySender{
		EnabledProtocolNetworks: values.NewDefaultStringSlice(
			enabledNetworks,
			s,
			"senders.ethereum-relay.enabled-networks",
		),
		BaseURL: values.NewDefaultString("https://relay.mailchain.xyz/", s, "senders.ethereum-relay.base-url"),
	}
}

func (r RelaySender) Supports() map[string]bool {
	m := map[string]bool{}
	for _, np := range r.EnabledProtocolNetworks.Get() {
		m[np] = true
	}
	return m
}

func (r RelaySender) Produce() (sender.Message, error) {
	return relayer.NewClient(r.BaseURL.Get())
}
