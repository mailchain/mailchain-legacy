package settings

import (
	"fmt"

	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/defaults"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/output"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/values"
	"github.com/mailchain/mailchain/internal/mailbox"
	"github.com/mailchain/mailchain/nameservice"
	"github.com/mailchain/mailchain/sender"
)

func network(s values.Store, protocol, network string, nd *defaults.NetworkDefaults) *Network {
	k := &Network{
		kind:     network,
		protocol: protocol,
		NameServiceAddress: values.NewDefaultString(nd.NameServiceAddress, s,
			fmt.Sprintf("protocols.%s.networks.%s.nameservice-address", protocol, network)),
		NameServiceDomainName: values.NewDefaultString(nd.NameServiceDomainName, s,
			fmt.Sprintf("protocols.%s.networks.%s.nameservice-domain-name", protocol, network)),
		PublicKeyFinder: values.NewDefaultString(nd.PublicKeyFinder, s,
			fmt.Sprintf("protocols.%s.networks.%s.public-key-finder", protocol, network)),
		Receiver: values.NewDefaultString(nd.Receiver, s,
			fmt.Sprintf("protocols.%s.networks.%s.receiver", protocol, network)),
		Sender: values.NewDefaultString(nd.Sender, s,
			fmt.Sprintf("protocols.%s.networks.%s.sender", protocol, network)),
		disabled: values.NewDefaultBool(nd.Disabled, s,
			fmt.Sprintf("protocols.%s.networks.%s.disabled", protocol, network)),
	}
	return k
}

type Network struct {
	kind                  string
	protocol              string
	NameServiceAddress    values.String
	NameServiceDomainName values.String
	PublicKeyFinder       values.String
	Receiver              values.String
	Sender                values.String
	disabled              values.Bool
}

func (s *Network) ProduceNameServiceDomain(dns *DomainNameServices) (nameservice.ForwardLookup, error) {
	return dns.Produce(s.NameServiceDomainName.Get())
}

func (s *Network) ProduceNameServiceAddress(ans *AddressNameServices) (nameservice.ReverseLookup, error) {
	return ans.Produce(s.NameServiceAddress.Get())
}

func (s *Network) ProduceSender(senders *Senders) (sender.Message, error) {
	return senders.Produce(s.Sender.Get())
}

func (s *Network) ProduceReceiver(receivers *Receivers) (mailbox.Receiver, error) {
	return receivers.Produce(s.Receiver.Get())
}

func (s *Network) ProducePublicKeyFinders(publicKeyFinders *PublicKeyFinders) (mailbox.PubKeyFinder, error) {
	return publicKeyFinders.Produce(s.PublicKeyFinder.Get())
}

func (s *Network) Disabled() bool {
	return s.disabled.Get()
}

func (s *Network) Kind() string {
	return s.kind
}

func (s *Network) Output() output.Element {
	return output.Element{
		FullName: fmt.Sprintf("protocols.%s.networks.%s", s.protocol, s.kind),
		Attributes: []output.Attribute{
			s.NameServiceAddress.Attribute(),
			s.NameServiceDomainName.Attribute(),
			s.PublicKeyFinder.Attribute(),
			s.Receiver.Attribute(),
			s.Sender.Attribute(),
			s.disabled.Attribute(),
		},
	}
}
