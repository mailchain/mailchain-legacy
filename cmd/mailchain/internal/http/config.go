package http

import (
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings"
	"github.com/mailchain/mailchain/internal/keystore"
	"github.com/mailchain/mailchain/internal/mailbox"
	"github.com/mailchain/mailchain/nameservice"
	"github.com/mailchain/mailchain/sender"
	"github.com/mailchain/mailchain/stores"

	"github.com/pkg/errors"
)

type config struct {
	mailboxStateStore stores.State
	cache             stores.Cache
	keystore          keystore.Store
	addressResolvers  map[string]nameservice.ReverseLookup
	domainResolvers   map[string]nameservice.ForwardLookup
	publicKeyFinders  map[string]mailbox.PubKeyFinder
	receivers         map[string]mailbox.Receiver
	balanceFinders    map[string]mailbox.BalanceFinder
	senders           map[string]sender.Message
	sentStore         stores.Sent
}

//nolint: gocyclo
func produceConfig(s *settings.Root, inbox stores.State) (*config, error) { //nolint: funlen
	keystorage, err := s.Keystore.Produce()
	if err != nil {
		return nil, errors.WithMessage(err, "could not create `keystore`")
	}
	sentStore, err := s.SentStore.Produce()
	if err != nil {
		return nil, errors.WithMessage(err, "Could not config sent store")
	}
	cacheStore, err := s.CacheStore.Produce()
	if err != nil {
		return nil, errors.WithMessage(err, "Could not configure cache")
	}

	nsAddressResolvers := map[string]nameservice.ReverseLookup{}
	nsDomainResolvers := map[string]nameservice.ForwardLookup{}
	publicKeyFinders := map[string]mailbox.PubKeyFinder{}
	balanceFinders := map[string]mailbox.BalanceFinder{}
	receivers := map[string]mailbox.Receiver{}
	senders := map[string]sender.Message{}

	for protocol := range s.Protocols {
		ans, err := s.Protocols[protocol].GetAddressNameServices(s.AddressNameServices)
		if err != nil {
			return nil, errors.WithMessage(err, "could not get address name service")
		}

		for k, v := range ans {
			nsAddressResolvers[k] = v
		}

		dns, err := s.Protocols[protocol].GetDomainNameServices(s.DomainNameServices)
		if err != nil {
			return nil, errors.WithMessage(err, "could not get domain name service")
		}

		for k, v := range dns {
			nsDomainResolvers[k] = v
		}

		name := s.Protocols[protocol].Kind
		protocolPubKeyFinders, err := s.Protocols[protocol].GetPublicKeyFinders(s.PublicKeyFinders)
		if err != nil {
			return nil, errors.WithMessagef(err, "could not get %q public key finders", name)
		}

		for k, v := range protocolPubKeyFinders {
			publicKeyFinders[k] = v
		}

		protocolGetBalances, err := s.Protocols[protocol].GetBalanceFinders(s.BalanceFinders)
		if err != nil {
			return nil, errors.WithMessagef(err, "could not get %q public key finders", name)
		}

		for k, v := range protocolGetBalances {
			balanceFinders[k] = v
		}

		protocolReceivers, err := s.Protocols[protocol].GetReceivers(s.Receivers)
		if err != nil {
			return nil, errors.WithMessagef(err, "Could not get %q receivers", name)
		}

		for k, v := range protocolReceivers {
			receivers[k] = v
		}

		protocolSenders, err := s.Protocols[protocol].GetSenders(s.Senders)
		if err != nil {
			return nil, errors.WithMessagef(err, "Could not get %q senders", name)
		}

		for k, v := range protocolSenders {
			senders[k] = v
		}
	}

	return &config{
		addressResolvers:  nsAddressResolvers,
		domainResolvers:   nsDomainResolvers,
		mailboxStateStore: inbox,
		cache:             cacheStore,
		keystore:          keystorage,
		publicKeyFinders:  publicKeyFinders,
		balanceFinders:    balanceFinders,
		receivers:         receivers,
		senders:           senders,
		sentStore:         sentStore,
	}, nil
}
