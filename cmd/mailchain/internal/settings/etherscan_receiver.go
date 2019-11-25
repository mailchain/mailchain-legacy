package settings //nolint: dupl

import (
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/defaults"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/output"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/values"
	"github.com/mailchain/mailchain/internal/clients/etherscan"
	"github.com/mailchain/mailchain/internal/mailbox"
	"github.com/mailchain/mailchain/internal/protocols/ethereum"
)

// EtherscanReceiver configuration element.
type EtherscanReceiver struct {
	kind                    string
	EnabledProtocolNetworks values.StringSlice
	APIKey                  values.String
}

func etherscanReceiverNoAuth(s values.Store) *EtherscanReceiver {
	return etherscanReceiverAny(s, defaults.ClientEtherscanNoAuth)
}

func etherscanReceiver(s values.Store) *EtherscanReceiver {
	return etherscanReceiverAny(s, defaults.ClientEtherscan)
}

func etherscanReceiverAny(s values.Store, kind string) *EtherscanReceiver {
	enabledNetworks := []string{}
	for _, n := range ethereum.Networks() {
		enabledNetworks = append(enabledNetworks, "ethereum/"+n)
	}
	return &EtherscanReceiver{
		kind: kind,
		EnabledProtocolNetworks: values.NewDefaultStringSlice(
			enabledNetworks,
			s,
			"receivers."+kind+".enabled-networks",
		),
		APIKey: values.NewDefaultString("", s, "receivers."+kind+".api-key"),
	}
}

// Supports a map of what protocol and network combinations are supported.
func (r EtherscanReceiver) Supports() map[string]bool {
	m := map[string]bool{}
	for _, np := range r.EnabledProtocolNetworks.Get() {
		m[np] = true
	}
	return m
}

// Produce `mailbox.Receiver` based on configuration settings.
func (r EtherscanReceiver) Produce() (mailbox.Receiver, error) {
	return etherscan.NewAPIClient(r.APIKey.Get())
}

// Output configuration as an `output.Element` for use in exporting configuration.
func (r EtherscanReceiver) Output() output.Element {
	return output.Element{
		FullName: r.kind,
		Attributes: []output.Attribute{
			r.EnabledProtocolNetworks.Attribute(),
			r.APIKey.Attribute(),
		},
	}
}
