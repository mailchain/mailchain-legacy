package settings //nolint: dupl

import (
	"github.com/mailchain/mailchain/cmd/internal/settings/output"
	"github.com/mailchain/mailchain/cmd/internal/settings/values"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/defaults"
	"github.com/mailchain/mailchain/internal/clients/etherscan"
	"github.com/mailchain/mailchain/internal/mailbox"
	"github.com/mailchain/mailchain/internal/protocols/ethereum"
)

// EtherscanBalanceFinder configuration settings.
type EtherscanBalanceFinder struct {
	EnabledProtocolNetworks values.StringSlice
	APIKey                  values.String
	kind                    string
}

func etherscanBalanceFinderNoAuth(s values.Store) *EtherscanBalanceFinder {
	return etherscanBalanceFinderAny(s, defaults.ClientEtherscanNoAuth)
}

func etherscanBalanceFinder(s values.Store) *EtherscanBalanceFinder {
	return etherscanBalanceFinderAny(s, defaults.ClientEtherscan)
}

func etherscanBalanceFinderAny(s values.Store, kind string) *EtherscanBalanceFinder {
	enabledNetworks := []string{}
	for _, n := range ethereum.Networks() {
		enabledNetworks = append(enabledNetworks, "ethereum/"+n)
	}
	return &EtherscanBalanceFinder{
		kind: kind,
		EnabledProtocolNetworks: values.NewDefaultStringSlice(
			enabledNetworks,
			s,
			"balance-finders."+kind+".enabled-networks",
		),
		APIKey: values.NewDefaultString("", s, "balance-finders."+kind+".api-key"),
	}
}

// Supports a map of what protocol and network combinations are supported.
func (r EtherscanBalanceFinder) Supports() map[string]bool {
	m := map[string]bool{}
	for _, np := range r.EnabledProtocolNetworks.Get() {
		m[np] = true
	}
	return m
}

// Produce `mailbox.PubKeyFinder` based on configuration settings.
func (r EtherscanBalanceFinder) Produce() (mailbox.BalanceFinder, error) {
	return etherscan.NewAPIClient(r.APIKey.Get())
}

// Output configuration as an `output.Element` for use in exporting configuration.
func (r EtherscanBalanceFinder) Output() output.Element {
	return output.Element{
		FullName: "balance-finders." + r.kind,
		Attributes: []output.Attribute{
			r.EnabledProtocolNetworks.Attribute(),
			r.APIKey.Attribute(),
		},
	}
}
