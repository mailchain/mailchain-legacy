package relay

import "github.com/mailchain/mailchain/internal/chains/ethereum"

// NewClient create new API client
func NewClient(apiKey string) (*Client, error) {
	return &Client{
		key: apiKey,
		networkConfigs: map[string]networkConfig{
			ethereum.Mainnet: {url: "https://api.etherscan.io/api"},
			ethereum.Ropsten: {url: "https://api-ropsten.etherscan.io/api"},
			ethereum.Kovan:   {url: "https://api-kovan.etherscan.io/api"},
			ethereum.Rinkeby: {url: "https://api-rinkeby.etherscan.io/api"},
			ethereum.Goerli:  {url: "https://api-goerli.etherscan.io/api"},
		},
	}, nil
}

// Client for talking to etherscan
type Client struct {
	key            string
	networkConfigs map[string]networkConfig
}

type networkConfig struct {
	url string
}
