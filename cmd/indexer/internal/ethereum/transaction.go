package ethereum

import (
	"context"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/mailchain/mailchain/cmd/indexer/internal/datastore"
	"github.com/mailchain/mailchain/cmd/indexer/internal/processor"
	"github.com/mailchain/mailchain/crypto/secp256k1"
	"github.com/mailchain/mailchain/internal/protocols/ethereum"
)

type Transaction struct {
	txStore    datastore.TransactionStore
	rawTxStore datastore.RawTransactionStore
	pkStore    datastore.PublicKeyStore

	networkID *big.Int
}

type txOptions struct {
	block *types.Block
}

func (t *Transaction) Run(ctx context.Context, protocol, network string, tx interface{}, txOpts processor.TransactionOptions) error {
	// blk *types.Block, ethTx *types.Transaction
	ethTx, ok := tx.(*types.Transaction)
	if !ok {
		return errors.New("tx must be go-ethereum/core/types.Transaction")
	}

	opts, ok := txOpts.(*txOptions)
	if !ok {
		return errors.New("tx must be ethereum.txOptions")
	}

	storeTx, err := t.toTransaction(opts.block, ethTx)
	if err != nil {
		return err
	}

	v, r, s := ethTx.RawSignatureValues()

	pubKeyBytes, err := ethereum.GetPublicKeyFromTransaction(r, s, v,
		ethTx.To().Bytes(),
		ethTx.Data(),
		ethTx.Nonce(),
		ethTx.GasPrice(),
		ethTx.Gas(),
		ethTx.Value())
	if err != nil {
		return err
	}

	pubKey, err := secp256k1.PublicKeyFromBytes(pubKeyBytes)
	if err != nil {
		return err
	}

	if err := t.pkStore.PutPublicKey(ctx, protocol, network, storeTx.From,
		&datastore.PublicKey{PublicKey: pubKey, BlockHash: storeTx.BlockHash, TxHash: storeTx.Hash}); err != nil {
		return err
	}

	return processor.StoreTransaction(ctx, t.txStore, t.rawTxStore, protocol, network, storeTx, ethTx)
}

func (t *Transaction) toTransaction(blk *types.Block, tx *types.Transaction) (*datastore.Transaction, error) {
	msg, err := tx.AsMessage(types.MakeSigner(&params.ChainConfig{ChainID: t.networkID}, blk.Number()))
	if err != nil {
		return nil, err
	}

	gasPrice := tx.GasPrice()
	value := tx.Value()
	gasUsed := big.NewInt(int64(tx.Gas()))

	return &datastore.Transaction{
		From:      msg.From().Bytes(),
		BlockHash: blk.Hash().Bytes(),
		Hash:      tx.Hash().Bytes(),
		Data:      tx.Data(),
		To:        tx.To().Bytes(),
		Value:     *value,
		GasUsed:   *gasUsed,
		GasPrice:  *gasPrice,
	}, nil
}
