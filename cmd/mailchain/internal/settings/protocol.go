package settings

import (
	"fmt"

	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/output"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/values"
	"github.com/mailchain/mailchain/internal/chains/ethereum"
	"github.com/mailchain/mailchain/internal/mailbox"
	"github.com/mailchain/mailchain/nameservice"
	"github.com/mailchain/mailchain/sender"
)

func protocol(s values.Store, protocol string) *Protocol {
	return &Protocol{
		Disabled: values.NewDefaultBool(false, s,
			fmt.Sprintf("protocols.%s.disabled", protocol)),
		Kind: protocol,
		Networks: map[string]NetworkClient{
			ethereum.Goerli:  network(s, protocol, ethereum.Goerli),
			ethereum.Kovan:   network(s, protocol, ethereum.Kovan),
			ethereum.Mainnet: network(s, protocol, ethereum.Mainnet),
			ethereum.Rinkeby: network(s, protocol, ethereum.Rinkeby),
			ethereum.Ropsten: network(s, protocol, ethereum.Ropsten),
		},
	}
}

type Protocol struct {
	Networks map[string]NetworkClient
	Kind     string
	Disabled values.Bool
}

func (p Protocol) GetSenders(senders *Senders) (map[string]sender.Message, error) {
	msg := map[string]sender.Message{}
	for network, v := range p.Networks {
		s, err := v.ProduceSender(senders)
		if err != nil {
			return nil, err
		}
		msg[p.Kind+"/"+network] = s
	}
	return msg, nil
}

func (p Protocol) GetReceivers(receivers *Receivers) (map[string]mailbox.Receiver, error) {
	msg := map[string]mailbox.Receiver{}
	for network, v := range p.Networks {
		s, err := v.ProduceReceiver(receivers)
		if err != nil {
			return nil, err
		}
		msg[p.Kind+"/"+network] = s
	}
	return msg, nil
}

func (p Protocol) GetPublicKeyFinders(publicKeyFinders *PublicKeyFinders) (map[string]mailbox.PubKeyFinder, error) {
	msg := map[string]mailbox.PubKeyFinder{}
	for network, v := range p.Networks {
		s, err := v.ProducePublicKeyFinders(publicKeyFinders)
		if err != nil {
			return nil, err
		}
		msg[p.Kind+"/"+network] = s
	}
	return msg, nil
}

func (p Protocol) GetAddressNameServices(ans *AddressNameServices) (map[string]nameservice.ReverseLookup, error) {
	msg := map[string]nameservice.ReverseLookup{}
	for network, v := range p.Networks {
		s, err := v.ProduceNameServiceAddress(ans)
		if err != nil {
			return nil, err
		}
		msg[p.Kind+"/"+network] = s
	}
	return msg, nil
}

func (p Protocol) GetDomainNameServices(ans *DomainNameServices) (map[string]nameservice.ForwardLookup, error) {
	msg := map[string]nameservice.ForwardLookup{}
	for network, v := range p.Networks {
		s, err := v.ProduceNameServiceDomain(ans)
		if err != nil {
			return nil, err
		}
		msg[p.Kind+"/"+network] = s
	}
	return msg, nil
}

func (p Protocol) Output() output.Element {
	elements := []output.Element{}
	for _, c := range p.Networks {
		elements = append(elements, c.Output())
	}
	return output.Element{
		FullName: "protocols." + p.Kind,
		Attributes: []output.Attribute{
			p.Disabled.Attribute(),
		},
		Elements: elements,
	}
}
