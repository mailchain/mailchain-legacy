package etherscan

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/mailchain/mailchain/internal/protocols/ethereum"
	"github.com/pkg/errors"
	"gopkg.in/resty.v1"
)

// NewAPIClient create new API client
func NewAPIClient(apiKey string) (*APIClient, error) {
	return &APIClient{
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

// GetTransactionByHash get transaction details from transaction hash via etherscan
func (c APIClient) getTransactionByHash(network string, hash common.Hash) (*types.Transaction, error) {
	config, ok := c.networkConfigs[network]
	if !ok {
		return nil, errors.Errorf("network not supported")
	}
	txListResponse, err := resty.R().
		SetQueryParams(map[string]string{
			"module": "proxy",
			"action": "eth_getTransactionByHash",
			"txhash": hash.Hex(),
			"apikey": c.key,
		}).Get(config.url)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	res := &jsonrpcMessage{}
	if err := json.Unmarshal(txListResponse.Body(), res); err != nil {
		return nil, errors.WithStack(err)
	}
	if res.Error != nil {
		return nil, errors.Errorf(res.Error.Message)
	}
	if res.Result == nil {
		return nil, errors.New("not found")
	}

	var ts *types.Transaction

	if err := json.Unmarshal(res.Result, &ts); err != nil {
		return nil, errors.WithStack(err)
	}
	return ts, nil
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
			"apikey":     c.key,
			"address":    common.BytesToAddress(address).Hex(),
		}).Get(config.url)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	txResult := &txList{}
	if err := json.Unmarshal(txListResponse.Body(), txResult); err != nil {
		return nil, errors.WithStack(err)
	}
	return txResult, nil
}

func (c APIClient) Decode(input string) ([]byte, error) {
	return hexutil.Decode(input)
}
