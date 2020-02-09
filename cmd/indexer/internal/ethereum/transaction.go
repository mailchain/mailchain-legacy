package ethereum

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/mailchain/mailchain/cmd/indexer/internal/actions"
	"github.com/mailchain/mailchain/cmd/internal/datastore"
	"github.com/mailchain/mailchain/crypto/secp256k1"
	"github.com/mailchain/mailchain/encoding"
	"github.com/mailchain/mailchain/internal/protocols/ethereum"
	"github.com/pkg/errors"
)

type Transaction struct {
	txStore    datastore.TransactionStore
	rawTxStore datastore.RawTransactionStore
	pkStore    datastore.PublicKeyStore

	networkID *big.Int
}

type TxOptions struct {
	Block *types.Block
}

func NewTransactionProcessor(store datastore.TransactionStore, rawStore datastore.RawTransactionStore, pkStore datastore.PublicKeyStore, networkID *big.Int) *Transaction {
	return &Transaction{
		txStore:    store,
		rawTxStore: rawStore,
		pkStore:    pkStore,
		networkID:  networkID,
	}
}

func (t *Transaction) Run(ctx context.Context, protocol, network string, tx interface{}, txOpts actions.TransactionOptions) error {
	ethTx, ok := tx.(*types.Transaction)
	if !ok {
		return errors.New("tx must be go-ethereum/core/types.Transaction")
	}

	opts, ok := txOpts.(*TxOptions)
	if !ok {
		return errors.New("tx must be ethereum.txOptions")
	}

	v, r, s := ethTx.RawSignatureValues()
	var to []byte
	if ethTx.To() != nil {
		to = ethTx.To().Bytes()
	}

	pubKeyBytes, err := ethereum.GetPublicKeyFromTransaction(r, s, v,
		to,
		ethTx.Data(),
		ethTx.Nonce(),
		ethTx.GasPrice(),
		ethTx.Gas(),
		ethTx.Value())
	if err != nil {
		return errors.WithStack(err)
	}

	pubKey, err := secp256k1.PublicKeyFromBytes(pubKeyBytes)
	if err != nil {
		return errors.WithStack(err)
	}

	from, err := t.From(opts.Block.Number(), ethTx)
	if err != nil {
		return errors.Wrapf(err, "failed to get sender from transaction hash: %s", encoding.EncodeHexZeroX(ethTx.Hash().Bytes()))
	}

	if err := t.pkStore.PutPublicKey(ctx, protocol, network, from,
		&datastore.PublicKey{PublicKey: pubKey, BlockHash: opts.Block.Hash().Bytes(), TxHash: ethTx.Hash().Bytes()}); err != nil {
		return errors.WithStack(err)
	}

	storeTx, err := t.ToTransaction(opts.Block, ethTx)
	if err != nil {
		return errors.WithStack(err)
	}

	return actions.StoreTransaction(ctx, t.txStore, t.rawTxStore, protocol, network, storeTx, ethTx)
}

func (t *Transaction) From(blockNo *big.Int, tx *types.Transaction) ([]byte, error) {
	msg, err := tx.AsMessage(types.MakeSigner(&params.ChainConfig{ChainID: t.networkID}, blockNo))
	if err != nil {
		return nil, err
	}
	return msg.From().Bytes(), nil
}

func (t *Transaction) ToTransaction(blk *types.Block, tx *types.Transaction) (*datastore.Transaction, error) {
	if blk.Transaction(tx.Hash()) == nil {
		return nil, errors.New("transaction hash not in block")
	}

	from, err := t.From(blk.Number(), tx)
	if err != nil {
		return nil, err
	}

	gasPrice := tx.GasPrice()
	value := tx.Value()
	gasUsed := big.NewInt(int64(tx.Gas()))

	var to []byte
	if tx.To() != nil {
		to = tx.To().Bytes()
	}

	return &datastore.Transaction{
		From:      from,
		BlockHash: blk.Hash().Bytes(),
		Hash:      tx.Hash().Bytes(),
		Data:      tx.Data(),
		To:        to,
		Value:     *value,
		GasUsed:   *gasUsed,
		GasPrice:  *gasPrice,
	}, nil
}
