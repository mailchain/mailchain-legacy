package settings

import (
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/output"
	"github.com/mailchain/mailchain/internal/mailbox"
	"github.com/mailchain/mailchain/nameservice"
	"github.com/mailchain/mailchain/sender"
)

//go:generate mockgen -source=producers.go -package=settingstest -destination=./settingstest/producers_mock.go

type Supporter interface {
	Supports() map[string]bool
}

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

type SenderClient interface {
	Produce() (sender.Message, error)
	Supporter
	Output() output.Element
}

type ReceiverClient interface {
	Produce() (mailbox.Receiver, error)
	Supporter
	Output() output.Element
}

type PublicKeyFinderClient interface {
	Produce() (mailbox.PubKeyFinder, error)
	Supporter
	Output() output.Element
}

type NameServiceAddressClient interface {
	Produce() (nameservice.ReverseLookup, error)
	Supporter
	Output() output.Element
}

type NameServiceDomainClient interface {
	Produce() (nameservice.ForwardLookup, error)
	Supporter
	Output() output.Element
}

type SentClient interface {
	Produce(client string) (sender.Message, error)
}
