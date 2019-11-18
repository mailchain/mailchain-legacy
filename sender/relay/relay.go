package relay

import (
	"context"
	"fmt"
	"strings"

	"github.com/mailchain/mailchain/internal/mailbox/signer"
	"github.com/mailchain/mailchain/internal/protocols"
	"github.com/mailchain/mailchain/internal/protocols/ethereum"
	"github.com/mailchain/mailchain/sender"
	"github.com/mailchain/mailchain/sender/ethrpc2"
	"github.com/pkg/errors"
)

// Send a transaction via the relay.
func (c Client) Send(ctx context.Context, network string, to, from, data []byte, txSigner signer.Signer, opts sender.SendOpts) error {
	s, ok := c.senders[network]
	if !ok {
		return errors.Errorf("no sender found for relay")
	}

	return s.Send(ctx, network, to, from, data, txSigner, opts)
}

// NewClient create new API client
func NewClient(baseURL string) (*Client, error) {
	senders := map[string]sender.Message{}

	for _, network := range []string{ethereum.Mainnet, ethereum.Ropsten, ethereum.Kovan, ethereum.Rinkeby, ethereum.Goerli} {
		client, err := ethrpc2.New(createAddress(baseURL, protocols.Ethereum, network))
		if err != nil {
			return nil, errors.WithStack(err)
		}
		senders[network] = client
	}

	return &Client{senders: senders}, nil
}

func createAddress(baseURL, chain, network string) string {
	return fmt.Sprintf("%s/json-rpc/%s/%s", strings.TrimSuffix(baseURL, "/"), strings.ToLower(chain), strings.ToLower(network))
}

// Client for talking to relay service
type Client struct {
	senders map[string]sender.Message
}
