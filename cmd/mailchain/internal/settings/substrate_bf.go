package settings //nolint: dupl

import (
	"fmt"

	"github.com/mailchain/mailchain/cmd/internal/settings/output"
	"github.com/mailchain/mailchain/cmd/internal/settings/values"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/defaults"
	"github.com/mailchain/mailchain/internal/mailbox"
	"github.com/mailchain/mailchain/internal/protocols"
	"github.com/mailchain/mailchain/internal/protocols/substrate"
)

// SubstratePublicKeyFinder configuration element.
type SubstrateBalanceFinder struct {
	EnabledProtocolNetworks values.StringSlice
	kind                    string
	Disabled                values.Bool
}

func substrateBalanceFinder(s values.Store) *SubstrateBalanceFinder {
	kind := defaults.SubstrateBalanceFinder

	enabledNetworks := []string{}
	for _, n := range substrate.Networks() {
		enabledNetworks = append(enabledNetworks, protocols.Substrate+"/"+n)
	}

	return &SubstrateBalanceFinder{
		kind: kind,
		Disabled: values.NewDefaultBool(false, s,
			fmt.Sprintf("balance-finders.%s.disabled", kind)),
		EnabledProtocolNetworks: values.NewDefaultStringSlice(
			enabledNetworks,
			s,
			"balance-finders."+kind+".enabled-networks",
		),
	}
}

// Supports a map of what protocol and network combinations are supported.
func (r SubstrateBalanceFinder) Supports() map[string]bool {
	m := map[string]bool{}
	for _, np := range r.EnabledProtocolNetworks.Get() {
		m[np] = true
	}

	return m
}

// Produce a `mailbox.PubKeyFinder` base on the configuration.
func (r SubstrateBalanceFinder) Produce() (mailbox.Balance, error) {
	return substrate.NewBalanceFinder(), nil
}

// Output configuration as an `output.Element` for use in exporting configuration.
func (r SubstrateBalanceFinder) Output() output.Element {
	return output.Element{
		FullName: "balance-finders." + r.kind,
		Attributes: []output.Attribute{
			r.Disabled.Attribute(),
			r.EnabledProtocolNetworks.Attribute(),
		},
	}
}
