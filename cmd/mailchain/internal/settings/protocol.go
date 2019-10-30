package settings

import (
	"fmt"

	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/output"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/values"
	"github.com/mailchain/mailchain/internal/mailbox"
	"github.com/mailchain/mailchain/nameservice"
	"github.com/mailchain/mailchain/sender"
)

func protocol(s values.Store, protocol string, networkClients map[string]NetworkClient) *Protocol {
	return &Protocol{
		Disabled: values.NewDefaultBool(false, s,
			fmt.Sprintf("protocols.%s.disabled", protocol)),
		Kind:     protocol,
		Networks: networkClients,
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
	networkElements := []output.Element{}
	for _, c := range p.Networks {
		networkElements = append(networkElements, c.Output())
	}
	return output.Element{
		FullName: "protocols." + p.Kind,
		Attributes: []output.Attribute{
			p.Disabled.Attribute(),
		},
		Elements: []output.Element{
			{
				FullName: "protocols." + p.Kind + ".networks",
				Elements: networkElements,
			},
		},
	}
}
