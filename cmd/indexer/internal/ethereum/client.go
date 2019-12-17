package ethereum

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
)

func NewRPC(address string) (*Client, error) {
	client, err := ethclient.Dial(address)
	if err != nil {
		return nil, err
	}
	return &Client{client: client}, nil
}

type Client struct {
	client *ethclient.Client
}

func (c *Client) Get(ctx context.Context, blockNo uint64) (blk interface{}, err error) {
	return c.client.BlockByNumber(ctx, big.NewInt(int64(blockNo)))
}
