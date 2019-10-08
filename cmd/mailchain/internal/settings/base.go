package settings

import (
	"io"

	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/output"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/values"
	"github.com/mailchain/mailchain/internal/chains"
)

func FromStore(s values.Store) *Base {
	return &Base{
		DomainNameServices:  domainNameServices(s),
		AddressNameServices: addressNameServices(s),
		Senders:             senders(s),
		Receivers:           receivers(s),
		PublicKeyFinders:    publicKeyFinders(s),
		// Protocols these contain the networks
		Protocols: map[string]*Protocol{
			chains.Ethereum: protocol(s, chains.Ethereum),
		},
		// other
		Keystore:     keystore(s),
		MailboxState: mailboxState(s),
		SentStore:    sentStore(s),
		Server:       server(s),
	}
}

type Base struct {
	AddressNameServices *AddressNameServices
	DomainNameServices  *DomainNameServices
	Senders             *Senders
	Receivers           *Receivers
	PublicKeyFinders    *PublicKeyFinders
	// protocol
	Protocols map[string]*Protocol
	// Ethereum  *Protocol
	// other
	Keystore     *Keystore
	MailboxState *MailboxState
	SentStore    *SentStore
	Server       *Server
}

func (o *Base) ToYaml(out io.Writer, tabsize int, commentDefaults, excludeDefaults bool) {
	protocols := []output.Element{}
	for _, v := range o.Protocols {
		protocols = append(protocols, v.Output())
	}

	output.ToYaml(output.Root{
		Elements: append(
			protocols,
			o.AddressNameServices.Output(),
			o.DomainNameServices.Output(),

			o.PublicKeyFinders.Output(),
			o.Receivers.Output(),
			o.Senders.Output(),

			o.Keystore.Output(),
			o.MailboxState.Output(),
			o.SentStore.Output(),
			o.Server.Output(),
		),
	}, out, 2, commentDefaults, excludeDefaults)
}
