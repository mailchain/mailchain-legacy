package fetching

import (
	"context"
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
