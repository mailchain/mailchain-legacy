package fetching

import (
	"context"
	"fmt"

	"github.com/mailchain/mailchain/internal/addressing"
	"github.com/mailchain/mailchain/internal/mailbox"
	"github.com/mailchain/mailchain/stores"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

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
