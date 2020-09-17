package substrate

import (
	"context"

	gsrpc "github.com/centrifuge/go-substrate-rpc-client"
)

func NewRPC(address string) (*BlockClient, error) {
	api, err := gsrpc.NewSubstrateAPI(address)
	if err != nil {
		return nil, err
	}

	return &BlockClient{api: api}, nil
}

type BlockClient struct {
	api *gsrpc.SubstrateAPI
}

func (c *BlockClient) BlockByNumber(ctx context.Context, blockNo uint64) (blk interface{}, err error) {
	blkHash, err := c.api.RPC.Chain.GetBlockHash(blockNo)
	if err != nil {
		return nil, err
	}

	sb, err := c.api.RPC.Chain.GetBlock(blkHash)
	if err != nil {
		return nil, err
	}

	return &sb.Block, nil
}

func (c *BlockClient) LatestBlockNumber(ctx context.Context) (blockNo uint64, err error) {
	signedBlock, err := c.api.RPC.Chain.GetBlockLatest()
	if err != nil {
		return 0, err
	}

	return uint64(signedBlock.Block.Header.Number), nil
}
