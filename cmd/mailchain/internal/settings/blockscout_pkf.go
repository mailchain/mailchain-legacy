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
type BlockscoutPublicKeyFinder struct {
	EnabledProtocolNetworks values.StringSlice
	APIKey                  values.String
	kind                    string
}

func blockscoutPublicKeyFinderNoAuth(s values.Store) *BlockscoutPublicKeyFinder {
	kind := defaults.ClientBlockscoutNoAuth

	return &BlockscoutPublicKeyFinder{
		kind: kind,
		EnabledProtocolNetworks: values.NewDefaultStringSlice(
			[]string{ethereum.Mainnet},
			s,
			"public-key-finders."+kind+".enabled-networks",
		),
		APIKey: values.NewDefaultString("", s, "public-key-finders."+kind+".api-key"),
	}
}

// Supports a map of what protocol and network combinations are supported.
func (r BlockscoutPublicKeyFinder) Supports() map[string]bool {
	m := map[string]bool{}
	for _, np := range r.EnabledProtocolNetworks.Get() {
		m[np] = true
	}

	return m
}

// Produce `mailbox.PubKeyFinder` based on configuration settings.
func (r BlockscoutPublicKeyFinder) Produce() (mailbox.PubKeyFinder, error) {
	return blockscout.NewAPIClient()
}

// Output configuration as an `output.Element` for use in exporting configuration.
func (r BlockscoutPublicKeyFinder) Output() output.Element {
	return output.Element{
		FullName: "public-key-finders." + r.kind,
		Attributes: []output.Attribute{
			r.EnabledProtocolNetworks.Attribute(),
		},
	}
}
