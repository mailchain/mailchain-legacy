package settings //nolint: dupl

import (
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/defaults"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/output"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/values"
	"github.com/mailchain/mailchain/internal/clients/etherscan"
	"github.com/mailchain/mailchain/internal/mailbox"
	"github.com/mailchain/mailchain/internal/protocols/ethereum"
)

// EtherscanPublicKeyFinder configuration settings.
type EtherscanPublicKeyFinder struct {
	EnabledProtocolNetworks values.StringSlice
	APIKey                  values.String
	kind                    string
}

func etherscanPublicKeyFinderNoAuth(s values.Store) *EtherscanPublicKeyFinder {
	return etherscanPublicKeyFinderAny(s, defaults.ClientEtherscanNoAuth)
}

func etherscanPublicKeyFinder(s values.Store) *EtherscanPublicKeyFinder {
	return etherscanPublicKeyFinderAny(s, defaults.ClientEtherscan)
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

// Supports a map of what protocol and network combinations are supported.
func (r EtherscanPublicKeyFinder) Supports() map[string]bool {
	m := map[string]bool{}
	for _, np := range r.EnabledProtocolNetworks.Get() {
		m[np] = true
	}
	return m
}

// Produce `mailbox.PubKeyFinder` based on configuration settings.
func (r EtherscanPublicKeyFinder) Produce() (mailbox.PubKeyFinder, error) {
	return etherscan.NewAPIClient(r.APIKey.Get())
}

// Output configuration as an `output.Element` for use in exporting configuration.
func (r EtherscanPublicKeyFinder) Output() output.Element {
	return output.Element{
		FullName: "public-key-finders." + r.kind,
		Attributes: []output.Attribute{
			r.EnabledProtocolNetworks.Attribute(),
			r.APIKey.Attribute(),
		},
	}
}
