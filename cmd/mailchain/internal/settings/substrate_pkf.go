package settings //nolint: dupl

import (
	"fmt"

	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/defaults"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/output"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/values"
	"github.com/mailchain/mailchain/internal/mailbox"
	"github.com/mailchain/mailchain/internal/protocols"
	"github.com/mailchain/mailchain/internal/protocols/substrate"
)

// SubstratePublicKeyFinder configuration element.
type SubstratePublicKeyFinder struct {
	EnabledProtocolNetworks values.StringSlice
	kind                    string
	Disabled                values.Bool
}

func substratePublicKeyFinder(s values.Store) *SubstratePublicKeyFinder {
	kind := defaults.SubstratePublicKeyFinder

	enabledNetworks := []string{}
	for _, n := range substrate.Networks() {
		enabledNetworks = append(enabledNetworks, protocols.Substrate+"/"+n)
	}

	return &SubstratePublicKeyFinder{
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
func (r SubstratePublicKeyFinder) Supports() map[string]bool {
	m := map[string]bool{}
	for _, np := range r.EnabledProtocolNetworks.Get() {
		m[np] = true
	}

	return m
}

// Produce a `mailbox.PubKeyFinder` base on the configuration.
func (r SubstratePublicKeyFinder) Produce() (mailbox.PubKeyFinder, error) {
	return substrate.NewPublicKeyFinder(), nil
}

// Output configuration as an `output.Element` for use in exporting configuration.
func (r SubstratePublicKeyFinder) Output() output.Element {
	return output.Element{
		FullName: "public-key-finders." + r.kind,
		Attributes: []output.Attribute{
			r.Disabled.Attribute(),
			r.EnabledProtocolNetworks.Attribute(),
		},
	}
}
