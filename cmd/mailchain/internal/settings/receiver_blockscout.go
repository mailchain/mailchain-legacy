package settings //nolint: dupl

import (
	"github.com/mailchain/mailchain/cmd/internal/settings/output"
	"github.com/mailchain/mailchain/cmd/internal/settings/values"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/defaults"
	"github.com/mailchain/mailchain/internal/clients/blockscout"
	"github.com/mailchain/mailchain/internal/mailbox"
	"github.com/mailchain/mailchain/internal/protocols/ethereum"
)

// BlockscoutReceiver configuration element.
type BlockscoutReceiver struct {
	kind                    string
	EnabledProtocolNetworks values.StringSlice
}

func blockscoutReceiverNoAuth(s values.Store) *BlockscoutReceiver {
	kind := defaults.ClientBlockscoutNoAuth

	return &BlockscoutReceiver{
		kind: kind,
		EnabledProtocolNetworks: values.NewDefaultStringSlice(
			[]string{"ethereum/" + ethereum.Mainnet},
			s,
			"receivers."+kind+".enabled-networks",
		),
	}
}

// Supports a map of what protocol and network combinations are supported.
func (r BlockscoutReceiver) Supports() map[string]bool {
	m := map[string]bool{}
	for _, np := range r.EnabledProtocolNetworks.Get() {
		m[np] = true
	}

	return m
}

// Produce `mailbox.Receiver` based on configuration settings.
func (r BlockscoutReceiver) Produce() (mailbox.Receiver, error) {
	return blockscout.NewAPIClient()
}

// Output configuration as an `output.Element` for use in exporting configuration.
func (r BlockscoutReceiver) Output() output.Element {
	return output.Element{
		FullName: r.kind,
		Attributes: []output.Attribute{
			r.EnabledProtocolNetworks.Attribute(),
		},
	}
}
