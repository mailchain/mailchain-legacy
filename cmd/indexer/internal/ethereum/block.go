package ethereum

import (
	"context"
	"errors"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/mailchain/mailchain/cmd/indexer/internal/actions"
)

type Block struct {
	txProcessor actions.Transaction
}

func NewBlockProcessor(tx actions.Transaction) *Block {
	return &Block{
		txProcessor: tx,
	}
}

func (b *Block) Run(ctx context.Context, protocol, network string, blk interface{}) error {
	ethBlk, ok := blk.(*types.Block)
	if !ok {
		return errors.New("tx must be go-ethereum/core/types.Block")
	}

	txs := ethBlk.Transactions()
	for i := range txs {
		if err := b.txProcessor.Run(ctx, protocol, network, txs[i], &TxOptions{Block: ethBlk}); err != nil {
			return err
		}
	}

	return nil
}
