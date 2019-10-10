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

type EtherscanPublicKeyFinder struct {
	EnabledProtocolNetworks values.StringSlice
	APIKey                  values.String
	kind                    string
}

func etherscanPublicKeyFinderNoAuth(s values.Store) *EtherscanPublicKeyFinder {
	return etherscanPublicKeyFinderAny(s, mailchain.ClientEtherscanNoAuth)
}

func etherscanPublicKeyFinder(s values.Store) *EtherscanPublicKeyFinder {
	return etherscanPublicKeyFinderAny(s, mailchain.ClientEtherscan)
}

func etherscanPublicKeyFinderAny(s values.Store, kind string) *EtherscanPublicKeyFinder {
	enabledNetworks := []string{}
	for _, n := range ethereum.Networks() {
		enabledNetworks = append(enabledNetworks, "ethereum/"+n)
	}
	return &EtherscanPublicKeyFinder{
		kind: kind,
		EnabledProtocolNetworks: values.NewDefaultStringSlice(
			enabledNetworks,
			s,
			"public-key-finders."+kind+".enabled-networks",
		),
		APIKey: values.NewDefaultString("", s, "public-key-finders."+kind+".api-key"),
	}
}

func (r EtherscanPublicKeyFinder) Supports() map[string]bool {
	m := map[string]bool{}
	for _, np := range r.EnabledProtocolNetworks.Get() {
		m[np] = true
	}
	return m
}

func (r EtherscanPublicKeyFinder) Produce() (mailbox.PubKeyFinder, error) {
	return etherscan.NewAPIClient(r.APIKey.Get())
}

func (r EtherscanPublicKeyFinder) Output() output.Element {
	return output.Element{
		FullName: "public-key-finders." + r.kind,
		Attributes: []output.Attribute{
			r.EnabledProtocolNetworks.Attribute(),
			r.APIKey.Attribute(),
		},
	}
}
