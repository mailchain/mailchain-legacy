package processor

import (
	"context"
	"errors"
	"fmt"
	"github.com/mailchain/mailchain/cmd/indexer/internal/actions"
	"github.com/mailchain/mailchain/cmd/indexer/internal/clients"
	"github.com/mailchain/mailchain/cmd/internal/datastore"
)

var (
	// ErrEndOfChain is returned when the processor reaches the end of the chain.
	ErrEndOfChain = errors.New("end of chain reached")
)

type Sequential struct {
	protocol string
	network  string

	syncStore      datastore.SyncStore
	blockProcessor actions.Block
	blockClient    clients.BlockByNumber

	processedAllBlocks bool
}

func NewSequential(protocol, network string, store datastore.SyncStore, proc actions.Block, client clients.BlockByNumber) *Sequential {
	return &Sequential{
		protocol:           protocol,
		network:            network,
		syncStore:          store,
		blockProcessor:     proc,
		blockClient:        client,
		processedAllBlocks: false,
	}
}

func (s *Sequential) NextBlock(ctx context.Context) error {
	if s.processedAllBlocks {
		return ErrEndOfChain
	}

	blkNo, err := s.syncStore.GetBlockNumber(ctx, s.protocol, s.network)
	if err != nil {
		return err
	}

	fmt.Println("block number: ", blkNo)

	blk, err := s.blockClient.Get(ctx, blkNo)
	if err != nil {
		return err
	}

	if err := s.blockProcessor.Run(ctx, s.protocol, s.network, blk); err != nil {
		return err
	}

	nextBlockNo := blkNo + 1

	fmt.Println("next block number: ", nextBlockNo)

	if err := s.syncStore.PutBlockNumber(ctx, s.protocol, s.network, nextBlockNo); err != nil {
		return err
	}

	return nil
}
