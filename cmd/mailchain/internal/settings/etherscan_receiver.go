// nolint: dupl
package settings

import (
	"github.com/mailchain/mailchain"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/output"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/values"
	"github.com/mailchain/mailchain/internal/protocols/ethereum"
	"github.com/mailchain/mailchain/internal/clients/etherscan"
	"github.com/mailchain/mailchain/internal/mailbox"
)

type EtherscanReceiver struct {
	kind                    string
	EnabledProtocolNetworks values.StringSlice
	APIKey                  values.String
}

func etherscanReceiverNoAuth(s values.Store) *EtherscanReceiver {
	return etherscanReceiverAny(s, mailchain.ClientEtherscanNoAuth)
}

func etherscanReceiver(s values.Store) *EtherscanReceiver {
	return etherscanReceiverAny(s, mailchain.ClientEtherscan)
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

func (r EtherscanReceiver) Supports() map[string]bool {
	m := map[string]bool{}
	for _, np := range r.EnabledProtocolNetworks.Get() {
		m[np] = true
	}
	return m
}

func (r EtherscanReceiver) Produce() (mailbox.Receiver, error) {
	return etherscan.NewAPIClient(r.APIKey.Get())
}

func (r EtherscanReceiver) Output() output.Element {
	return output.Element{
		FullName: r.kind,
		Attributes: []output.Attribute{
			r.EnabledProtocolNetworks.Attribute(),
			r.APIKey.Attribute(),
		},
	}
}
