package substrate

import (
	"bytes"
	"context"

	"github.com/centrifuge/go-substrate-rpc-client/scale"
	"github.com/centrifuge/go-substrate-rpc-client/types"
	"github.com/mailchain/mailchain/cmd/indexer/internal/actions"
	"github.com/mailchain/mailchain/cmd/internal/datastore"
	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/ed25519"
	"github.com/mailchain/mailchain/crypto/secp256k1"
	"github.com/mailchain/mailchain/crypto/sr25519"
	"github.com/mailchain/mailchain/encoding"
	"github.com/mailchain/mailchain/internal/protocols/substrate"
	"github.com/minio/blake2b-simd"
	"github.com/pkg/errors"
)

type Extrinsic struct {
	txStore    datastore.TransactionStore
	rawTxStore datastore.RawTransactionStore
	pkStore    datastore.PublicKeyStore
	// chainConfig *params.ChainConfig
}

func NewExtrinsicProcessor(store datastore.TransactionStore, rawStore datastore.RawTransactionStore, pkStore datastore.PublicKeyStore) *Extrinsic {
	return &Extrinsic{
		txStore:    store,
		rawTxStore: rawStore,
		pkStore:    pkStore,
		// chainConfig: chainConfig,
	}
}

func (t *Extrinsic) Run(ctx context.Context, protocol, network string, tx interface{}, txOpts actions.TransactionOptions) error {
	subEx, ok := tx.(*types.Extrinsic)
	if !ok {
		return errors.Errorf("tx must be github.com/centrifuge/go-substrate-rpc-client/types.Extrinsic")
	}

	opts, ok := txOpts.(*TxOptions)
	if !ok {
		return errors.Errorf("tx must be substrate.ExOptions")
	}

	storeTx, err := t.ToTransaction(network, opts.Block, subEx)
	if err != nil {
		return errors.WithStack(err)
	}

	return actions.StoreTransaction(ctx, t.txStore, t.rawTxStore, protocol, network, storeTx, subEx)
}

func getFromPublicKey(sig *types.ExtrinsicSignatureV4) (crypto.PublicKey, error) {
	if !sig.Signer.IsAccountID {
		return nil, errors.Errorf("must be signed by account ID")
	}

	if sig.Signature.IsEcdsa { //nolint: gocritic
		return secp256k1.PublicKeyFromBytes(sig.Signature.AsEcdsa)
	} else if sig.Signature.IsEd25519 {
		return ed25519.PublicKeyFromBytes(sig.Signer.AsAccountID[:])
	} else if sig.Signature.IsSr25519 {
		return sr25519.PublicKeyFromBytes(sig.Signer.AsAccountID[:])
	} else {
		return nil, errors.Errorf("invalid signature")
	}
}

func getFromAddress(network string, sig *types.ExtrinsicSignatureV4) ([]byte, error) {
	pubKey, err := getFromPublicKey(sig)
	if err != nil {
		return nil, err
	}

	return substrate.SS58AddressFormat(network, pubKey)
}

func getToAddress(network string, dataPart []byte) ([]byte, error) {
	if len(dataPart) < 33 {
		return []byte{}, nil
	}

	toPubKey, err := sr25519.PublicKeyFromBytes(dataPart[1:33])
	if err != nil {
		return nil, err
	}

	return substrate.SS58AddressFormat(network, toPubKey)
}

func ExtrinsicHash(ex *types.Extrinsic) ([]byte, error) {
	w := bytes.NewBuffer([]byte{})
	encoder := scale.NewEncoder(w)

	if err := ex.Encode(*encoder); err != nil {
		return nil, err
	}

	h := blake2b.Sum256(w.Bytes())

	return h[:], nil
}

func (t *Extrinsic) ToTransaction(network string, blk *types.Block, ex *types.Extrinsic) (*datastore.Transaction, error) {
	w := bytes.NewBuffer([]byte{})
	encoder := scale.NewEncoder(w)

	if err := ex.Method.Args.Encode(*encoder); err != nil {
		return nil, err
	}

	txInfo, data := getParts(w.Bytes())

	from, err := getFromAddress(network, &ex.Signature)
	if err != nil {
		return nil, err
	}

	to, err := getToAddress(network, txInfo)
	if err != nil {
		return nil, err
	}

	hash, err := ExtrinsicHash(ex)
	if err != nil {
		return nil, err
	}

	return &datastore.Transaction{
		From: from,
		// BlockHash: blk.Hash().Bytes(),
		Hash: hash,
		Data: data,
		To:   to,
		// Value:    *value,
		// GasUsed:  *gasUsed,
		// GasPrice: *gasPrice,
	}, nil
}

func getParts(data []byte) (txInfo, dataField []byte) {
	// mailchain hex encoded -> 0x6d61696c636861696e
	// 0x6d61696c636861696e hex encoded -> 3078366436313639366336333638363136393665
	pieces := bytes.Split(data, encoding.DataPrefix())
	txInfo = pieces[0]

	if len(pieces) == 2 {
		dataField = append(
			// 0x30, 0x78 -> 0x
			encoding.DataPrefix(),
			pieces[1]...,
		)
	}

	return
}
