package settings

import (
	"github.com/mailchain/mailchain/internal/mailbox"
	"github.com/mailchain/mailchain/nameservice"
	"github.com/mailchain/mailchain/sender"
)

//go:generate mockgen -source=producers.go -package=settingstest -destination=./settingstest/producers_mock.go

type Supporter interface {
	Supports() map[string]bool
}

type SenderClient interface {
	Produce() (sender.Message, error)
	Supporter
}

type ReceiverClient interface {
	Produce() (mailbox.Receiver, error)
	Supporter
}

type PublicKeyFinderClient interface {
	Produce() (mailbox.PubKeyFinder, error)
	Supporter
}

type NameServiceAddressClient interface {
	Produce() (nameservice.ReverseLookup, error)
	Supporter
}

type SentClient interface {
	Produce(client string) (sender.Message, error)
}
