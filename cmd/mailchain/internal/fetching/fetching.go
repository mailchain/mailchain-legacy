package fetching

import (
	"runtime"
	"time"

	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings"
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

	for protocolName := range config.Protocols {
		logger := log.With().Str("component", "Fetching").Str("protocol", protocolName).Logger()
		protocol := config.Protocols[protocolName]

		if protocol.Disabled.Get() {
			logger.Debug().Msg("protocol disabled skipping")

			continue
		}

		for networkName := range protocol.Networks {
			logger := logger.With().Str("network", networkName).Logger()
			network := protocol.Networks[networkName]

			if network.Disabled() {
				logger.Debug().Msg("network disabled skipping")

				continue
			}

			addresses, err := ks.GetAddresses(protocolName, networkName)
			if err != nil {
				return nil, nil, nil, errors.WithMessagef(err, "failed to get addresses for %s.%s", protocolName, networkName)
			}

			if len(addresses[protocolName][networkName]) == 0 {
				logger.Debug().Msg("no addresses found skipping")

				continue
			}

			addressesProtocolsNetworks[protocolName+"."+networkName] = addresses[protocolName][networkName]

			r, err := network.ProduceReceiver(config.Receivers)
			if err != nil {
				return nil, nil, nil, errors.WithMessagef(err, "failed to get receiver for %s.%s", protocolName, networkName)
			}

			if r == nil {
				logger.Info().Msg("no receiver configured skipping")

				continue
			}

			receiverByKind[r.Kind()] = r
			kindProtocolsNetworks[r.Kind()] = appendListMap(kindProtocolsNetworks, r.Kind(), protocolName+"."+networkName)
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
	case "algod":
		return time.Second * 60, nil
	default:
		log.Logger.Warn().Str("component", "FetchGroup").Str("kind", kind).Msg("unknown kind using default wait time 500 seconds")

		return time.Second * 500, nil
	}
}

func Do(config *settings.Root, inbox stores.State) error {
	if config.Fetcher.Disabled.Get() {
		logger := log.With().Str("component", "fetching").Logger()
		logger.Info().Msg("background fetching disabled")
		return nil
	}

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
