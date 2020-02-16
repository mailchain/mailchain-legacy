package blockscout

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	"github.com/mailchain/mailchain/internal/protocols/ethereum"
	"github.com/pkg/errors"
	"gopkg.in/resty.v1"
)

// NewAPIClient create new API client
func NewAPIClient(apiKey string) (*APIClient, error) {
	return &APIClient{
		key: apiKey,
		networkConfigs: map[string]networkConfig{
			ethereum.Mainnet: {url: "https://blockscout.com/eth/mainnet/api"},
		},
	}, nil
}

// APIClient for talking to etherscan
type APIClient struct {
	key            string
	networkConfigs map[string]networkConfig
}

type networkConfig struct {
	url string
}

// IsNetworkSupported checks if the network is supported by etherscan API
func (c APIClient) isNetworkSupported(network string) bool {
	_, ok := c.networkConfigs[network]
	return ok
}

// GetTransactionsByAddress get transactions from address via etherscan
func (c APIClient) getTransactionsByAddress(network string, address []byte) (*txList, error) {
	config, ok := c.networkConfigs[network]
	if !ok {
		return nil, errors.Errorf("network not supported")
	}
	txListResponse, err := resty.R().
		SetQueryParams(map[string]string{
			"module":     "account",
			"action":     "txlist",
			"startblock": "0",
			"endblock":   "99999999",
			"sort":       "desc",
			"address":    common.BytesToAddress(address).Hex(),
		}).Get(config.url)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	txResult := &txList{}
	if err := json.Unmarshal(txListResponse.Body(), txResult); err != nil {
		return nil, errors.WithStack(err)
	}
	if txResult.Status == "0" {
		return nil, errors.Errorf(txResult.Message)
	}
	return txResult, nil
}
