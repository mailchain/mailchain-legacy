package settings //nolint: dupl

import (
	"github.com/mailchain/mailchain/cmd/internal/settings/output"
	"github.com/mailchain/mailchain/cmd/internal/settings/values"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/defaults"
	"github.com/mailchain/mailchain/internal/clients/blockscout"
	"github.com/mailchain/mailchain/internal/mailbox"
	"github.com/mailchain/mailchain/internal/protocols/ethereum"
)

// BlockscoutPublicKeyFinder configuration settings.
type BlockscoutBalanceFinder struct {
	EnabledProtocolNetworks values.StringSlice
	APIKey                  values.String
	kind                    string
}

func blockscoutBalanceFinderNoAuth(s values.Store) *BlockscoutBalanceFinder {
	kind := defaults.ClientBlockscoutNoAuth

	return &BlockscoutBalanceFinder{
		kind: kind,
		EnabledProtocolNetworks: values.NewDefaultStringSlice(
			[]string{"ethereum/" + ethereum.Mainnet},
			s,
			"balance-finder."+kind+".enabled-networks",
		),
		APIKey: values.NewDefaultString("", s, "balance-finder."+kind+".api-key"),
	}
}

// Supports a map of what protocol and network combinations are supported.
func (r BlockscoutBalanceFinder) Supports() map[string]bool {
	m := map[string]bool{}
	for _, np := range r.EnabledProtocolNetworks.Get() {
		m[np] = true
	}

	return m
}

// Produce `mailbox.PubKeyFinder` based on configuration settings.
func (r BlockscoutBalanceFinder) Produce() (mailbox.BalanceFinder, error) {
	return blockscout.NewAPIClient()
}

// Output configuration as an `output.Element` for use in exporting configuration.
func (r BlockscoutBalanceFinder) Output() output.Element {
	return output.Element{
		FullName: "balance-finders." + r.kind,
		Attributes: []output.Attribute{
			r.EnabledProtocolNetworks.Attribute(),
		},
	}
}
