package settings //nolint: dupl

import (
	"fmt"

	"github.com/mailchain/mailchain/cmd/internal/settings/output"
	"github.com/mailchain/mailchain/cmd/internal/settings/values"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/defaults"
	"github.com/mailchain/mailchain/internal/mailbox"
	"github.com/mailchain/mailchain/internal/protocols"
	"github.com/mailchain/mailchain/internal/protocols/algorand"
	"github.com/mailchain/mailchain/internal/protocols/substrate"
)

// AlgorandPublicKeyFinder configuration element.
type AlgorandPublicKeyFinder struct {
	EnabledProtocolNetworks values.StringSlice
	kind                    string
	Disabled                values.Bool
}

func algorandPublicKeyFinder(s values.Store) *AlgorandPublicKeyFinder {
	kind := defaults.AlgorandPublicKeyFinder

	enabledNetworks := []string{}
	for _, n := range algorand.Networks() {
		enabledNetworks = append(enabledNetworks, protocols.Algorand+"/"+n)
	}

	return &AlgorandPublicKeyFinder{
		kind: kind,
		Disabled: values.NewDefaultBool(false, s,
			fmt.Sprintf("public-key-finders.%s.disabled", kind)),
		EnabledProtocolNetworks: values.NewDefaultStringSlice(
			enabledNetworks,
			s,
			"public-key-finders."+kind+".enabled-networks",
		),
	}
}

// Supports a map of what protocol and network combinations are supported.
func (r AlgorandPublicKeyFinder) Supports() map[string]bool {
	m := map[string]bool{}
	for _, np := range r.EnabledProtocolNetworks.Get() {
		m[np] = true
	}

	return m
}

// Produce a `mailbox.PubKeyFinder` base on the configuration.
func (r AlgorandPublicKeyFinder) Produce() (mailbox.PubKeyFinder, error) {
	return substrate.NewPublicKeyFinder(), nil
}

// Output configuration as an `output.Element` for use in exporting configuration.
func (r AlgorandPublicKeyFinder) Output() output.Element {
	return output.Element{
		FullName: "public-key-finders." + r.kind,
		Attributes: []output.Attribute{
			r.Disabled.Attribute(),
			r.EnabledProtocolNetworks.Attribute(),
		},
	}
}
