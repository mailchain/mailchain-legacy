package settings //nolint: dupl similar to etherscan but should maintain separate implementation

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
	APIKey                  values.String
}

func blockscoutReceiverNoAuth(s values.Store) *BlockscoutReceiver {
	return blockscoutReceiverAny(s, defaults.ClientBlockscoutNoAuth)
}

func blockscoutReceiverAny(s values.Store, kind string) *BlockscoutReceiver {
	enabledNetworks := []string{"ethereum/" + ethereum.Mainnet}
	return &BlockscoutReceiver{
		kind: kind,
		EnabledProtocolNetworks: values.NewDefaultStringSlice(
			enabledNetworks,
			s,
			"receivers."+kind+".enabled-networks",
		),
		APIKey: values.NewDefaultString("", s, "receivers."+kind+".api-key"),
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
	return blockscout.NewAPIClient(r.APIKey.Get())
}

// Output configuration as an `output.Element` for use in exporting configuration.
func (r BlockscoutReceiver) Output() output.Element {
	return output.Element{
		FullName: r.kind,
		Attributes: []output.Attribute{
			r.EnabledProtocolNetworks.Attribute(),
			r.APIKey.Attribute(),
		},
	}
}
