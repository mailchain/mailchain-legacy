package settings

import (
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/output"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/values"
	"github.com/mailchain/mailchain/internal/protocols/ethereum"
	"github.com/mailchain/mailchain/sender"
	relayer "github.com/mailchain/mailchain/sender/relay"
)

// RelaySender configuration element
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

// Supports a map of what protocol and network combinations are supported.
func (r RelaySender) Supports() map[string]bool {
	m := map[string]bool{}
	for _, np := range r.EnabledProtocolNetworks.Get() {
		m[np] = true
	}
	return m
}

// Produce `sender.Message` based on configuration settings.
func (r RelaySender) Produce() (sender.Message, error) {
	return relayer.NewClient(r.BaseURL.Get())
}

// Output configuration as an `output.Element` for use in exporting configuration.
func (r RelaySender) Output() output.Element {
	return output.Element{
		FullName: "senders.ethereum-relay",
		Attributes: []output.Attribute{
			r.BaseURL.Attribute(),
			r.EnabledProtocolNetworks.Attribute(),
		},
	}
}
