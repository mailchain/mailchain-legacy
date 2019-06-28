package settings

import (
	"fmt"

	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/defaults"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/values"
	"github.com/mailchain/mailchain/internal/mailbox"
	"github.com/mailchain/mailchain/sender"
)

func network(s values.Store, protocol, network string) *Network {
	k := &Network{
		Kind: network,
		PublicKeyFinder: values.NewDefaultString(defaults.EthereumReceiver, s,
			fmt.Sprintf("protocols.%s.networks.%s.public-key-finder", protocol, network)),
		Receiver: values.NewDefaultString(defaults.EthereumReceiver, s,
			fmt.Sprintf("protocols.%s.networks.%s.receiver", protocol, network)),
		Sender: values.NewDefaultString(fmt.Sprintf("%s-relay", protocol), s,
			fmt.Sprintf("protocols.%s.networks.%s.sender", protocol, network)),
		Disabled: values.NewDefaultBool(false, s,
			fmt.Sprintf("protocols.%s.networks.%s.disabled", protocol, network)),
	}
	return k
}

type Network struct {
	Kind            string
	PublicKeyFinder values.String
	Receiver        values.String
	Sender          values.String
	Disabled        values.Bool
}

func (s Network) ProduceSender(senders *Senders) (sender.Message, error) {
	return senders.Produce(s.Sender.Get())
}

func (s Network) ProduceReceiver(receivers *Receivers) (mailbox.Receiver, error) {
	return receivers.Produce(s.Receiver.Get())
}

func (s Network) ProducePublicKeyFinders(publicKeyFinders *PublicKeyFinders) (mailbox.PubKeyFinder, error) {
	return publicKeyFinders.Produce(s.PublicKeyFinder.Get())
}
