package algod

import (
	"github.com/mailchain/mailchain/internal/protocols/algorand"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// New algod client is created that stores config for each supported network.
func New(algodAddress, algodToken string) (*Client, error) {
	return &Client{
		networkConfigs: map[string]networkConfig{
			algorand.Mainnet: {url: "https://api.algoexplorer.io/v2"},
			algorand.Testnet: {url: "https://api.testnet.algoexplorer.io/v2"},
			algorand.Betanet: {url: "https://api.betanet.algoexplorer.io/v2"},
		},
		algodToken: algodToken,
		logger:     log.With().Str("component", "sender").Str("kind", "algod").Logger(),
	}, nil
}

// Client ethereum JSON-RPC2 client that is used to send transactions.
type Client struct {
	networkConfigs map[string]networkConfig
	algodToken     string
	logger         zerolog.Logger
}

type networkConfig struct {
	url string
}

// IsNetworkSupported checks if the network is supported by etherscan API.
func (c *Client) isNetworkSupported(network string) bool {
	_, ok := c.networkConfigs[network]

	return ok
}
