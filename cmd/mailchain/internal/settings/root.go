package settings

import (
	"io"

	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/defaults"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/output"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/values"
	"github.com/mailchain/mailchain/internal/protocols"
	"github.com/mailchain/mailchain/internal/protocols/ethereum"
	"github.com/mailchain/mailchain/internal/protocols/substrate"
)

// FromStore creates root settings from a configuration storage
func FromStore(s values.Store) *Root {
	return &Root{
		DomainNameServices:  domainNameServices(s),
		AddressNameServices: addressNameServices(s),
		Senders:             senders(s),
		Receivers:           receivers(s),
		PublicKeyFinders:    publicKeyFinders(s),
		// Protocols these contain the networks
		Protocols: map[string]*Protocol{
			protocols.Ethereum: protocol(s, protocols.Ethereum, map[string]NetworkClient{
				ethereum.Goerli:  network(s, protocols.Ethereum, ethereum.Goerli, defaults.EthereumNetworkAny()),
				ethereum.Kovan:   network(s, protocols.Ethereum, ethereum.Kovan, defaults.EthereumNetworkAny()),
				ethereum.Mainnet: network(s, protocols.Ethereum, ethereum.Mainnet, defaults.EthereumNetworkAny()),
				ethereum.Rinkeby: network(s, protocols.Ethereum, ethereum.Rinkeby, defaults.EthereumNetworkAny()),
				ethereum.Ropsten: network(s, protocols.Ethereum, ethereum.Ropsten, defaults.EthereumNetworkAny()),
			}),
			protocols.Substrate: protocol(s, protocols.Substrate, map[string]NetworkClient{
				substrate.EdgewareTestnet: network(s, protocols.Substrate, substrate.EdgewareTestnet, defaults.SubstrateNetworkAny()),
			}),
		},
		// other
		Keystore:     keystore(s),
		MailboxState: mailboxState(s),
		SentStore:    sentStore(s),
		Server:       server(s),
	}
}

// Root configuration element.
type Root struct {
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

// ToYaml converts settings to yaml
func (o *Root) ToYaml(out io.Writer, tabsize int, commentDefaults, excludeDefaults bool) {
	protocolElements := []output.Element{}
	for _, v := range o.Protocols {
		protocolElements = append(protocolElements, v.Output())
	}

	output.ToYaml(output.Root{
		Elements: []output.Element{
			{
				FullName: "protocols",
				Elements: protocolElements,
			},

			o.AddressNameServices.Output(),
			o.DomainNameServices.Output(),

			o.PublicKeyFinders.Output(),
			o.Receivers.Output(),
			o.Senders.Output(),

			o.Keystore.Output(),
			o.MailboxState.Output(),
			o.SentStore.Output(),
			o.Server.Output(),
		},
	}, out, 2, commentDefaults, excludeDefaults)
}
