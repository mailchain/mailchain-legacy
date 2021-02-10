package fetching

import (
	"context"
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings"
	"github.com/mailchain/mailchain/internal/addressing"
	"github.com/mailchain/mailchain/internal/mailbox"
	"github.com/mailchain/mailchain/stores"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func getReceivers(config *settings.Root) (receiverByKind map[string]mailbox.Receiver, kindProtocolsNetworks map[string][]string, addressesProtocolsNetworks map[string][][]byte, err error) {
	receiverByKind = map[string]mailbox.Receiver{}
	kindProtocolsNetworks = map[string][]string{}
	addressesProtocolsNetworks = map[string][][]byte{}

	ks, err := config.Keystore.Produce()
	if err != nil {
		return nil, nil, nil, err
	}

	for p := range config.Protocols {
		protocol := config.Protocols[p]
		if protocol.Disabled.Get() {
			continue
		}

		for n := range protocol.Networks {
			network := protocol.Networks[n]
			if network.Disabled() {
				continue
			}

			r, err := network.ProduceReceiver(config.Receivers)
			if err != nil {
				return nil, nil, nil, errors.WithMessagef(err, "failed to get receiver for %s.%s", p, n)
			}

			receiverByKind[r.Kind()] = r
			kindProtocolsNetworks[r.Kind()] = appendListMap(kindProtocolsNetworks, r.Kind(), p+"."+n)

			addresses, err := ks.GetAddresses(p, n)
			if err != nil {
				return nil, nil, nil, errors.WithMessagef(err, "failed to get addresses for %s.%s", p, n)
			}

			addressesProtocolsNetworks[p+"."+n] = addresses
		}
	}

	return receiverByKind, kindProtocolsNetworks, addressesProtocolsNetworks, nil
}

func appendListMap(m map[string][]string, key, value string) []string {
	if items, ok := m[key]; ok {
		return append(items, value)
	}

	return []string{value}
}

func waitByKind(kind string) (time.Duration, error) {
	switch kind {
	case "etherscan":
		return time.Second * 120, nil
	case "etherscan-no-auth":
		return time.Second * 300, nil
	case "mailchain":
		return time.Second * 60, nil
	}

	return time.Nanosecond, errors.Errorf("unknown kind: %s", kind)
}

func Do(config *settings.Root, inbox stores.State) error {
	receiversByKind, kindProtocolsNetworks, addressesProtocolsNetworks, err := getReceivers(config)
	if err != nil {
		return err
	}

	runtime.GOMAXPROCS(len(receiversByKind))

	for kind := range receiversByKind {
		fg, err := NewFetchGroup(config, inbox, receiversByKind[kind], kindProtocolsNetworks[kind], addressesProtocolsNetworks)
		if err != nil {
			return errors.WithMessagef(err, "failed to create fetch group")
		}

		go func(fg *FetchGroup) {
			fg.Fetch()
		}(fg)
	}

	return nil
}

func NewFetchGroup(config *settings.Root, inbox stores.State, receiver mailbox.Receiver, protocolsNetworks []string, addressesProtocolsNetworks map[string][][]byte) (*FetchGroup, error) {
	logger := log.With().Str("component", "FetchGroup").Str("kind", receiver.Kind()).Logger()

	wait, err := waitByKind(receiver.Kind())
	if err != nil {
		return nil, errors.WithMessagef(err, "can't determine wait time")
	}

	logger.Debug().Msgf("initial wait %s", wait)

	b := backoff.NewExponentialBackOff()
	b.InitialInterval = wait / 1000 // Needed as backoff expects microseconds not seconds
	b.MaxElapsedTime = 0            // prevent it from stoping

	return &FetchGroup{
		fetcher: &Fetcher{
			inbox:    inbox,
			receiver: receiver,
		},
		backoff:                    b,
		protocolsNetworks:          protocolsNetworks,
		addressesProtocolsNetworks: addressesProtocolsNetworks,
	}, nil
}

// FetchGroup used to get all messages for a receiver type.
type FetchGroup struct {
	fetcher                    *Fetcher
	backoff                    backoff.BackOff
	protocolsNetworks          []string
	addressesProtocolsNetworks map[string][][]byte
}

func (f *FetchGroup) Fetch() {
	logger := log.With().Str("component", "Fetcher").Str("kind", f.fetcher.receiver.Kind()).Logger()
	for {
		logger.Debug().Msg("starting check")

		for _, protocolNetwork := range f.protocolsNetworks {
			parts := strings.Split(protocolNetwork, ".")
			if len(parts) != 2 {
				logger.Error().Str("protocolNetwork", protocolNetwork).Msg("bad protocols network")
				continue
			}

			protocol, network := parts[0], parts[1]

			for _, addr := range f.addressesProtocolsNetworks[protocolNetwork] {
				subLogger := logger.With().Str("protocol", protocol).Str("network", network).Logger()
				wait := f.backoff.NextBackOff() * 1000 // returning milliseconds needs converting
				subLogger.Debug().Stringer("wait", wait).Msg("waiting")
				time.Sleep(wait * 1000)

				encodedAddress, _, _ := addressing.EncodeByProtocol(addr, protocol)
				subLogger.Debug().Str("encoded address", encodedAddress).Msg("fetching")

				if err := f.fetcher.Fetch(context.Background(), protocol, network, addr); err != nil {
					subLogger.Error().Err(err).Msg("")
					continue
				}

				f.backoff.Reset()
			}
		}
	}
}

type Fetcher struct {
	inbox    stores.State
	receiver mailbox.Receiver
}

func (f *Fetcher) Fetch(ctx context.Context, protocol, network string, addr []byte) error {
	encodedAddress, _, _ := addressing.EncodeByProtocol(addr, protocol)
	logger := log.With().Str("protocol", protocol).Str("network", network).Str("encoded address", encodedAddress).Logger()

	transactions, err := f.receiver.Receive(ctx, protocol, network, addr)
	if mailbox.IsNetworkNotSupportedError(err) {
		return errors.Errorf("network `%s.%s` does not have etherscan client configured", protocol, network)
	}

	if err != nil {
		return errors.WithStack(err)
	}

	logger.Info().Str("found transactions", fmt.Sprint(len(transactions))).Msg("fetched message transactions")

	for i := range transactions {
		tx := transactions[i]
		if err := f.inbox.PutTransaction(protocol, network, addr, tx); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}
