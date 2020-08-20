package settings //nolint: dupl

import (
	"github.com/mailchain/mailchain/cmd/internal/settings/output"
	"github.com/mailchain/mailchain/cmd/internal/settings/values"
	"github.com/mailchain/mailchain/internal/clients/mailchain"
	"github.com/mailchain/mailchain/internal/mailbox"
	"github.com/mailchain/mailchain/internal/protocols"
	"github.com/mailchain/mailchain/internal/protocols/substrate"
)

// MailchainReceiver configuration element.
type MailchainReceiver struct {
	EnabledProtocolNetworks values.StringSlice
	Address                 values.String
}

func mailchainReceiverNoAuth(s values.Store) *MailchainReceiver {
	return &MailchainReceiver{
		Address: values.NewDefaultString("http://localhost:8081", s, "receivers.mailchain.address"),
		EnabledProtocolNetworks: values.NewDefaultStringSlice(
			[]string{
				protocols.Substrate + "/" + substrate.EdgewareBeresheet,
				protocols.Substrate + "/" + substrate.EdgewareMainnet,
			},
			s,
			"receivers.mailchain.enabled-networks",
		),
	}
}

// Supports a map of what protocol and network combinations are supported.
func (r MailchainReceiver) Supports() map[string]bool {
	m := map[string]bool{}
	for _, np := range r.EnabledProtocolNetworks.Get() {
		m[np] = true
	}

	return m
}

// Produce `mailbox.Receiver` based on configuration settings.
func (r MailchainReceiver) Produce() (mailbox.Receiver, error) {
	return mailchain.NewReceiver(r.Address.Get())
}

// Output configuration as an `output.Element` for use in exporting configuration.
func (r MailchainReceiver) Output() output.Element {
	return output.Element{
		FullName: "mailchain",
		Attributes: []output.Attribute{
			r.EnabledProtocolNetworks.Attribute(),
			r.Address.Attribute(),
		},
	}
}
