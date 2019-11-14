package settings

import (
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/output"
	"github.com/mailchain/mailchain/internal/mailbox"
	"github.com/mailchain/mailchain/nameservice"
	"github.com/mailchain/mailchain/sender"
)

//go:generate mockgen -source=producers.go -package=settingstest -destination=./settingstest/producers_mock.go

// Supporter is used to determine if a protocol/network combination is supported by a configured resource.
type Supporter interface {
	Supports() map[string]bool
}

// NetworkClient full configuration for a network.
type NetworkClient interface {
	ProduceNameServiceDomain(dns *DomainNameServices) (nameservice.ForwardLookup, error)
	ProduceNameServiceAddress(ans *AddressNameServices) (nameservice.ReverseLookup, error)
	ProduceSender(senders *Senders) (sender.Message, error)
	ProduceReceiver(receivers *Receivers) (mailbox.Receiver, error)
	ProducePublicKeyFinders(publicKeyFinders *PublicKeyFinders) (mailbox.PubKeyFinder, error)
	Disabled() bool
	Kind() string
	Output() output.Element
}

// SenderClient configuration.
type SenderClient interface {
	Produce() (sender.Message, error)
	Supporter
	Output() output.Element
}

// ReceiverClient configuration.
type ReceiverClient interface {
	Produce() (mailbox.Receiver, error)
	Supporter
	Output() output.Element
}

// PublicKeyFinderClient configuration.
type PublicKeyFinderClient interface {
	Produce() (mailbox.PubKeyFinder, error)
	Supporter
	Output() output.Element
}

// NameServiceAddressClient configuration.
type NameServiceAddressClient interface {
	Produce() (nameservice.ReverseLookup, error)
	Supporter
	Output() output.Element
}

// NameServiceDomainClient configuration.
type NameServiceDomainClient interface {
	Produce() (nameservice.ForwardLookup, error)
	Supporter
	Output() output.Element
}

// SentClient configuration.
type SentClient interface {
	Produce(client string) (sender.Message, error)
}
