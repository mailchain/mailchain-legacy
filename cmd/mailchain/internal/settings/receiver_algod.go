package settings //nolint: dupl

import (
	"github.com/mailchain/mailchain/cmd/internal/settings/output"
	"github.com/mailchain/mailchain/cmd/internal/settings/values"
	"github.com/mailchain/mailchain/internal/clients/algod"
	"github.com/mailchain/mailchain/internal/mailbox"
	"github.com/mailchain/mailchain/internal/protocols"
	"github.com/mailchain/mailchain/internal/protocols/algorand"
)

// AlgodReceiver configuration element.
type AlgodReceiver struct {
	EnabledProtocolNetworks values.StringSlice
	Token                   values.String
}

func algodReceiver(s values.Store) *AlgodReceiver {
	enabledNetworks := []string{
		protocols.Algorand + "/" + algorand.Mainnet,
		protocols.Algorand + "/" + algorand.Betanet,
		protocols.Algorand + "/" + algorand.Testnet,
	}
	return &AlgodReceiver{
		EnabledProtocolNetworks: values.NewDefaultStringSlice(
			enabledNetworks,
			s,
			"receivers.algod.enabled-networks",
		),
		Token: values.NewDefaultString("", s, "receivers.algod.token"),
	}
}

// Supports a map of what protocol and network combinations are supported.
func (r AlgodReceiver) Supports() map[string]bool {
	m := map[string]bool{}
	for _, np := range r.EnabledProtocolNetworks.Get() {
		m[np] = true
	}
	return m
}

// Produce `mailbox.Receiver` based on configuration settings.
func (r AlgodReceiver) Produce() (mailbox.Receiver, error) {
	return algod.New(r.Token.Get())
}

// Output configuration as an `output.Element` for use in exporting configuration.
func (r AlgodReceiver) Output() output.Element {
	return output.Element{
		FullName: "receivers.algod",
		Attributes: []output.Attribute{
			r.EnabledProtocolNetworks.Attribute(),
			r.Token.Attribute(),
		},
	}
}
