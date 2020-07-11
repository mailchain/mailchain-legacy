package ethereum

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
)

func NewRPC(address string) (*BlockClient, error) {
	client, err := ethclient.Dial(address)
	if err != nil {
		return nil, err
	}

	return &BlockClient{client: client}, nil
}

type BlockClient struct {
	client *ethclient.Client
}

func (c *BlockClient) BlockByNumber(ctx context.Context, blockNo uint64) (blk interface{}, err error) {
	return c.client.BlockByNumber(ctx, big.NewInt(int64(blockNo)))
}

func (c *BlockClient) NetworkID(ctx context.Context) (*big.Int, error) {
	return c.client.NetworkID(ctx)
}

func (c *BlockClient) LatestBlockNumber(ctx context.Context) (blockNo uint64, err error) {
	hdr, err := c.client.HeaderByNumber(ctx, nil)
	if err != nil {
		return 0, err
	}

	return hdr.Number.Uint64(), nil
}

func (c *BlockClient) GetLatest(ctx context.Context) (blk interface{}, err error) {
	return c.client.BlockByNumber(ctx, nil)
}
