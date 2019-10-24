package settings

import (
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/output"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/values"
	"github.com/mailchain/mailchain/internal/mailbox"
	"github.com/mailchain/mailchain/nameservice"
	"github.com/mailchain/mailchain/sender"
)

func protocol(s values.Store, protocol string, networkClients map[string]NetworkClient) *Protocol

type Protocol struct {
	Networks map[string]NetworkClient
	Kind     string
	Disabled values.Bool
}

func (p Protocol) GetSenders(senders *Senders) (map[string]sender.Message, error)

func (p Protocol) GetReceivers(receivers *Receivers) (map[string]mailbox.Receiver, error)

func (p Protocol) GetPublicKeyFinders(publicKeyFinders *PublicKeyFinders) (map[string]mailbox.PubKeyFinder, error)

func (p Protocol) GetAddressNameServices(ans *AddressNameServices) (map[string]nameservice.ReverseLookup, error)

func (p Protocol) GetDomainNameServices(ans *DomainNameServices) (map[string]nameservice.ForwardLookup, error)

func (p Protocol) Output() output.Element
